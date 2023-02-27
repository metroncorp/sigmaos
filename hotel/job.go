package hotel

import (
	"path"
	"strconv"

	"sigmaos/cacheclnt"
	db "sigmaos/debug"
	"sigmaos/fslib"
	"sigmaos/proc"
	"sigmaos/protdev"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
	"sigmaos/test"
)

const (
	HOTEL      = "hotel"
	HOTELDIR   = "name/hotel/"
	MEMFS      = "memfs"
	HTTP_ADDRS = "http-addr"
)

func JobDir(job string) string {
	return path.Join(HOTELDIR, job)
}

func JobHTTPAddrsPath(job string) string {
	return path.Join(JobDir(job), HTTP_ADDRS)
}

func MemFsPath(job string) string {
	return path.Join(JobDir(job), MEMFS)
}

func GetJobHTTPAddrs(fsl *fslib.FsLib, job string) (sp.Taddrs, error) {
	mnt, err := fsl.ReadMount(JobHTTPAddrsPath(job))
	if err != nil {
		return nil, err
	}
	return mnt.Addr, err
}

func InitHotelFs(fsl *fslib.FsLib, jobname string) {
	fsl.MkDir(HOTELDIR, 0777)
	if err := fsl.MkDir(JobDir(jobname), 0777); err != nil {
		db.DFatalf("Mkdir %v err %v\n", JobDir(jobname), err)
	}
}

type Srv struct {
	Name   string
	Public bool
}

func MkHotelSvc(public bool) []Srv {
	return []Srv{Srv{"hotel-userd", public}, Srv{"hotel-rated", public},
		Srv{"hotel-geod", public}, Srv{"hotel-profd", public},
		Srv{"hotel-searchd", public}, Srv{"hotel-reserved", public},
		Srv{"hotel-recd", public}, Srv{"hotel-wwwd", public}}
}

var ncores = []int{0, 1,
	1, 1, 3,
	3, 0, 2}

var cacheNcore = 2

//var ncores = []int{0, 2,
//	2, 2, 3,
//	3, 0, 2}

type HotelJob struct {
	*sigmaclnt.SigmaClnt
	cacheClnt *cacheclnt.CacheClnt
	cacheMgr  *cacheclnt.CacheMgr
	pids      []proc.Tpid
	cache     string
}

func MakeHotelJob(sc *sigmaclnt.SigmaClnt, job string, srvs []Srv, cache string, ncache int) (*HotelJob, error) {
	var cc *cacheclnt.CacheClnt
	var cm *cacheclnt.CacheMgr
	var err error

	// Init fs.
	InitHotelFs(sc.FsLib, job)

	// Create a cache clnt.
	if ncache > 0 {
		cm, err = cacheclnt.MkCacheMgr(sc, job, ncache, proc.Tcore(cacheNcore), test.Overlays)
		if err != nil {
			db.DFatalf("Error MkCacheMgr %v", err)
			return nil, err
		}
		cc, err = cacheclnt.MkCacheClnt(sc.FsLib, job)
		if err != nil {
			db.DFatalf("Error cacheclnt %v", err)
			return nil, err
		}
	}

	pids := make([]proc.Tpid, 0, len(srvs))

	for i, srv := range srvs {
		p := proc.MakeProc(srv.Name, []string{job, strconv.FormatBool(srv.Public), cache})
		p.SetNcore(proc.Tcore(ncores[i]))
		if _, errs := sc.SpawnBurst([]*proc.Proc{p}); len(errs) > 0 {
			db.DFatalf("Error burst-spawnn proc %v: %v", p, errs)
			return nil, err
		}
		if err = sc.WaitStart(p.GetPid()); err != nil {
			db.DFatalf("Error spawn proc %v: %v", p, err)
			return nil, err
		}
		pids = append(pids, p.GetPid())
	}

	return &HotelJob{sc, cc, cm, pids, "cached"}, nil
}

func (hj *HotelJob) Stop() error {
	for _, pid := range hj.pids {
		if err := hj.Evict(pid); err != nil {
			return err
		}
		if _, err := hj.WaitExit(pid); err != nil {
			return err
		}
	}
	if hj.cacheMgr != nil {
		hj.cacheMgr.Stop()
	}
	return nil
}

func (hj *HotelJob) StatsSrv() ([]*protdev.Stats, error) {
	return hj.cacheClnt.StatsSrv()
}
