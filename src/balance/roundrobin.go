/**
 * roundrobin.go - roundrobin balance impl
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

package balance

import (
	"errors"
	"sort"
	"../logging"
	"../core"
)

/**
 * Roundrobin balancer
 */
type RoundrobinBalancer struct {

	/* Current backend position */
	current int
}

/**
 * Elect backend using roundrobin strategy
 */
func (b *RoundrobinBalancer) Elect(context core.Context, backends []*core.Backend) (*core.Backend, error) {

	log := logging.For("RoundrobinBalancer")
	if len(backends) == 0 {
		return nil, errors.New("Can't elect backend, Backends empty")
	}

	for _, backend := range backends {
		if backend.Priority  <= 0 {
			return nil, errors.New("Invalid backend Priority 0")
		}
	}

	sorted := make([]*core.Backend, len(backends))
	copy(sorted, backends)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})

	if len(sorted) > 1 {
		if sorted[0].Priority > sorted[1].Priority {
			backend := sorted[0]
			b.current = 0
			log.Debug("elected backend :",backend)
			return backend, nil
		}

		// var pri_sorted []*core.Backend
		// for i := 0; i < len(sorted); i++ {
		// 	if sorted[i].Priority == sorted[0].Priority {
		// 		pri_sorted = append(pri_sorted, sorted[i])
		// 	}
		// }
		// log.Debug("sorted Priority :",pri_sorted)
		// sorted = pri_sorted
		for i := 0; i < len(sorted); i++ {
			if sorted[i].Priority != sorted[0].Priority {
				sorted = sorted[:i]
				break
			}
		}
		}

	log.Debug("Priority sorted backends :",sorted)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Target.String() < sorted[j].Target.String()
	})

	log.Debug("rr backends :",sorted)
	if b.current >= len(sorted) {
		b.current = 0
	}

	backend := sorted[b.current]
	b.current += 1
	log.Debug("elected backend :",backend)
	return backend, nil
}
