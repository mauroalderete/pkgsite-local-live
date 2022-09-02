package interceptor

import "net/http"

type InterceptorRuler func(*http.Response) bool
type InterceptorHandler func(*http.Response) error

type Interceptor interface {
	Rules() []InterceptorRuler
	Handler() InterceptorHandler
}
