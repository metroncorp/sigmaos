package fss3

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"sigmaos/auth"
	db "sigmaos/debug"
	"sigmaos/fs"
	"sigmaos/path"
	"sigmaos/perf"
	proc "sigmaos/proc"
	"sigmaos/serr"
	sp "sigmaos/sigmap"
	"sigmaos/sigmasrv"
)

var fss3 *Fss3

type Fss3 struct {
	*sigmasrv.SigmaSrv
	mu      sync.Mutex
	clients map[string]*s3.Client
}

func (fss3 *Fss3) getClient(ctx fs.CtxI) (*s3.Client, *serr.Err) {
	fss3.mu.Lock()
	defer fss3.mu.Unlock()

	var clnt *s3.Client
	var ok bool
	if clnt, ok = fss3.clients[ctx.Principal().ID]; ok {
		return clnt, nil
	}
	s3secrets, ok := ctx.Claims().GetSecrets()["s3"]
	// If this principal doesn't carry any s3 secrets, return EPERM
	if !ok {
		return nil, serr.NewErr(serr.TErrPerm, fmt.Errorf("Principal %v has no S3 secrets", ctx.Principal().ID))
	}
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			auth.NewAWSCredentialsProvider(s3secrets),
		),
	)
	if err != nil {
		db.DFatalf("Failed to load SDK configuration %v", err)
	}
	clnt = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	fss3.clients[ctx.Principal().ID] = clnt
	return clnt, nil
}

func RunFss3(buckets []string) {
	fss3 = &Fss3{
		clients: make(map[string]*s3.Client),
	}
	root := newDir("", path.Path{}, sp.DMDIR)
	pe := proc.GetProcEnv()
	addr := sp.NewTaddrAnyPort(sp.INNER_CONTAINER_IP, pe.GetNet())
	ssrv, err := sigmasrv.NewSigmaSrvRoot(root, sp.S3, addr, pe)
	if err != nil {
		db.DFatalf("Error NewSigmaSrv: %v", err)
	}
	p, err := perf.NewPerf(ssrv.MemFs.SigmaClnt().ProcEnv(), perf.S3)
	if err != nil {
		db.DFatalf("Error NewPerf: %v", err)
	}
	defer p.Done()

	fss3.SigmaSrv = ssrv
	ssrv.RunServer()
}
