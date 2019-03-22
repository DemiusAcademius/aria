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
package resources

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	anyhost   = "0.0.0.0"
	localhost = "127.0.0.1"
	// XdsCluster is the cluster name for the control server (used by non-ADS set-up)
	XdsCluster       = "xds_cluster"
	AccesslogCluster = "accesslog_cluster"
)

type ServicesInfo struct {
	snapshotVersion uint64
	mu              *sync.Mutex

	coreInterface corev1.CoreV1Interface
	// snapshotCache cache.SnapshotCache

	namespaces map[string]ServicesMap
	// listeners  []cache.Resource
}

type ServicesMap = map[string]ServiceInfo

type ServiceInfo struct {
	ServiceName  string
	FrontendType string
	// node-name -> IO
	Endpoints map[string]EndpointInfo
}

type EndpointInfo struct {
	IP   string
	Port int32
}

type QualifiedServiceName struct {
	Namespace    string
	ServiceName  string
	FrontendType string
	Endpoints    []EndpointInfo
}

func NewServicesInfo(coreInterface corev1.CoreV1Interface) *ServicesInfo {
	var snapshotVersion uint64
	namespaces := map[string]ServicesMap{}

	mu := &sync.Mutex{}
	// listeners := MakeHTTPListeners()

	info := ServicesInfo{snapshotVersion, mu, coreInterface, namespaces}

	// snapshotCache.SetSnapshot("envoy-proxy", info.Generate())

	return &info
}

type ResourceChangeType = int

const (
	ResourceCreated ResourceChangeType = iota + 1
	ResourceDeleted
)

func ChangeTypeString(changeType ResourceChangeType) string {
	if changeType == ResourceCreated {
		return "created"
	}
	return "deleted"
}

func (info *ServicesInfo) HandleServiceStatusChange(srv *v1.Service, key string, changeType ResourceChangeType) {
	ct := ChangeTypeString(changeType)

	namespace := srv.Namespace
	name := srv.ObjectMeta.Name
	annotations := srv.Annotations

	log.Printf("Srv: %s.%s -> %s", name, namespace, ct)

	if ariaProxyCOnfig, ok := annotations["aria-proxy/config"]; ok {		
		log.Printf("   Config: %s", ariaProxyCOnfig)
	}
}

func (info *ServicesInfo) HandleEndpointStatusChange(endpoints *v1.Endpoints, key string, changeType ResourceChangeType) {
	ns, ok := info.namespaces[endpoints.Namespace]
	if ok {
		// namespace exists
		srvInfo, ok := ns[endpoints.Name]
		if ok {
			// service for endpoints exists
			// change internal database
			info.mu.Lock()
			defer info.mu.Unlock()

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
			// info.snapshotCache.SetSnapshot("envoy-proxy", info.Generate())
		}
	}
}

func (info *ServicesInfo) MakeQualifiedNames() []QualifiedServiceName {
	var numClusters = 0
	for _, ns := range info.namespaces {
		numClusters += len(ns)
	}

	clusters := make([]QualifiedServiceName, numClusters)

	i := 0
	for ns, nsMap := range info.namespaces {
		for _, srv := range nsMap {
			endpoints := make([]EndpointInfo, 0, len(srv.Endpoints))
			for _, e := range srv.Endpoints {
				endpoints = append(endpoints, e)
			}

			clusters[i] =
				QualifiedServiceName{
					Namespace:    ns,
					ServiceName:  srv.ServiceName,
					FrontendType: srv.FrontendType,
					Endpoints:    endpoints,
				}
			i++
		}
	}

	return clusters
}

/*
// MakeEndpoint creates a localhost endpoint on a given port.
func MakeEndpoint(clusterName, ns string, endpoints []EndpointInfo) *v2.ClusterLoadAssignment {
	lbEndpoints := make([]endpoint.LbEndpoint, len(endpoints))

	for i, ep := range endpoints {

		log.Println("      endpont: ", clusterName,ep.IP,ep.Port)

		lbEndpoints[i] = endpoint.LbEndpoint{
			Endpoint: &endpoint.Endpoint{
				Address: &core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Protocol: core.TCP,
							Address:  ep.IP,
							PortSpecifier: &core.SocketAddress_PortValue{
								PortValue: uint32(ep.Port),
							},
						},
					},
				},
			},
		}
	}

	return &v2.ClusterLoadAssignment{
		ClusterName: clusterName + "." + ns,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: lbEndpoints,
		}},
	}
}

// MakeCluster creates a cluster using either ADS or EDS.
func MakeCluster(clusterName, ns string) *v2.Cluster {
	edsSource := &core.ConfigSource{
		ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
			ApiConfigSource: &core.ApiConfigSource{
				ApiType: core.ApiConfigSource_GRPC,
				GrpcServices: []*core.GrpcService{{
					TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
					},
				}},
			},
		},
	}

	return &v2.Cluster{
		Name:           clusterName + "." + ns,
		ConnectTimeout: 30 * time.Second,

		TlsContext: &auth.UpstreamTlsContext{
			CommonTlsContext: &auth.CommonTlsContext{
				TlsCertificates: []*auth.TlsCertificate{
					&auth.TlsCertificate{
						CertificateChain: &core.DataSource{Specifier: &core.DataSource_Filename{Filename: "/etc/envoy/acc.io.crt"}},
						PrivateKey:       &core.DataSource{Specifier: &core.DataSource_Filename{Filename: "/etc/envoy/acc.io.key"}},
					},
				},
			},
		},

		Type: v2.Cluster_EDS,
		EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
			EdsConfig: edsSource,
		},
	}
}

// MakeHTTPListeners creates a HTTP listeners for a cluster (redirect for HTTP and HTTPS)
func MakeHTTPListeners() []cache.Resource {
	// access log service configuration
	alsConfig := &als.FileAccessLog{
		// Path: "dev/stdout",
		Path: "/var/log/envoy/https_access.log",
	}
	alsConfigPbst, err := util.MessageToStruct(alsConfig)
	if err != nil {
		panic(err)
	}

	rdsSource := core.ConfigSource{}
	rdsSource.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			ApiType: core.ApiConfigSource_GRPC,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
				},
			}},
		},
	}

	// HTTP filter configuration
	httpsManager := &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: "ingress_https",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: "https-routes",
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: util.Router,
		}},
		AccessLog: []*alf.AccessLog{{
			// Name:   util.HTTPGRPCAccessLog,
			Name:   util.FileAccessLog,
			Config: alsConfigPbst,
		}},
	}

	httpsPbst, err := util.MessageToStruct(httpsManager)
	if err != nil {
		panic(err)
	}

	httpsListener := &v2.Listener{
		Name: "https-listener",
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  anyhost,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 443,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:   util.HTTPConnectionManager,
				Config: httpsPbst,
			}},

			TlsContext: &auth.DownstreamTlsContext{
				CommonTlsContext: &auth.CommonTlsContext{
					TlsCertificates: []*auth.TlsCertificate{
						&auth.TlsCertificate{
							CertificateChain: &core.DataSource{Specifier: &core.DataSource_Filename{Filename: "/etc/envoy/acc.io.crt"}},
							PrivateKey:       &core.DataSource{Specifier: &core.DataSource_Filename{Filename: "/etc/envoy/acc.io.key"}},
						},
					},
				},
			},
		}},
	}
	listeners := make([]cache.Resource, 1)
	// listeners[0] = redirectListener
	listeners[0] = httpsListener

	return listeners
}

// Generate produces a snapshot from the parameters.
func (info *ServicesInfo) Generate() cache.Snapshot {
	qualifiedNames := info.MakeQualifiedNames()
	numClusters := len(qualifiedNames)

	clusters := make([]cache.Resource, numClusters)
	endpoints := make([]cache.Resource, numClusters)

	for i, qn := range qualifiedNames {
		clusters[i] = MakeCluster(qn.ServiceName, qn.Namespace)
		endpoints[i] = MakeEndpoint(qn.ServiceName, qn.Namespace, qn.Endpoints)
	}

	routes := MakeRoutes(qualifiedNames)

	// atomic.AddUint64(&info.snapshotVersion, 1)
	// version := atomic.LoadUint64(&info.snapshotVersion)
	info.snapshotVersion++

	return cache.NewSnapshot(strconv.FormatUint(info.snapshotVersion, 10), endpoints, clusters, routes, info.listeners)
}
*/