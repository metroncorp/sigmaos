package sortedmap

import (
	"fmt"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"

	"sigmaos/rand"
)

type SortedMap[K constraints.Ordered, V any] struct {
	sync.Mutex
	dents  map[K]V
	sorted []K
	rr     int
}

func NewSortedMap[K constraints.Ordered, V any]() *SortedMap[K, V] {
	sd := &SortedMap[K, V]{}
	sd.dents = make(map[K]V)
	sd.sorted = make([]K, 0)
	return sd
}

func (sd *SortedMap[K, V]) roundrobinIndexL() int {
	if len(sd.sorted) == 0 {
		return 0
	}
	return sd.rr % len(sd.sorted)
}

func (sd *SortedMap[K, V]) roundrobinNextL() int {
	i := sd.roundrobinIndexL()
	sd.rr += 1
	return i
}

func (sd *SortedMap[K, V]) Iter(f func(key K, val V) bool) {
	sd.Lock()
	defer sd.Unlock()

	for _, k := range sd.sorted {
		b := f(k, sd.dents[k])
		if !b {
			return
		}
	}
}

func (sd *SortedMap[K, V]) String() string {
	s := "["
	sd.Iter(func(k K, v V) bool {
		s += fmt.Sprintf("%v , ", k)
		return true
	})
	return s + "]"
}

func (sd *SortedMap[K, V]) Len() int {
	sd.Lock()
	defer sd.Unlock()

	return len(sd.dents)
}

func (sd *SortedMap[K, V]) Lookup(n K) (V, bool) {
	sd.Lock()
	defer sd.Unlock()

	e, ok := sd.dents[n]
	return e, ok
}

func (sd *SortedMap[K, V]) Random() (K, bool) {
	sd.Lock()
	defer sd.Unlock()

	if len(sd.sorted) == 0 {
		var k K
		return k, false
	}
	k := sd.sorted[rand.Int64(int64(len(sd.sorted)))]
	return k, true
}

func (sd *SortedMap[K, V]) RoundRobin() (K, bool) {
	sd.Lock()
	defer sd.Unlock()

	if len(sd.sorted) == 0 {
		var k K
		return k, false
	}
	k := sd.sorted[sd.roundrobinNextL()]
	return k, true
}

func (sd *SortedMap[K, V]) Keys(s int) []K {
	sd.Lock()
	defer sd.Unlock()

	keys := make([]K, len(sd.sorted[s:]))
	copy(keys, sd.sorted[s:])
	return keys
}

func (sd *SortedMap[K, V]) insertSortL(name K) {
	i := sort.Search(len(sd.sorted), func(i int) bool { return sd.sorted[i] >= name })
	new := make([]K, 1)
	sd.sorted = append(sd.sorted, new...)
	copy(sd.sorted[i+1:], sd.sorted[i:])
	sd.sorted[i] = name
	if i < sd.roundrobinIndexL() {
		sd.rr += 1
	}
}

func (sd *SortedMap[K, V]) delSortL(name K) {
	i := sort.Search(len(sd.sorted), func(i int) bool { return sd.sorted[i] >= name })
	sd.sorted = append(sd.sorted[:i], sd.sorted[i+1:]...)
	if i < sd.roundrobinIndexL() {
		if sd.rr > 0 {
			sd.rr -= 1
		}
	}
}

// Return true if K was inserted
func (sd *SortedMap[K, V]) Insert(name K, e V) bool {
	sd.Lock()
	defer sd.Unlock()
	if _, ok := sd.dents[name]; !ok {
		sd.dents[name] = e
		sd.insertSortL(name)
		return true
	}
	return false
}

func (sd *SortedMap[K, V]) Delete(name K) bool {
	sd.Lock()
	defer sd.Unlock()

	if _, ok := sd.dents[name]; !ok {
		return false
	}
	delete(sd.dents, name)
	sd.delSortL(name)
	return true
}
