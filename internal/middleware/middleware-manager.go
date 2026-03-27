package middleware

type MiddlewareManager struct {
	jwtkey []byte
}

func NewMiddlewareManager(jwtkey []byte) *MiddlewareManager {
	return &MiddlewareManager{
		jwtkey: jwtkey,
	}
}
