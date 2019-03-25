// Copyright 2018 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package resource creates test xDS resources
package kubehandler

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	env_api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	env_auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	env_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	env_endpnt "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	env_lsnr "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	env_als "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	env_alf "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	env_hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	env_cache "github.com/envoyproxy/go-control-plane/pkg/cache"
	env_util "github.com/envoyproxy/go-control-plane/pkg/util"

	proxyconf "demius/envoy-proxy-manager/proxy-config"
)

const (
	anyhost   = "0.0.0.0"
	localhost = "127.0.0.1"
	// XdsCluster is the cluster name for the control server (used by non-ADS set-up)
	XdsCluster       = "xds_cluster"
	AccesslogCluster = "accesslog_cluster"
)

// KubeHandler connect envoy-proxy management with k8s info
type KubeHandler struct {
	snapshotVersion uint64
	mu              *sync.Mutex

	coreInterface corev1.CoreV1Interface
	snapshotCache env_cache.SnapshotCache

	namespaces map[string]Namespace
	listeners  []env_cache.Resource
}

// Namespace map of k8s service-names -> config
type Namespace = map[string]ServiceInfo

// ServiceInfo k8s service routing envoy-proxy config
type ServiceInfo struct {
	ServiceName string
	ProxyConfig *proxyconf.Config
	// node-name -> IO
	Endpoints map[string]EndpointInfo
}

// EndpointInfo logical socket address of endpoint
type EndpointInfo struct {
	IP   string
	Port int32
}

// QualifiedServiceInfo contains service info with namespace
type QualifiedServiceInfo struct {
	Namespace   string
	ServiceName string
	ProxyConfig *proxyconf.Config
	Endpoints   []EndpointInfo
}

type certKeysPath struct {
	CertificateChain string
	PrivateKey       string
}

// Certificates for internal IO
var ioKeys = certKeysPath{
	CertificateChain: "/certs/acc.io/acc.io.crt",
	PrivateKey:       "/certs/acc.io/acc.io.key",
}

// New create new KubeHandler
func New(snapshotCache env_cache.SnapshotCache, coreInterface corev1.CoreV1Interface) *KubeHandler {
	var snapshotVersion uint64
	namespaces := map[string]Namespace{}

	mu := &sync.Mutex{}
	listeners := MakeHTTPListeners()

	info := KubeHandler{snapshotVersion, mu, coreInterface, snapshotCache, namespaces, listeners}

	snapshotCache.SetSnapshot("envoy-proxy", info.Generate())

	return &info
}

// ResourceChangeType kinds of resource changing
type ResourceChangeType = int

const (
	// ResourceCreated when service or endpoint created
	ResourceCreated ResourceChangeType = iota + 1
	// ResourceDeleted when service or endpoint deleted
	ResourceDeleted
)

// ChangeTypeString convert changeType enum to string representation
func ChangeTypeString(changeType ResourceChangeType) string {
	if changeType == ResourceCreated {
		return "created"
	}
	return "deleted"
}

// HandleServiceStatusChange handle creation or deletion of service
func (handler *KubeHandler) HandleServiceStatusChange(srv *v1.Service, key string, changeType ResourceChangeType) {
	ct := ChangeTypeString(changeType)

	namespace := srv.Namespace
	name := srv.ObjectMeta.Name
	annotations := srv.Annotations

	log.Printf("Srv: %s.%s -> %s", name, namespace, ct)

	if ariaProxyCOnfig, ok := annotations["aria.io/proxy-config"]; ok {
		log.Printf("   Config: %s", ariaProxyCOnfig)
		proxyconfig, err := proxyconf.Parse([]byte(ariaProxyCOnfig))
		if err == nil {
			ns, ok := handler.namespaces[srv.Namespace]
			if changeType == ResourceCreated {
				// service created
				var endpoints map[string]EndpointInfo
				if ok {
					// namespace exists
					srvInfo, ok := ns[srv.ObjectMeta.Name]
					if ok {
						// service exists
						endpoints = srvInfo.Endpoints
					} else {
						endpoints = map[string]EndpointInfo{}
					}
				} else {
					// new namespace
					ns = Namespace{}
					handler.namespaces[srv.Namespace] = ns
					endpoints = map[string]EndpointInfo{}
				}
				ns[srv.ObjectMeta.Name] =
					ServiceInfo{
						ServiceName: srv.ObjectMeta.Name,
						ProxyConfig: proxyconfig,
						Endpoints:   endpoints,
					}
			} else {
				// service deleted
				ns, ok := handler.namespaces[srv.Namespace]
				if ok {
					delete(ns, srv.ObjectMeta.Name)
					if len(ns) == 0 {
						delete(handler.namespaces, srv.Namespace)
					}
				}
			}
			handler.snapshotCache.SetSnapshot("envoy-proxy", handler.Generate())
		}
	}
}

// HandleEndpointStatusChange handle creation or deletion of endpoint
func (handler *KubeHandler) HandleEndpointStatusChange(endpoints *v1.Endpoints, key string, changeType ResourceChangeType) {
	ns, ok := handler.namespaces[endpoints.Namespace]
	if ok {
		// namespace exists
		srvInfo, ok := ns[endpoints.Name]
		if ok {
			// service for endpoints exists
			// change internal database
			handler.mu.Lock()
			defer handler.mu.Unlock()

			if changeType == ResourceDeleted {
				srvInfo.Endpoints = map[string]EndpointInfo{}
			} else {
				for _, ss := range endpoints.Subsets {
					for _, pp := range ss.Ports {
						for _, a := range ss.Addresses {
							srvInfo.Endpoints[*a.NodeName] = EndpointInfo{
								IP:   a.IP,
								Port: pp.Port,
							}
						}
					}
				}
			}

			ct := ChangeTypeString(changeType)
			fmt.Printf("Endpoints %s: %s, key: %s, ns: %s\n", ct, endpoints.Name, key, endpoints.Namespace)
			handler.snapshotCache.SetSnapshot("envoy-proxy", handler.Generate())
		}
	}
}

// MakeQualifiedSrvInfos expand namespaces with services
func (handler *KubeHandler) MakeQualifiedSrvInfos() []QualifiedServiceInfo {
	var numClusters = 0
	for _, ns := range handler.namespaces {
		numClusters += len(ns)
	}

	clusters := make([]QualifiedServiceInfo, numClusters)

	i := 0
	for ns, nsMap := range handler.namespaces {
		for _, srv := range nsMap {
			endpoints := make([]EndpointInfo, 0, len(srv.Endpoints))
			for _, e := range srv.Endpoints {
				endpoints = append(endpoints, e)
			}

			clusters[i] =
				QualifiedServiceInfo{
					Namespace:   ns,
					ServiceName: srv.ServiceName,
					ProxyConfig: srv.ProxyConfig,
					Endpoints:   endpoints,
				}
			i++
		}
	}

	return clusters
}

// MakeEndpoint creates a localhost endpoint on a given port.
func MakeEndpoint(clusterName, ns string, endpoints []EndpointInfo) *env_api.ClusterLoadAssignment {
	lbEndpoints := make([]env_endpnt.LbEndpoint, len(endpoints))

	for i, ep := range endpoints {

		log.Println("      endpont: ", clusterName, ep.IP, ep.Port)

		lbEndpoints[i] = env_endpnt.LbEndpoint{
			HostIdentifier: &env_endpnt.LbEndpoint_Endpoint{
				Endpoint: &env_endpnt.Endpoint{
					Address: &env_core.Address{
						Address: &env_core.Address_SocketAddress{
							SocketAddress: &env_core.SocketAddress{
								Protocol: env_core.TCP,
								Address:  ep.IP,
								PortSpecifier: &env_core.SocketAddress_PortValue{
									PortValue: uint32(ep.Port),
								},
							},
						},
					},
				},
			},
		}
	}

	return &env_api.ClusterLoadAssignment{
		ClusterName: clusterName + "." + ns,
		Endpoints: []env_endpnt.LocalityLbEndpoints{{
			LbEndpoints: lbEndpoints,
		}},
	}
}

// MakeCluster creates a cluster using either ADS or EDS.
func MakeCluster(clusterName, ns string) *env_api.Cluster {
	edsSource := &env_core.ConfigSource{
		ConfigSourceSpecifier: &env_core.ConfigSource_ApiConfigSource{
			ApiConfigSource: &env_core.ApiConfigSource{
				ApiType: env_core.ApiConfigSource_GRPC,
				GrpcServices: []*env_core.GrpcService{{
					TargetSpecifier: &env_core.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &env_core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
					},
				}},
			},
		},
	}

	return &env_api.Cluster{
		Name:           clusterName + "." + ns,
		ConnectTimeout: 30 * time.Second,

		TlsContext: upstreamTLSContext(&ioKeys),

		ClusterDiscoveryType: &env_api.Cluster_Type{Type: env_api.Cluster_EDS},
		EdsClusterConfig: &env_api.Cluster_EdsClusterConfig{
			EdsConfig: edsSource,
		},
	}
}

// MakeHTTPListeners creates a HTTP listeners for a cluster (redirect for HTTP and HTTPS)
func MakeHTTPListeners() []env_cache.Resource {
	// access log service configuration
	alsConfig := &env_als.FileAccessLog{
		Path: "dev/stdout",
		// Path: "/var/log/envoy/https_access.log",
	}
	alsConfigPbst, err := env_util.MessageToStruct(alsConfig)
	if err != nil {
		panic(err)
	}

	rdsSource := env_core.ConfigSource{}
	rdsSource.ConfigSourceSpecifier = &env_core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &env_core.ApiConfigSource{
			ApiType: env_core.ApiConfigSource_GRPC,
			GrpcServices: []*env_core.GrpcService{{
				TargetSpecifier: &env_core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &env_core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
				},
			}},
		},
	}

	// HTTP filter configuration
	httpsManager := &env_hcm.HttpConnectionManager{
		CodecType:  env_hcm.AUTO,
		StatPrefix: "ingress_https",
		RouteSpecifier: &env_hcm.HttpConnectionManager_Rds{
			Rds: &env_hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: "https-routes",
			},
		},
		HttpFilters: []*env_hcm.HttpFilter{{
			Name: env_util.Router,
		}},
		AccessLog: []*env_alf.AccessLog{{
			// Name:   util.HTTPGRPCAccessLog,
			Name:       env_util.FileAccessLog,
			ConfigType: &env_alf.AccessLog_Config{Config: alsConfigPbst},
		}},
	}

	httpsPbst, err := env_util.MessageToStruct(httpsManager)
	if err != nil {
		panic(err)
	}

	httpsListener := &env_api.Listener{
		Name: "https-listener",
		Address: env_core.Address{
			Address: &env_core.Address_SocketAddress{
				SocketAddress: &env_core.SocketAddress{
					Protocol: env_core.TCP,
					Address:  anyhost,
					PortSpecifier: &env_core.SocketAddress_PortValue{
						PortValue: 443,
					},
				},
			},
		},
		FilterChains: []env_lsnr.FilterChain{{
			Filters: []env_lsnr.Filter{{
				Name:       env_util.HTTPConnectionManager,
				ConfigType: &env_lsnr.Filter_Config{Config: httpsPbst},
			}},

			TlsContext: downstreamTLSContext(&ioKeys),
		}},
	}
	listeners := make([]env_cache.Resource, 1)
	// listeners[0] = redirectListener
	listeners[0] = httpsListener

	return listeners
}

func downstreamTLSContext(keysPath *certKeysPath) *env_auth.DownstreamTlsContext {
	return &env_auth.DownstreamTlsContext{
		CommonTlsContext: commonTLSContext(keysPath),
	}
}

func upstreamTLSContext(keysPath *certKeysPath) *env_auth.UpstreamTlsContext {
	return &env_auth.UpstreamTlsContext{
		CommonTlsContext: commonTLSContext(keysPath),
	}
}

func commonTLSContext(keysPath *certKeysPath) *env_auth.CommonTlsContext {
	return &env_auth.CommonTlsContext{
		TlsCertificates: []*env_auth.TlsCertificate{
			&env_auth.TlsCertificate{
				CertificateChain: &env_core.DataSource{
					Specifier: &env_core.DataSource_Filename{Filename: keysPath.CertificateChain},
				},
				PrivateKey:       &env_core.DataSource{
					Specifier: &env_core.DataSource_Filename{Filename: keysPath.PrivateKey},
				},
			},
		},
	}
}

// Generate produces a snapshot from the parameters.
func (handler *KubeHandler) Generate() env_cache.Snapshot {
	qualifiedServices := handler.MakeQualifiedSrvInfos()
	numClusters := len(qualifiedServices)

	clusters := make([]env_cache.Resource, numClusters)
	endpoints := make([]env_cache.Resource, numClusters)

	for i, qn := range qualifiedServices {
		clusters[i] = MakeCluster(qn.ServiceName, qn.Namespace)
		endpoints[i] = MakeEndpoint(qn.ServiceName, qn.Namespace, qn.Endpoints)
	}

	routes := MakeRoutes(qualifiedServices)

	atomic.AddUint64(&handler.snapshotVersion, 1)
	version := atomic.LoadUint64(&handler.snapshotVersion)

	return env_cache.NewSnapshot(strconv.FormatUint(version, 10), endpoints, clusters, routes, handler.listeners)
}
