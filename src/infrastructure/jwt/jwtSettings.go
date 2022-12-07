package jwt

import "crypto/rsa"

type JWTSettings struct {
	PublicKey *rsa.PublicKey
}

func NewJWTSettings(publickey *rsa.PublicKey) *JWTSettings {
	settings := JWTSettings{
		PublicKey: publickey,
	}
	return &settings
}
