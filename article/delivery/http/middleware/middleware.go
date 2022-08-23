package middleware

import "github.com/labstack/echo"

type Middleware struct {
}

func (m *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		context.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(context)
	}
}

func InitMiddleware() *Middleware {
	return &Middleware{}
}
