package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"sigmaos/proc"
	sp "sigmaos/sigmap"
)

const (
	ISSUER = "sigmaos"
)

var ALL_PATHS []string = []string{"*"}

type ProcClaims struct {
	PrincipalID sp.TprincipalID `json:"principal_id"`
	//	PrincipalID  string   `json:"principal_id"`
	AllowedPaths []string           `json:"allowed_paths"`
	Secrets      map[string]*Secret `json:"secrets"`
	jwt.StandardClaims
}

// Construct proc claims from a proc env
func NewProcClaims(pe *proc.ProcEnv) *ProcClaims {
	secrets := make(map[string]*Secret)
	for svc, s := range pe.GetClaims().GetSecrets() {
		secrets[svc] = NewSecret(s.ID, s.Key)
	}
	return &ProcClaims{
		PrincipalID:  pe.GetClaims().GetPrincipalID(),
		AllowedPaths: pe.GetClaims().GetAllowedPaths(),
		Secrets:      secrets,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 999).Unix(), // TODO: set expiry properly
			Issuer:    ISSUER,
		},
	}
}

func (pc *ProcClaims) GetSecrets() map[string]*Secret {
	return pc.Secrets
}

// Add a secret to a proc claim
func (pc *ProcClaims) AddSecret(svc string, s *Secret) {
	pc.Secrets[svc] = s
}

func (pc *ProcClaims) String() string {
	return fmt.Sprintf("&{ PrincipalID:%v AllowedPaths:%v Secrets:%v }", pc.PrincipalID, pc.AllowedPaths, pc.Secrets)
}
