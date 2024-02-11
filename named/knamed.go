package named

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"

	"sigmaos/auth"
	db "sigmaos/debug"
	"sigmaos/keys"
	"sigmaos/perf"
	"sigmaos/proc"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
)

func RunKNamed(args []string) error {
	pe := proc.GetProcEnv()
	db.DPrintf(db.NAMED, "%v: knamed %v\n", pe.GetPID(), args)
	if len(args) != 4 {
		return fmt.Errorf("%v: wrong number of arguments %v", args[0], args)
	}
	masterKey := auth.SymmetricKey(args[3])
	// Self-sign token for bootstrapping purposes
	kmgr := keys.NewSymmetricKeyMgr(keys.WithConstGetKeyFn(masterKey))
	kmgr.AddKey(sp.Tsigner(pe.GetPID()), masterKey)
	kmgr.AddKey(auth.SIGMA_DEPLOYMENT_MASTER_SIGNER, masterKey)
	as, err1 := auth.NewAuthSrv[*jwt.SigningMethodHMAC](jwt.SigningMethodHS256, sp.Tsigner(pe.GetPID()), proc.NOT_SET, kmgr)
	if err1 != nil {
		db.DPrintf(db.ERROR, "Error bootstrapping auth srv: %v", err1)
		return err1
	}
	pc := auth.NewProcClaims(pe)
	token, err1 := as.NewToken(pc)
	if err1 != nil {
		db.DPrintf(db.ERROR, "Error NewToken: %v", err1)
		return err1
	}
	pe.SetToken(token)

	nd := &Named{}
	nd.realm = sp.Trealm(args[1])

	p, err := perf.NewPerf(pe, perf.KNAMED)
	if err != nil {
		db.DFatalf("Error NewPerf: %v", err)
	}
	defer p.Done()

	sc, err := sigmaclnt.NewSigmaClntFsLib(pe)
	if err != nil {
		db.DFatalf("NewSigmaClntFsLib: err %v", err)
	}
	nd.SigmaClnt = sc

	init := args[2]

	nd.masterKey = masterKey

	db.DPrintf(db.NAMED, "started %v %v", pe.GetPID(), nd.realm)

	w := os.NewFile(uintptr(3), "pipew")
	r := os.NewFile(uintptr(4), "piper")
	w2 := os.NewFile(uintptr(5), "pipew")
	w2.Close()

	if init == "start" {
		fmt.Fprintf(w, init)
		w.Close()
	}

	if err := nd.startLeader(); err != nil {
		db.DFatalf("Error startLeader %v\n", err)
	}
	defer nd.fs.Close()

	mnt, err := nd.newSrv()
	if err != nil {
		db.DFatalf("Error newSrv %v\n", err)
	}

	db.DPrintf(db.NAMED, "newSrv %v mnt %v", nd.realm, mnt)

	if err := nd.fs.SetRootNamed(mnt); err != nil {
		db.DFatalf("SetNamed: %v", err)
	}

	if init == "init" {
		nd.initfs()
		fmt.Fprintf(w, init)
		w.Close()
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		db.DPrintf(db.ALWAYS, "pipe read err %v", err)
		return err
	}
	r.Close()

	db.DPrintf(db.NAMED, "%v: knamed done %v %v %v\n", pe.GetPID(), nd.realm, mnt, string(data))

	nd.resign()

	return nil
}

var InitRootDir = []string{sp.BOOT, sp.KPIDS, sp.MEMFS, sp.LCSCHED, sp.PROCQ, sp.SCHEDD, sp.UX, sp.S3, sp.DB, sp.MONGO, sp.REALM, sp.KEYS}

// If initial root dir doesn't exist, create it.
func (nd *Named) initfs() error {
	for _, n := range InitRootDir {
		if _, err := nd.SigmaClnt.Create(n, 0777|sp.DMDIR, sp.OREAD); err != nil {
			db.DPrintf(db.ALWAYS, "Error create [%v]: %v", n, err)
			return err
		}
	}
	return nil
}
