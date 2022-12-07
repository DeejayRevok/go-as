package auth

type AccessTokenDeserializer interface {
	Deserialize(serializedToken string) (*AccessToken, error)
}
