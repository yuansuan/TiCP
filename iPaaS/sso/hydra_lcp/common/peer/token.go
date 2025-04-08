package peer

import (
	"github.com/dgrijalva/jwt-go"
)

// JWTClaims JWTClaims
type JWTClaims struct {
	ApplicantPeer     string `json:"applicant_peer"`
	ResourceOwnerPeer string `json:"resource_owner_peer"`

	ResourceKind string `json:"resource_kind"`
	Resource     string `json:"resource"`

	*jwt.StandardClaims
}
