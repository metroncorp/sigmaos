package kv

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"math/big"
	"strconv"
	"strings"
	"time"

	db "ulambda/debug"
	"ulambda/fenceclnt1"
	"ulambda/fslib"
	"ulambda/group"
	np "ulambda/ninep"
	"ulambda/procclnt"
	"ulambda/reader"
	"ulambda/writer"
)

//
// Clerk for sharded kv service, which repeatedly reads/writes keys.
//

const (
	NKEYS  = 100
	WAITMS = 50
)

func key2shard(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	shard := int(h.Sum32() % NSHARD)
	return shard
}

func keyPath(kvd, shard string, k string) string {
	d := shardPath(kvd, shard)
	return d + "/" + k
}

func shard(shard int) string {
	// return strconv.Itoa(shard)
	return fmt.Sprintf("%03d", shard)
}

func Key(k uint64) string {
	return "key" + strconv.FormatUint(k, 16)
}

func nrand() uint64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Uint64()
	return x
}

type KvClerk struct {
	*fslib.FsLib
	*procclnt.ProcClnt
	fclnt *fenceclnt1.FenceClnt
	conf  *Config
}

func MakeClerk(name string, namedAddr []string) (*KvClerk, error) {
	kc := &KvClerk{}
	kc.FsLib = fslib.MakeFsLibAddr(name, namedAddr)
	kc.ProcClnt = procclnt.MakeProcClnt(kc.FsLib)
	kc.conf = &Config{}
	kc.fclnt = fenceclnt1.MakeLeaderFenceClnt(kc.FsLib, KVBALANCER)
	if err := kc.switchConfig(); err != nil {
		return nil, err
	}
	return kc, nil
}

// Read config, and retry if we have a stale group fence
func (kc KvClerk) switchConfig() error {
	for {
		err := kc.GetFileJsonWatch(KVCONFIG, kc.conf)
		if err != nil {
			db.DLPrintf("KVCLERK_ERR", "GetFileJsonWatch %v err %v\n", KVCONFIG, err)
			return err
		}
		db.DLPrintf("KVCLERK0", "Conf %v\n", kc.conf)
		kvs := makeKvs(kc.conf.Shards).mkKvs()
		dirs := make([]string, 0, len(kvs)+1)
		for _, kvd := range kvs {
			dirs = append(dirs, group.GRPDIR+"/"+kvd)
		}
		if err := kc.fclnt.FenceAtEpoch(kc.conf.Epoch, dirs); err != nil {
			if np.IsErrVersion(err) {
				db.DLPrintf("KVCLERK_ERR", "version mismatch; retry\n")
				time.Sleep(WAITMS * time.Millisecond)
				continue
			}

			db.DLPrintf("KVCLERK_ERR", "FenceAtEpoch %v failed %v\n", dirs, err)
			return err
		}
		break
	}
	return nil
}

// Try to fix err; if return is nil, retry.
func (kc *KvClerk) fixRetry(err error) error {

	// Shard dir hasn't been created yet (config 0) or hasn't moved
	// yet, so wait a bit, and retry.  XXX make sleep time
	// dynamic?

	if np.IsErrNotfound(err) && strings.HasPrefix(np.ErrPath(err), "shard") {
		db.DLPrintf("KVCLERK_ERR", "Wait for shard %v\n", np.ErrPath(err))
		time.Sleep(WAITMS * time.Millisecond)
		return nil
	}
	if np.IsErrStale(err) || np.IsErrUnreachable(err) {
		db.DLPrintf("KVCLERK_ERR", "fixRetry %v\n", err)
		return kc.switchConfig()
	}

	// if && strings.Contains(np.ErrPath(err), KVCONF) {
	//}
	return err
}

// Do an operation. If an error, try to fix the error (e.g., rereading
// config), and on success, retry.
func (kc *KvClerk) doop(o *op) {
	s := key2shard(o.k)
	for {
		db.DLPrintf("KVCLERK", "o %v conf %v\n", o.kind, kc.conf)
		fn := keyPath(kc.conf.Shards[s], shard(s), o.k)
		o.do(kc.FsLib, fn)
		if o.err == nil { // success?
			return
		}
		o.err = kc.fixRetry(o.err)
		if o.err != nil {
			return
		}
	}
}

type opT int

const (
	GETVAL opT = iota + 1
	PUT
	SET
	GETRD
)

type op struct {
	kind opT
	b    []byte
	k    string
	off  np.Toffset
	rdr  *reader.Reader
	err  error
}

func (o *op) do(fsl *fslib.FsLib, fn string) {
	switch o.kind {
	case GETVAL:
		o.b, o.err = fsl.GetFile(fn)
	case GETRD:
		o.rdr, o.err = fsl.OpenReader(fn)
	case PUT:
		_, o.err = fsl.PutFile(fn, 0777, np.OWRITE, o.b)
	case SET:
		_, o.err = fsl.SetFile(fn, o.b, o.off)
	}
	db.DLPrintf("KVCLERK", "op %v fn %v err %v\n", o.kind, fn, o.err)
}

func (kc *KvClerk) Get(k string, off np.Toffset) ([]byte, error) {
	op := &op{GETVAL, []byte{}, k, off, nil, nil}
	kc.doop(op)
	return op.b, op.err
}

func (kc *KvClerk) GetReader(k string) (*reader.Reader, error) {
	op := &op{GETRD, []byte{}, k, 0, nil, nil}
	kc.doop(op)
	return op.rdr, op.err
}

func (kc *KvClerk) Set(k string, b []byte, off np.Toffset) error {
	op := &op{SET, b, k, off, nil, nil}
	kc.doop(op)
	return op.err
}

func (kc *KvClerk) Append(k string, b []byte) error {
	op := &op{SET, b, k, np.NoOffset, nil, nil}
	kc.doop(op)
	return op.err
}

func (kc *KvClerk) Put(k string, b []byte) error {
	op := &op{PUT, b, k, 0, nil, nil}
	kc.doop(op)
	return op.err
}

func (kc *KvClerk) AppendJson(k string, v interface{}) error {
	b, err := writer.JsonRecord(v)
	if err != nil {
		return err
	}
	op := &op{SET, b, k, np.NoOffset, nil, nil}
	kc.doop(op)
	return op.err
}
