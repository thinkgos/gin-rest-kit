package web

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxCarrierKey struct{}

// Carrier is an HTTP Carrier.
type Carrier interface {
	// Bind checks the Method and Content-Type to select codec.Marshaler automatically,
	// Depending on the "Content-Type" header different bind are used.
	Bind(*gin.Context, any) error
	// BindQuery binds the passed struct pointer using the query codec.Marshaler.
	BindQuery(*gin.Context, any) error
	// BindUri binds the passed struct pointer using the uri codec.Marshaler.
	BindUri(*gin.Context, any) error
	// ShouldBind checks the Method and Content-Type to select codec.Marshaler automatically then validate the request,
	// Depending on the "Content-Type" header different bind are used.
	ShouldBind(*gin.Context, any) error
	// ShouldBindQuery binds the passed struct pointer using the query codec.Marshaler then validate the request.
	ShouldBindQuery(*gin.Context, any) error
	// ShouldBindUri binds the passed struct pointer using the uri codec.Marshaler then validate the request.
	ShouldBindUri(*gin.Context, any) error
	// ShouldBindQueryUri binds the passed struct pointer using the query and uri codec.Marshaler then validate the request.
	ShouldBindQueryUri(*gin.Context, any) error
	// ShouldBindBodyUri binds the passed struct pointer using the body and uri codec.Marshaler then validate the request.
	ShouldBindBodyUri(*gin.Context, any) error
	// ShouldBindQueryBody binds the passed struct pointer using the query and body codec.Marshaler then validate the request.
	ShouldBindQueryBody(*gin.Context, any) error
	// ShouldBindQueryBodyUri binds the passed struct pointer using query, body, uri codec.Marshaler if necessary.
	ShouldBindQueryBodyUri(*gin.Context, any) error
	// ShouldAutoBind auto binds the passed struct pointer using query, body, uri codec.Marshaler if necessary.
	ShouldAutoBind(*gin.Context, any) error
	// Error encode error response.
	Error(*gin.Context, error)
	// Render encode response.
	Render(*gin.Context, any)
	// Validate the request.
	Validate(context.Context, any) error
}

// WithValueCarrier returns the value associated with ctxCarrierKey is
// Carrier.
func WithValueCarrier(ctx context.Context, c Carrier) context.Context {
	return context.WithValue(ctx, ctxCarrierKey{}, c)
}

// FromCarrier returns the Carrier value stored in ctx, if not exist cause panic.
func FromCarrier(ctx context.Context) Carrier {
	c, ok := ctx.Value(ctxCarrierKey{}).(Carrier)
	if !ok {
		panic("carrier: must be set Carrier into context but it is not!!!")
	}
	return c
}

// CarrierInterceptor carrier middleware.
func CarrierInterceptor(carrier Carrier) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(WithValueCarrier(c.Request.Context(), carrier))
		c.Next()
	}
}
