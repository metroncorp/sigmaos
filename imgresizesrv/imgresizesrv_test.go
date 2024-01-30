package imgresizesrv_test

import (
	"fmt"
	"image/jpeg"
	"os"
	"path"
	"testing"
	"time"

	"github.com/nfnt/resize"
	"github.com/stretchr/testify/assert"

	db "sigmaos/debug"
	"sigmaos/fsetcd"
	"sigmaos/fttasks"
	"sigmaos/groupmgr"
	"sigmaos/imgresizesrv"
	"sigmaos/proc"
	rd "sigmaos/rand"
	sp "sigmaos/sigmap"
	"sigmaos/test"
)

const (
	IMG_RESIZE_MCPU proc.Tmcpu = 100
	IMG_RESIZE_MEM  proc.Tmem  = 0
)

func TestCompile(t *testing.T) {
}

func TestResizeImg(t *testing.T) {
	fn := "/tmp/thumb.jpeg"

	os.Remove(fn)

	in, err := os.Open("1.jpg")
	assert.Nil(t, err)
	img, err := jpeg.Decode(in)
	assert.Nil(t, err)

	start := time.Now()

	img1 := resize.Resize(160, 0, img, resize.Lanczos3)

	db.DPrintf(db.TEST, "resize %v\n", time.Since(start))

	out, err := os.Create(fn)
	assert.Nil(t, err)
	jpeg.Encode(out, img1, nil)
}

func TestResizeProc(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}
	in := path.Join(sp.S3, "~local/9ps3/img-save/6.jpg")
	//	in := path.Join(sp.S3, "~local/9ps3/img-save/6.jpg")
	out := path.Join(sp.S3, "~local/9ps3/img/6-thumb-xxx.jpg")
	ts.Remove(out)
	p := proc.NewProc("imgresize", []string{in, out, "1"})
	err := ts.Spawn(p)
	assert.Nil(t, err, "Spawn")
	err = ts.WaitStart(p.GetPid())
	assert.Nil(t, err, "WaitStart error")
	status, err := ts.WaitExit(p.GetPid())
	assert.Nil(t, err, "WaitExit error %v", err)
	assert.True(t, status.IsStatusOK(), "WaitExit status error: %v", status)
	ts.Shutdown()
}

type Tstate struct {
	job string
	*test.Tstate
	ch chan bool
	ft *fttasks.FtTasks
}

func newTstate(t *test.Tstate) *Tstate {
	ts := &Tstate{}
	ts.Tstate = t
	ts.job = rd.String(4)
	ts.ch = make(chan bool)
	ts.cleanup1()

	ft, err := fttasks.MkFtTasks(ts.SigmaClnt.FsLib, imgresizesrv.IMG, ts.job)
	assert.Nil(ts.T, err)
	ts.ft = ft
	return ts
}

func (ts *Tstate) cleanup1() {
	ts.RmDir(imgresizesrv.IMG)
	imgresizesrv.Cleanup(ts.FsLib, path.Join(sp.S3, "~local/9ps3/img-save"))
}

func (ts *Tstate) shutdown() {
	ts.ch <- true
	ts.Shutdown()
}

func (ts *Tstate) progress() {
	for true {
		select {
		case <-ts.ch:
			return
		case <-time.After(1 * time.Second):
			if n, err := ts.ft.NTaskDone(); err != nil {
				assert.Nil(ts.T, err)
			} else {
				fmt.Printf("%d..", n)
			}
		}
	}
}

func TestImgdFatal(t *testing.T) {
	t1, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}
	ts := newTstate(t1)

	imgd := imgresizesrv.StartImgd(ts.SigmaClnt, ts.job, IMG_RESIZE_MCPU, IMG_RESIZE_MEM, false, 1, 0)

	fn := path.Join(sp.S3, "~local/9ps3/img-save/", "yyy.jpg")

	err := ts.ft.SubmitTask(fn)
	assert.Nil(ts.T, err)

	err = ts.ft.SubmitTask(fttasks.STOP)
	assert.Nil(ts.T, err)

	gs := imgd.WaitGroup()
	for _, s := range gs {
		assert.True(ts.T, s.IsStatusFatal(), s)
	}
	db.DPrintf(db.TEST, "shutdown\n")
	ts.Shutdown()
}

func (ts *Tstate) imgdJob(paths []string) {
	imgd := imgresizesrv.StartImgd(ts.SigmaClnt, ts.job, IMG_RESIZE_MCPU, IMG_RESIZE_MEM, false, 1, 0)

	for _, pn := range paths {
		db.DPrintf(db.TEST, "submit %v\n", pn)
		err := ts.ft.SubmitTask(pn)
		assert.Nil(ts.T, err)
	}

	err := ts.ft.SubmitTask(fttasks.STOP)
	assert.Nil(ts.T, err)

	go ts.progress()

	gs := imgd.WaitGroup()
	for _, s := range gs {
		assert.True(ts.T, s.IsStatusOK(), s)
	}
}

func TestImgdOne(t *testing.T) {
	t1, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}
	ts := newTstate(t1)
	fn := path.Join(sp.S3, "~local/9ps3/img-save/1.jpg")
	ts.imgdJob([]string{fn})
	ts.shutdown()
}

func TestImgdMany(t *testing.T) {
	t1, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}
	ts := newTstate(t1)

	sts, err := ts.GetDir(path.Join(sp.S3, "~local/9ps3/img-save"))
	assert.Nil(t, err)

	paths := make([]string, 0, len(sts))
	for _, st := range sts {
		fn := path.Join(sp.S3, "~local/9ps3/img-save/", st.Name)
		paths = append(paths, fn)
	}

	ts.imgdJob(paths)
	ts.shutdown()
}

func TestImgdRestart(t *testing.T) {
	t1, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}
	ts := newTstate(t1)

	fn := path.Join(sp.S3, "~local/9ps3/img-save/1.jpg")

	err := ts.ft.SubmitTask(fn)
	assert.Nil(t, err)

	imgd := imgresizesrv.StartImgd(ts.SigmaClnt, ts.job, IMG_RESIZE_MCPU, IMG_RESIZE_MEM, true, 1, 0)

	time.Sleep(2 * time.Second)

	imgd.StopGroup()

	ts.Shutdown()

	time.Sleep(2 * fsetcd.LeaseTTL * time.Second)

	db.DPrintf(db.TEST, "Restart")

	t2, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}
	ts.Tstate = t2

	gms, err := groupmgr.Recover(ts.SigmaClnt)
	assert.Nil(ts.T, err, "Recover")
	assert.Equal(ts.T, 1, len(gms))

	err = ts.ft.SubmitTask(fttasks.STOP)
	assert.Nil(t, err)

	go ts.progress()

	gms[0].WaitGroup()

	ts.shutdown()
}
