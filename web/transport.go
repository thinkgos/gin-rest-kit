package web

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxTransportKey struct{}

type Transporter interface {
	Gin() *gin.Context
}

type Transport struct {
	c *gin.Context
}

func (r *Transport) Gin() *gin.Context { return r.c }

func WithValueTransporter(ctx context.Context, c Transporter) context.Context {
	return context.WithValue(ctx, ctxTransportKey{}, c)
}

// FromCarrier returns the Carrier value stored in ctx, if not exist cause panic.
func FromTransporter(ctx context.Context) Transporter {
	c, ok := ctx.Value(ctxTransportKey{}).(Transporter)
	if !ok {
		panic("transport: must be set Transporter into context but it is not!!!")
	}
	return c
}

func TransportInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(WithValueTransporter(c.Request.Context(), &Transport{c}))
		c.Next()
	}
}
