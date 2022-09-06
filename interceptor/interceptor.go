// Package interceptor exports a interfaces that defines a interceptor.
package interceptor

import "net/http"

// InterceptorRuler defines a function that contains the rules to determine if a request must be injected or not.
//
// Receives an *http.Response with the requested content to be evaluated by the rule.
// Returns true if a request pass the rule, otherwise must be return false.
type InterceptorRuler func(*http.Response) bool

// InterceptorHandler defines a function that execution the injection.
//
// Receives an *http.Response with the request content to modify.
type InterceptorHandler func(*http.Response) error

// Interceptor exposes two methods to handle the modification that must be apply to the content requested.
type Interceptor interface {

	// Rules returns a list of interceptor.InterceptorRuler with the functions that it will validate
	// if the content must be modify or not.
	Rules() []InterceptorRuler

	// Handler returns a interceptor.InterceptorHandler with all operations to modify a content requested.
	Handler() InterceptorHandler
}
