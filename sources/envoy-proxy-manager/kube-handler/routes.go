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

	var routesCnt = calcRoutesCnt(services)

	routes := make([]env_route.Route, routesCnt)
	var currentRoute = 0

	for _, srv := range services {
		clusterName := srv.ServiceName + "." + srv.Namespace

		for _, config := range srv.ProxyConfig.Routes {
			duration := calcDuration(config.Route.Timeout)

			routes[currentRoute] =  env_route.Route{
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
			currentRoute++		
		}

	}

	routeResources := make([]env_cache.Resource, 1)
	routeResources[0] = &env_api.RouteConfiguration{
		Name: "https-routes",
		VirtualHosts: []env_route.VirtualHost{
			{
				Name:    "backend-vhost",
				Domains: []string{"*"},
				Routes:  routes,
			},
		},
	}
	return routeResources
}

func calcRoutesCnt(services []QualifiedServiceInfo) int {
	var cnt = 0
	for _, srv := range services {
		cnt += len(srv.ProxyConfig.Routes)
	}
	return cnt
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
