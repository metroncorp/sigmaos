package crash

import (
	"os"
	"strconv"
	"time"

	db "ulambda/debug"
	"ulambda/fslib"
	"ulambda/rand"
)

//
// Crash/partition testing
//

func GetEnv() int64 {
	crash := os.Getenv("SIGMACRASH")
	n, err := strconv.Atoi(crash)
	if err != nil {
		n = 0
	}
	return int64(n)
}

func Crasher(fsl *fslib.FsLib) {
	crash := GetEnv()
	if crash == 0 {
		return
	}
	go func() {
		for true {
			ms := rand.Int64(crash)
			// log.Printf("%v: ms %v\n", proc.GetProgram(), ms)
			time.Sleep(time.Duration(ms) * time.Millisecond)
			r := rand.Int64(1000)
			// log.Printf("%v: r = %v\n", proc.GetProgram(), r)
			if r < 330 {
				Crash(fsl)
			} else if r < 660 {
				Partition(fsl)
			}
		}
	}()
}

func Crash(fsl *fslib.FsLib) {
	db.DLPrintf(db.ALWAYS, "crash.Crash %v\n", os.Args)
	os.Exit(1)
}

func Partition(fsl *fslib.FsLib) {
	db.DLPrintf(db.ALWAYS, "crash.Partition %v\n", os.Args)
	if err := fsl.Disconnect("name"); err != nil {
		db.DLPrintf(db.ALWAYS, "Disconnect %v name fails err %v\n", os.Args, err)
	}
	time.Sleep(time.Duration(5) * time.Millisecond)
}

func MaybePartition(fsl *fslib.FsLib) bool {
	r := rand.Int64(1000)
	if r < 330 {
		Partition(fsl)
		return true
	}
	return false
}
