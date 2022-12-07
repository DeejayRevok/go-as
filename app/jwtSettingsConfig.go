package app

import (
	"go-as/src/infrastructure/iam"
	"go-as/src/infrastructure/jwt"
	"os"

	"go.uber.org/zap"
)

func LoadJWTSettings(logger *zap.Logger) *jwt.JWTSettings {
	iamBasePath := os.Getenv("IAM_BASE_PATH")
	rsaPublicKeyBuilder := iam.NewIAMJWTRSAPublicKeyBuilder(iamBasePath)
	publicKey, err := rsaPublicKeyBuilder.Build()
	if err != nil {
		logger.Fatal("Error loading jwt rsa public key")
	}
	return jwt.NewJWTSettings(publicKey)
}
