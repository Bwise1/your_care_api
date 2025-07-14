package tracing

import (
	"github.com/bwise1/your_care_api/util/values"
	"github.com/lucsky/cuid"
)

type Context struct {
	// RequestID specifies the request ID if empty, a new request ID should be generated
	RequestID string
}

// New creates a new tracing context
func New() *Context {
	return &Context{
		RequestID: cuid.New(),
	}
}

// OutgoingHeaders returns the tracing information for response headers
func (tc *Context) OutgoingHeaders() map[string]string {
	return map[string]string{
		values.HeaderRequestID: tc.RequestID,
	}
}
