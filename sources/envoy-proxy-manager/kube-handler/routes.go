package kubehandler

import (
	"time"

	env_api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	env_route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	env_cache "github.com/envoyproxy/go-control-plane/pkg/cache"
)

const (
	// ShortDuration 10 seconds response timeout
	ShortDuration = 10 * time.Second
	// NormalDuraion 30 seconds response timeout
	NormalDuraion = 30 * time.Second
	// LongDuration 5 minutes response timeout
	LongDuration = 5 * time.Minute
	// InfiniteDuration 1 hour response timeout
	InfiniteDuration = 1 * time.Hour
)

// MakeRoutes creates an HTTP route that routes to a given cluster.
func MakeRoutes(services []QualifiedServiceInfo) []env_cache.Resource {
	// first route is ACC, second route is IO
	var routesCnt = calcRoutesCnt(services)
	routes := make([][]env_route.Route, 2)
	curRoutes := make([]int, 2)

	for i := 0; i <= 1; i++ {
		routes[i] = make([]env_route.Route, routesCnt[i])
		curRoutes[i] = 0
	}

	for _, srv := range services {
		clusterName := srv.ServiceName + "." + srv.Namespace

		for _, config := range srv.ProxyConfig.Routes {
			idx := calcListenerIdx(srv.ProxyConfig.Listener)
			if idx > -1 {
				duration := calcDuration(config.Route.Timeout)

				routes[idx][curRoutes[idx]] =  env_route.Route{
					Match: env_route.RouteMatch{
						PathSpecifier: &env_route.RouteMatch_Prefix{
							Prefix: config.Match.Prefix,
						},
					},
					Action: &env_route.Route_Route{
						Route: &env_route.RouteAction{
							PrefixRewrite:    config.Route.PrefixRewrite,
							ClusterSpecifier: &env_route.RouteAction_Cluster{Cluster: clusterName},
							Timeout:          &duration,
						},
					},
				}
				curRoutes[idx]++			
			}
		}

	}

	routeResources := make([]env_cache.Resource, 2)
	for i := 0; i <= 1; i++ {
		routeResources[i] = &env_api.RouteConfiguration{
			Name: calcRouteConfigName(i),
			VirtualHosts: []env_route.VirtualHost{
				{
					Name:    calcVHostName(i),
					Domains: []string{"*"},
					Routes:  routes[i],
				},
			},
		}
		}
	return routeResources
}

// first route is ACC, second route is IO
func calcRoutesCnt(services []QualifiedServiceInfo) []int {
	var cnt = []int {0,0}
	for _, srv := range services {
		idx := calcListenerIdx(srv.ProxyConfig.Listener)
		if idx > -1 {
			cnt[idx] += len(srv.ProxyConfig.Routes);
		}
	}
	return cnt
}

// first route is ACC, second route is IO
func calcListenerIdx(listener string) int {
	switch listener {
	case "acc": return 0
	case "io": return 1
	}
	return -1;
}

// first route is ACC, second route is IO
func calcRouteConfigName(idx int) string {
	switch idx {
	case 0: return "acc-routes"
	case 1: return "io-routes"
	}
	return ""
}

// first route is ACC, second route is IO
func calcVHostName(idx int) string {
	switch idx {
	case 0: return "backend-acc"
	case 1: return "backend-io"
	}
	return ""
}


func calcDuration(timeout string) time.Duration {
	switch timeout {
	case "short": return ShortDuration
	case "normal": return NormalDuraion
	case "long": return LongDuration
	case "infinite": return InfiniteDuration
	}
	return NormalDuraion
}
