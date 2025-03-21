package lb

import (
	"sigmaos/simulation/simms"
)

// Round-robin load balancer
type RoundRobinLB struct {
}

// TODO: fix, and eventually subsume into OmniscientLB
func NewRoundRobinLB(*uint64, simms.LoadBalancerStateCache, simms.NewLoadBalancerMetricFn, simms.AssignRequestsToLoadBalancerShardsFn) simms.LoadBalancer {
	return &RoundRobinLB{}
}

func (rr *RoundRobinLB) SteerRequests(reqs []*simms.Request, instances []*simms.MicroserviceInstance) [][]*simms.Request {
	steeredReqs := make([][]*simms.Request, len(instances))
	for i := range steeredReqs {
		steeredReqs[i] = []*simms.Request{}
	}
	lastIdx := 0
	// For each request
	for _, r := range reqs {
		// Find a ready instance to process that request
		for instanceIdx := range instances {
			idx := (lastIdx + 1 + instanceIdx) % len(instances)
			if instances[idx].IsReady() {
				// For the next request, start at the following instance
				lastIdx = idx
				steeredReqs[idx] = append(steeredReqs[idx], r)
				break
			}
		}
	}
	return steeredReqs
}
