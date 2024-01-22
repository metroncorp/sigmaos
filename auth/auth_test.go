package auth_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	db "sigmaos/debug"
	"sigmaos/fslib"
	"sigmaos/proc"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
	"sigmaos/test"
)

const (
	REALM1 sp.Trealm = "testrealm1"
)

func TestStartStop(t *testing.T) {
	rootts := test.NewTstateWithRealms(t)
	db.DPrintf(db.TEST, "Started successfully")
	rootts.Shutdown()
}

func TestOK(t *testing.T) {
	rootts := test.NewTstateWithRealms(t)

	sts, err := rootts.GetDir(sp.NAMED)
	assert.Nil(t, err)

	db.DPrintf(db.TEST, "realm named root %v", sp.Names(sts))

	assert.True(t, fslib.Present(sts, []string{sp.UXREL}), "initfs")

	sts, err = rootts.GetDir(sp.SCHEDD)
	assert.Nil(t, err)

	db.DPrintf(db.TEST, "realm names sched %v", sp.Names(sts))

	sts2, err := rootts.GetDir(path.Join(sp.SCHEDD, sts[0].Name) + "/")
	assert.Nil(t, err, "Err getdir: %v", err)

	db.DPrintf(db.TEST, "sched contents %v", sp.Names(sts2))

	rootts.Shutdown()
}

func TestMaliciousPrincipalFail(t *testing.T) {
	rootts := test.NewTstateWithRealms(t)

	// Create a new sigma clnt, with an unexpected principal
	pe := proc.NewAddedProcEnv(rootts.ProcEnv(), 1)
	pe.SetPrincipal(&sp.Tprincipal{
		ID:           "malicious-user",
		TokenPresent: false,
	})
	sc1, err := sigmaclnt.NewSigmaClnt(pe)
	assert.Nil(t, err, "Err NewClnt: %v", err)

	_, err = sc1.GetDir(sp.NAMED)
	assert.NotNil(t, err)

	sts, err := rootts.GetDir(sp.SCHEDD)
	assert.Nil(t, err)

	db.DPrintf(db.TEST, "realm names sched %v", sp.Names(sts))

	_, err = sc1.GetDir(path.Join(sp.SCHEDD, sts[0].Name) + "/")
	assert.NotNil(t, err)

	rootts.Shutdown()
}

func TestNoDelegationPrincipalFail(t *testing.T) {
	rootts := test.NewTstateWithRealms(t)

	// Create a new sigma clnt, with an unexpected principal
	pe := proc.NewAddedProcEnv(rootts.ProcEnv(), 1)
	pe.SetPrincipal(&sp.Tprincipal{
		ID:           "malicious-user",
		TokenPresent: false,
	})
	sc1, err := sigmaclnt.NewSigmaClnt(pe)
	assert.Nil(t, err, "Err NewClnt: %v", err)

	_, err = sc1.GetDir(sp.NAMED)
	assert.NotNil(t, err)

	sts, err := rootts.GetDir(sp.SCHEDD)
	assert.Nil(t, err)

	db.DPrintf(db.TEST, "realm names sched %v", sp.Names(sts))

	_, err = sc1.GetDir(path.Join(sp.SCHEDD, sts[0].Name) + "/")
	assert.NotNil(t, err)

	rootts.Shutdown()
}
