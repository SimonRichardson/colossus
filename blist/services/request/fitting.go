package request

import (
	"math/rand"
	"time"

	"github.com/SimonRichardson/colossus/blist/permitters"
)

type Provider func([]Cluster) Fitting

type Fitting interface {
	Clusters() []Cluster
	Required() Required
}

type Required []int

func (r Required) Len() int {
	return len(r)
}

func (r Required) Contains(i int) bool {
	for _, v := range r {
		if v == i {
			return true
		}
	}
	return false
}

func First() Provider {
	return func(clusters []Cluster) Fitting {
		return fitting{
			clusters: clusters,
			required: Required([]int{0}),
		}
	}
}

func Random(index int, probability float64) Provider {
	return func(clusters []Cluster) Fitting {
		c := filter(clusters, func(cluster Cluster) bool {
			if cluster.Index() == index {
				return true
			}

			return rand.Float64() > probability
		})

		return fitting{
			clusters: c,
			required: Required([]int{index}),
		}
	}
}

func RateLimit(index int, duration time.Duration) Provider {
	permits := permitters.New(int64(1), duration)
	return func(clusters []Cluster) Fitting {
		c := clusters
		if !permits.Allowed(1) {
			c = filter(clusters, func(cluster Cluster) bool {
				return cluster.Index() == index
			})
		}
		return fitting{
			clusters: c,
			required: Required([]int{index}),
		}
	}
}

func Race() Provider {
	return func(clusters []Cluster) Fitting {
		return fitting{
			clusters: clusters,
			required: index(len(clusters)),
		}
	}
}

func Ignore() Provider {
	return func(clusters []Cluster) Fitting {
		return fitting{
			clusters: clusters[0:1],
			required: Required([]int{0}),
		}
	}
}

type fitting struct {
	clusters []Cluster
	required Required
}

func (s fitting) Clusters() []Cluster {
	return s.clusters
}

func (s fitting) Required() Required {
	return s.required
}

func index(v int) Required {
	var res []int
	for i := 0; i < v; i++ {
		res = append(res, i)
	}
	return Required(res)
}

func filter(clusters []Cluster, fn func(Cluster) bool) []Cluster {
	var res []Cluster
	for _, v := range clusters {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}
