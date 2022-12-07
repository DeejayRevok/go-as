package iam

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
)

const jwtKeySetPath = "/jwks"

type IAMJWTRSAPublicKeyBuilder struct {
	iamBasePath string
}

func (builder *IAMJWTRSAPublicKeyBuilder) Build() (*rsa.PublicKey, error) {
	jwyKeySetFullPath := builder.iamBasePath + jwtKeySetPath

	response, err := http.Get(jwyKeySetFullPath)
	if err != nil {
		return nil, err
	}

	responseJson, err := builder.getJsonFromResponse(response)
	if err != nil {
		return nil, err
	}

	signingKey := builder.getSigningKey(responseJson)
	if signingKey == nil {
		return nil, errors.New("signing key not found")
	}
	E, err := builder.getEFromKey(signingKey)
	if err != nil {
		return nil, err
	}
	N, err := builder.getNFromKey(signingKey)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: N,
		E: E,
	}, nil
}

func (*IAMJWTRSAPublicKeyBuilder) getJsonFromResponse(response *http.Response) (map[string]interface{}, error) {
	var parsedJson map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&parsedJson)
	if err != nil {
		return nil, err
	}
	return parsedJson, nil
}

func (*IAMJWTRSAPublicKeyBuilder) getSigningKey(responseJson map[string]interface{}) map[string]interface{} {
	keys := responseJson["keys"].([]interface{})
	for _, key := range keys {
		keyJson := key.(map[string]interface{})
		if keyJson["use"].(string) == "sig" {
			return keyJson
		}
	}
	return nil
}

func (*IAMJWTRSAPublicKeyBuilder) getNFromKey(key map[string]interface{}) (*big.Int, error) {
	n := new(big.Int)
	n, ok := n.SetString(key["n"].(string), 10)
	if !ok {
		return nil, fmt.Errorf("error converting N %s to int", n)
	}
	return n, nil
}

func (*IAMJWTRSAPublicKeyBuilder) getEFromKey(key map[string]interface{}) (int, error) {
	return strconv.Atoi(key["e"].(string))
}

func NewIAMJWTRSAPublicKeyBuilder(iamBasePath string) *IAMJWTRSAPublicKeyBuilder {
	return &IAMJWTRSAPublicKeyBuilder{
		iamBasePath: iamBasePath,
	}
}
