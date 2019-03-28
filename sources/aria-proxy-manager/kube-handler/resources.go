package kubehandler

import (
	"time"

	log "github.com/sirupsen/logrus"

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
)

type certKeysPath struct {
	CertificateChain string
	PrivateKey       string
}

// Certificates for internal IO
var accKeys = certKeysPath{
	CertificateChain: "/certs/acc.md/acc.md.crt",
	PrivateKey:       "/certs/acc.md/acc.md.key",
}

var ioKeys = certKeysPath{
	CertificateChain: "/certs/acc.io/acc.io.crt",
	PrivateKey:       "/certs/acc.io/acc.io.key",
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
func MakeCluster(clusterName, ns string, listenerIdx int) *env_api.Cluster {
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

		TlsContext: upstreamTLSContext(listenerIdx),

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

	listeners := make([]env_cache.Resource, 2)
	for i := 0; i <= 1; i++ {
		routeConfigHame := calcRouteConfigName(i)
		listenerName := calcListenerName(i)
		listenerPort := calcListenerPort(i)

		// HTTP filter configuration
		httpsManager := &env_hcm.HttpConnectionManager{
			CodecType:  env_hcm.AUTO,
			StatPrefix: "ingress_https",
			RouteSpecifier: &env_hcm.HttpConnectionManager_Rds{
				Rds: &env_hcm.Rds{
					ConfigSource:    rdsSource,
					RouteConfigName: routeConfigHame,
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

		// HTTP listener configuration
		listeners[i] = &env_api.Listener{
			Name: listenerName,
			Address: env_core.Address{
				Address: &env_core.Address_SocketAddress{
					SocketAddress: &env_core.SocketAddress{
						Protocol: env_core.TCP,
						Address:  anyhost,
						PortSpecifier: &env_core.SocketAddress_PortValue{
							PortValue: listenerPort,
						},
					},
				},
			},
			FilterChains: []env_lsnr.FilterChain{{
				Filters: []env_lsnr.Filter{{
					Name:       env_util.HTTPConnectionManager,
					ConfigType: &env_lsnr.Filter_Config{Config: httpsPbst},
				}},

				TlsContext: downstreamTLSContext(i),
			}},
		}
	}

	return listeners
}

// first route is ACC, second route is IO
func calcListenerName(idx int) string {
	switch idx {
	case 0: return "acc-listener"
	case 1: return "io-listener"
	}
	return ""
}

// first route is ACC, second route is IO
func calcListenerPort(idx int) uint32 {
	switch idx {
	case 0: return 443
	case 1: return 8081
	}
	return 443
}

func downstreamTLSContext(listenerIdx int) *env_auth.DownstreamTlsContext {
	return &env_auth.DownstreamTlsContext{
		CommonTlsContext: commonTLSContext(listenerIdx),
	}
}

func upstreamTLSContext(listenerIdx int) *env_auth.UpstreamTlsContext {
	return &env_auth.UpstreamTlsContext{
		CommonTlsContext: commonTLSContext(listenerIdx),
	}
}

func commonTLSContext(listenerIdx int) *env_auth.CommonTlsContext {
	var keys *certKeysPath
	switch listenerIdx {
	case 0: keys = &accKeys
	case 1: keys = &ioKeys
	}
	return &env_auth.CommonTlsContext{
		TlsCertificates: []*env_auth.TlsCertificate{
			&env_auth.TlsCertificate{
				CertificateChain: &env_core.DataSource{
					Specifier: &env_core.DataSource_Filename{Filename: keys.CertificateChain},
				},
				PrivateKey: &env_core.DataSource{
					Specifier: &env_core.DataSource_Filename{Filename: keys.PrivateKey},
				},
			},
		},
	}
}
