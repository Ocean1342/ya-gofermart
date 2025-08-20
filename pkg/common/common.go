package common

var CtxUser User

func init() {
	CtxUser = User{}
}

type User struct {
}

const (
	XRequestID              = "X-Request-ID"
	AuthorizationHeaderName = "Authorization"
)
