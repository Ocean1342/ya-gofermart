package common

import "math"

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

func MoneyFloatToInt(in float64) int {
	return int(math.Round(in * 100))
}

func MoneyIntToFloat(in int) float64 {
	return float64(in / 100)
}
