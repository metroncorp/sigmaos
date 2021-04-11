package kv

//
// Shard coordinator: assigns shards to KVs.  Assumes no KV failures
// This is a short-lived daemon: it rebalances shards and then exists.
//

import (
	"errors"
	"fmt"
	"log"
	"sync"

	db "ulambda/debug"
	"ulambda/fsclnt"
	"ulambda/fslib"
	"ulambda/memfsd"
)

const (
	NSHARD       = 10
	KVDIR        = "name/kv"
	SHARDER      = KVDIR + "/sharder"
	KVCONFIG     = KVDIR + "/config"
	KVNEXTCONFIG = KVDIR + "/nextconfig"
	KVCOMMIT     = KVDIR + "/commit/"
)

var ErrWrongKv = errors.New("ErrWrongKv")
var ErrRetry = errors.New("ErrRetry")

type Config struct {
	N      int
	Shards []string // maps shard # to server
}

func makeConfig(n int) *Config {
	cf := &Config{n, make([]string, NSHARD)}
	return cf
}

type Sharder struct {
	mu sync.Mutex
	*fslib.FsLibSrv
	ch       chan bool
	pid      string
	args     []string
	kvs      []string // the kv servers in this configuration
	nextKvs  []string // the available kvs for the next config
	nkvd     int      // # KVs in reconfiguration
	conf     *Config
	nextConf *Config
	done     bool
}

func MakeSharder(args []string) (*Sharder, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("MakeSharder: too few arguments %v\n", args)
	}
	sh := &Sharder{}
	sh.ch = make(chan bool)
	sh.pid = args[0]
	sh.args = args[1:]
	db.Name("sharder")
	ip, err := fsclnt.LocalIP()
	if err != nil {
		return nil, fmt.Errorf("MakeSharder: no IP %v\n", err)
	}
	fsd := memfsd.MakeFsd(ip + ":0")
	db.DLPrintf("SHARDER", "New sharder %v", args)
	fls, err := fslib.InitFs(SHARDER, fsd, nil)
	if err != nil {
		return nil, err
	}
	sh.FsLibSrv = fls
	sh.Started(sh.pid)

	err = sh.Mkdir(KVCOMMIT, 0777)
	if err != nil {
		db.DLPrintf("SHARDER", "MkDir %v failed %v\n", KVCOMMIT, err)
	}

	return sh, nil
}

// Caller holds lock
// XXX minimize movement
func (sh *Sharder) balance() *Config {
	j := 0
	conf := makeConfig(sh.conf.N + 1)

	db.DLPrintf("SHARDER", "shards %v (len %v) kvs %v\n", sh.conf.Shards,
		len(sh.conf.Shards), sh.nextKvs)

	if len(sh.nextKvs) == 0 {
		return conf
	}
	for i, _ := range sh.conf.Shards {
		conf.Shards[i] = sh.nextKvs[j]
		j = (j + 1) % len(sh.nextKvs)
	}
	return conf
}

func (sh *Sharder) Exit() {
	sh.ExitFs(SHARDER)
	sh.Exiting(sh.pid, "OK")
}

func (sh *Sharder) readConfig(conffile string) *Config {
	conf := Config{}
	err := sh.ReadFileJson(conffile, &conf)
	if err != nil {
		return nil
	}
	sh.kvs = make([]string, 0)
	for _, kv := range conf.Shards {
		present := false
		if kv == "" {
			continue
		}
		for _, k := range sh.kvs {
			if k == kv {
				present = true
				break
			}
		}
		if !present {
			sh.kvs = append(sh.kvs, kv)
		}
	}
	return &conf
}

func (sh *Sharder) Init() {
	sh.conf = makeConfig(0)
	err := sh.MakeFileJson(KVCONFIG, *sh.conf)
	if err != nil {
		log.Fatalf("Sharder: cannot make file  %v %v\n", KVCONFIG, err)
	}
}

func (sh *Sharder) watchPrepared(p string) {
	sh.ch <- true
}

func (sh *Sharder) Work() {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.conf = sh.readConfig(KVCONFIG)
	if sh.conf == nil {
		sh.Init()
	}

	db.DLPrintf("SHARDER", "Sharder: %v %v\n", sh.conf, sh.args)
	if sh.args[0] == "add" {
		sh.nextKvs = append(sh.kvs, sh.args[1:]...)
	} else {
		sh.nextKvs = make([]string, len(sh.kvs))
		copy(sh.nextKvs, sh.kvs)
		for _, del := range sh.args[1:] {
			for i, kv := range sh.nextKvs {
				if del == kv {
					sh.nextKvs = append(sh.nextKvs[:i],
						sh.nextKvs[i+1:]...)
				}
			}
		}
	}

	sh.nextConf = sh.balance()
	db.DLPrintf("SHARDER", "Sharder next conf: %v %v\n", sh.nextConf, sh.nextKvs)

	sts, err := sh.ReadDir(KVCOMMIT)
	if err != nil {
		log.Fatalf("SHARDER: ReadDir commit error %v\n", err)
	}

	for _, st := range sts {
		fn := KVCOMMIT + st.Name
		err = sh.Remove(fn)
		if err != nil {
			db.DLPrintf("SHARDER", "Remove %v failed %v\n", fn, err)
		}
	}

	if sh.args[0] == "del" {
		sh.nextKvs = append(sh.nextKvs, sh.args[1:]...)

	}

	sh.nkvd = len(sh.nextKvs)
	for _, kv := range sh.nextKvs {
		fn := KVCOMMIT + kv
		_, err := sh.ReadFileWatch(fn, sh.watchPrepared)
		if err == nil {
			log.Fatalf("SHARDER: watch failed %v", err)
		}
	}

	err = sh.MakeFileJson(KVNEXTCONFIG, *sh.nextConf)
	if err != nil {
		log.Printf("SHARDER: %v error %v\n", KVNEXTCONFIG, err)
		return
	}

	for i := 0; i < sh.nkvd; i++ {
		<-sh.ch
	}

	db.DLPrintf("SHARDER", "Commit to %v\n", sh.nextConf)
	// commit to new config
	err = sh.Rename(KVNEXTCONFIG, KVCONFIG)
	if err != nil {
		db.DLPrintf("SHARDER", "SHARDER: rename %v -> %v: error %v\n",
			KVNEXTCONFIG, KVCONFIG, err)
		return
	}
}
