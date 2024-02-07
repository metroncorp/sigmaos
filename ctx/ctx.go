package ctx

import (
	"sigmaos/auth"
	"sigmaos/clntcond"
	"sigmaos/fs"
	"sigmaos/sessp"
	sp "sigmaos/sigmap"
)

type Ctx struct {
	principal *sp.Tprincipal
	claims    *auth.ProcClaims
	sessid    sessp.Tsession
	clntid    sp.TclntId
	sct       *clntcond.ClntCondTable
	fencefs   fs.Dir
}

func NewCtx(principal *sp.Tprincipal, claims *auth.ProcClaims, sessid sessp.Tsession, clntid sp.TclntId, sct *clntcond.ClntCondTable, fencefs fs.Dir) *Ctx {
	return &Ctx{
		principal: principal,
		claims:    claims,
		sessid:    sessid,
		clntid:    clntid,
		sct:       sct,
		fencefs:   fencefs,
	}
}

func NewCtxNull() *Ctx {
	return NewCtx(sp.NoPrincipal(), nil, 0, sp.NoClntId, nil, nil)
}

func (ctx *Ctx) Principal() *sp.Tprincipal {
	return ctx.principal
}

func (ctx *Ctx) Claims() *auth.ProcClaims {
	return ctx.claims
}

func (ctx *Ctx) SessionId() sessp.Tsession {
	return ctx.sessid
}

func (ctx *Ctx) ClntId() sp.TclntId {
	return ctx.clntid
}

func (ctx *Ctx) ClntCondTable() *clntcond.ClntCondTable {
	return ctx.sct
}

func (ctx *Ctx) FenceFs() fs.Dir {
	return ctx.fencefs
}
