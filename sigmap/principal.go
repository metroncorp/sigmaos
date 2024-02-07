package sigmap

func NewPrincipal(id TprincipalID, token Ttoken) *Tprincipal {
	return &Tprincipal{
		IDStr:     id.String(),
		TokenStr2: token.String(),
	}
}

func (p *Tprincipal) SetID(principalID TprincipalID) {
	p.IDStr = principalID.String()
}

func (p *Tprincipal) GetID() TprincipalID {
	return TprincipalID(p.IDStr)
}

func (p *Tprincipal) GetToken() Ttoken {
	return Ttoken(p.TokenStr2)
}

func (p *Tprincipal) SetToken(t Ttoken) {
	p.TokenStr2 = t.String()
}

func (t Ttoken) String() string {
	return string(t)
}

func (id TprincipalID) String() string {
	return string(id)
}
