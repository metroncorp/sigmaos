// binsrv daemon. it takes as argument a local directory where to
// cache binaries.
package main

import (
	"os"

	db "sigmaos/debug"
	"sigmaos/uprocsrv/binsrv"
)

func main() {
	if len(os.Args) < 4 {
		db.DFatalf("%s: Usage <kernelid> <uprocpid> <mnt>\n", os.Args[0])
	}
	if err := binsrv.RunBinFS(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		db.DFatalf("RunBinFs %v err %v\n", os.Args, err)
	}
}
