package web

import (
	"github.com/go-playground/validator/v10"
	"github.com/thinkgos/encoding"
)

type Applier interface {
	setEncoding(*encoding.Encoding)
	setValidation(*validator.Validate)
	setTransformError(TransformError)
	setTransformBody(TransformBody)
}

type Option func(Applier)

func WithEncoding(e *encoding.Encoding) Option {
	return func(cy Applier) {
		cy.setEncoding(e)
	}
}
func WithValidation(v *validator.Validate) Option {
	return func(cy Applier) {
		cy.setValidation(v)
	}
}

func WithTransformError(t TransformError) Option {
	return func(cy Applier) {
		cy.setTransformError(t)
	}
}
func WithTransformBody(t TransformBody) Option {
	return func(cy Applier) {
		cy.setTransformBody(t)
	}
}
