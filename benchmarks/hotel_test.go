package benchmarks_test

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/stretchr/testify/assert"

	"sigmaos/apps/hotel"
	"sigmaos/benchmarks/loadgen"
	db "sigmaos/debug"
	"sigmaos/proc"
	sp "sigmaos/sigmap"
	"sigmaos/test"
	"sigmaos/util/perf"
	rd "sigmaos/util/rand"
)

const (
	RAND_INIT = 12345
)

type hotelFn func(wc *hotel.WebClnt, r *rand.Rand)

type HotelJobInstance struct {
	sigmaos             bool
	justCli             bool
	k8ssrvaddr          string
	job                 string
	dur                 []time.Duration
	maxrps              []int
	ncache              int
	cachetype           string
	scaleCacheDelay     time.Duration
	manuallyScaleCaches bool
	nCachesToAdd        int
	scaleGeoDelay       time.Duration
	manuallyScaleGeo    bool
	nGeoToAdd           int
	nGeo                int
	geoNIdx             int
	geoSearchRadius     int
	geoNResults         int
	ready               chan bool
	fn                  hotelFn
	hj                  *hotel.HotelJob
	lgs                 []*loadgen.LoadGenerator
	p                   *perf.Perf
	wc                  *hotel.WebClnt
	*test.RealmTstate
}

func NewHotelJob(ts *test.RealmTstate, p *perf.Perf, sigmaos bool, durs string, maxrpss string, fn hotelFn, justCli bool, ncache int, cachetype string, cacheMcpu proc.Tmcpu, manuallyScaleCaches bool, scaleCacheDelay time.Duration, nCachesToAdd int, nGeo int, geoNIndex int, geoSearchRadius int, geoNResults int, manuallyScaleGeo bool, scaleGeoDelay time.Duration, nGeoToAdd int) *HotelJobInstance {
	ji := &HotelJobInstance{}
	ji.sigmaos = sigmaos
	ji.job = rd.String(8)
	ji.ready = make(chan bool)
	ji.fn = fn
	ji.RealmTstate = ts
	ji.p = p
	ji.justCli = justCli
	ji.ncache = ncache
	ji.cachetype = cachetype
	ji.manuallyScaleCaches = manuallyScaleCaches
	ji.scaleCacheDelay = scaleCacheDelay
	ji.nCachesToAdd = nCachesToAdd
	ji.manuallyScaleGeo = manuallyScaleGeo
	ji.scaleGeoDelay = scaleGeoDelay
	ji.nGeoToAdd = nGeoToAdd
	ji.nGeo = nGeo
	ji.geoNIdx = geoNIndex
	ji.geoSearchRadius = geoSearchRadius
	ji.geoNResults = geoNResults

	durslice := strings.Split(durs, ",")
	maxrpsslice := strings.Split(maxrpss, ",")
	assert.Equal(ts.Ts.T, len(durslice), len(maxrpsslice), "Non-matching lengths: %v %v", durs, maxrpss)

	ji.dur = make([]time.Duration, 0, len(durslice))
	ji.maxrps = make([]int, 0, len(durslice))

	for i := range durslice {
		d, err := time.ParseDuration(durslice[i])
		assert.Nil(ts.Ts.T, err, "Bad duration %v", err)
		n, err := strconv.Atoi(maxrpsslice[i])
		assert.Nil(ts.Ts.T, err, "Bad duration %v", err)
		ji.dur = append(ji.dur, d)
		ji.maxrps = append(ji.maxrps, n)
	}

	var err error
	var svcs []*hotel.Srv
	if sigmaos {
		svcs = hotel.NewHotelSvc()
	}

	if ji.justCli {
		// Read job name
		sts, err := ji.GetDir("name/hotel/")
		assert.Nil(ji.Ts.T, err, "Err Get hotel dir %v", err)
		var l int
		for _, st := range sts {
			// Dumb heuristic, but will always be the longest name....
			if len(st.Name) > l {
				ji.job = st.Name
				l = len(st.Name)
			}
		}
	}

	if !ji.justCli {
		var nc = ncache
		// Only start one cache if autoscaling.
		if sigmaos && CACHE_TYPE == "cached" && HOTEL_CACHE_AUTOSCALE {
			nc = 1
		}
		if !sigmaos {
			nc = 0
		}
		ji.hj, err = hotel.NewHotelJob(ts.SigmaClnt, ji.job, svcs, N_HOTEL, cachetype, cacheMcpu, nc, CACHE_GC, HOTEL_IMG_SZ_MB, nGeo, geoNIndex, geoSearchRadius, geoNResults)
		assert.Nil(ts.Ts.T, err, "Error NewHotelJob: %v", err)
	}

	if !sigmaos {
		ji.k8ssrvaddr = K8S_ADDR
		// Write a file for clients to discover the server's address.
		if !ji.justCli {
			pn := hotel.JobHTTPAddrsPath(ji.job)
			h, p, err := net.SplitHostPort(K8S_ADDR)
			assert.Nil(ts.Ts.T, err, "Err split host port %v: %v", ji.k8ssrvaddr, err)
			port, err := strconv.Atoi(p)
			assert.Nil(ts.Ts.T, err, "Err parse port %v: %v", p, err)
			addr := sp.NewTaddr(sp.Tip(h), sp.Tport(port))
			mnt := sp.NewEndpoint(sp.EXTERNAL_EP, []*sp.Taddr{addr})
			if err = ts.MkEndpointFile(pn, mnt); err != nil {
				db.DFatalf("MkEndpointFile mnt %v err %v", mnt, err)
			}
		}
	}

	if sigmaos {
		if HOTEL_CACHE_AUTOSCALE && cachetype == "cached" && !ji.justCli {
			ji.hj.CacheAutoscaler.Run(1*time.Second, ncache)
		}
	}

	wc, err := hotel.NewWebClnt(ts.FsLib, ji.job)
	assert.Nil(ts.Ts.T, err, "Err NewWebClnt: %v", err)
	ji.wc = wc
	// Make a load generators.
	ji.lgs = make([]*loadgen.LoadGenerator, 0, len(ji.dur))
	for i := range ji.dur {
		ji.lgs = append(ji.lgs, loadgen.NewLoadGenerator(ji.dur[i], ji.maxrps[i], func(r *rand.Rand) (time.Duration, bool) {
			// Run a single request.
			ji.fn(ji.wc, r)
			return 0, false
		}))
	}
	return ji
}

func (ji *HotelJobInstance) StartHotelJob() {
	db.DPrintf(db.ALWAYS, "StartHotelJob dur %v ncache %v maxrps %v kubernetes (%v,%v) manuallyScaleCaches %v scaleCacheDelay %v nCachesToAdd %v manuallyScaleGeo %v scaleGeoDelay %v nGeoToAdd %v nGeoInit %v geoNIndex %v geoSearchRadius: %v geoNResults: %v", ji.dur, ji.ncache, ji.maxrps, !ji.sigmaos, ji.k8ssrvaddr, ji.manuallyScaleCaches, ji.scaleCacheDelay, ji.nCachesToAdd, ji.manuallyScaleGeo, ji.scaleGeoDelay, ji.nGeoToAdd, ji.nGeo, ji.geoNIdx, ji.geoSearchRadius, ji.geoNResults)
	var wg sync.WaitGroup
	for _, lg := range ji.lgs {
		wg.Add(1)
		go func(lg *loadgen.LoadGenerator, wg *sync.WaitGroup) {
			defer wg.Done()
			lg.Calibrate()
		}(lg, &wg)
	}
	wg.Wait()
	_, err := ji.wc.StartRecording()
	if err != nil {
		db.DFatalf("Can't start recording: %v", err)
	}
	if !ji.justCli && ji.manuallyScaleGeo {
		go func() {
			time.Sleep(ji.scaleGeoDelay)
			if ji.sigmaos {
				for i := 0; i < ji.nGeoToAdd; i++ {
					err := ji.hj.AddGeoSrv()
					assert.Nil(ji.Ts.T, err, "Add Geo srv: %v", err)
				}
			} else {
				if ji.nGeoToAdd > 0 {
					err := k8sScaleUpGeo()
					assert.Nil(ji.Ts.T, err, "K8s scale up Geo srv: %v", err)
				} else {
					db.DPrintf(db.ALWAYS, "No geos meant to be added. Skip scaling up")
				}
			}
		}()
	}
	if !ji.justCli && ji.manuallyScaleCaches {
		go func() {
			time.Sleep(ji.scaleCacheDelay)
			ji.hj.CacheAutoscaler.AddServers(ji.nCachesToAdd)
		}()
	}
	for i, lg := range ji.lgs {
		db.DPrintf(db.TEST, "Run load generator rps %v dur %v", ji.maxrps[i], ji.dur[i])
		lg.Run()
		//    ji.printStats()
	}
	db.DPrintf(db.ALWAYS, "Done running HotelJob")
}

func (ji *HotelJobInstance) printStats() {
	if ji.sigmaos && !ji.justCli {
		for _, s := range hotel.HOTELSVC {
			stats, err := ji.ReadStats(s)
			assert.Nil(ji.Ts.T, err, "error get stats [%v] %v", s, err)
			fmt.Printf("= %s: %v\n", s, stats)
		}
		cs, err := ji.hj.StatsSrv()
		assert.Nil(ji.Ts.T, err)
		for i, cstat := range cs {
			fmt.Printf("= cache-%v: %v\n", i, cstat)
		}
	}
}

func (ji *HotelJobInstance) Wait() {
	db.DPrintf(db.TEST, "extra sleep")
	time.Sleep(20 * time.Second)
	if ji.p != nil {
		ji.p.Done()
	}
	db.DPrintf(db.TEST, "Evicting hotel procs")
	if ji.sigmaos && !ji.justCli {
		ji.printStats()
		err := ji.hj.Stop()
		assert.Nil(ji.Ts.T, err, "stop %v", err)
	}
	db.DPrintf(db.TEST, "Done evicting hotel procs")
	for _, lg := range ji.lgs {
		db.DPrintf(db.ALWAYS, "Data:\n%v", lg.StatsDataString())
	}
	for _, lg := range ji.lgs {
		lg.Stats()
	}
}

func (ji *HotelJobInstance) requestK8sStats() {
	rep, err := ji.wc.SaveResults()
	assert.Nil(ji.Ts.T, err, "Save results: %v", err)
	assert.Equal(ji.Ts.T, rep, "Done!", "Save results not ok: %v", rep)
}
