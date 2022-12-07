package auth

type AccessToken struct {
	Iss   string
	Sub   string
	Exp   int64
	Iat   int64
	Scope string
}
