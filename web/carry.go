package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thinkgos/encoding"
	"github.com/thinkgos/encoding/form"
)

var _ Carrier = (*Carry)(nil)
var _ Applier = (*Carry)(nil)

type Carry struct {
	encoding       *encoding.Encoding
	validation     *validator.Validate
	transformError TransformError
	transformBody  TransformBody
}

func NewCarry(opts ...Option) *Carry {
	e := encoding.New()
	err := e.Register(encoding.Mime_Query, &form.QueryCodec{
		Codec: form.New("json").RegisterBuiltinTypeDecoderCommaStringToSlice(),
	})
	if err != nil {
		panic(fmt.Errorf("carry: 初始化Query编解码器失败. %w", err))
	}
	err = e.Register(encoding.Mime_Uri, &form.QueryCodec{
		Codec: form.New("json").RegisterBuiltinTypeDecoderCommaStringToSlice(),
	})
	if err != nil {
		panic(fmt.Errorf("carry: 初始化URI编解码器失败. %w", err))
	}

	cy := &Carry{
		encoding: e,
		validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
	}
	for _, opt := range opts {
		opt(cy)
	}

	return cy
}

func (cy *Carry) setEncoding(e *encoding.Encoding) {
	cy.encoding = e
}
func (cy *Carry) setValidation(v *validator.Validate) {
	cy.validation = v
}

func (cy *Carry) setTransformError(e TransformError) {
	cy.transformError = e
}

func (cy *Carry) setTransformBody(e TransformBody) {
	cy.transformBody = e
}

func (cy *Carry) Bind(c *gin.Context, v any) error {
	return cy.encoding.Bind(c.Request, v)
}
func (cy *Carry) BindQuery(c *gin.Context, v any) error {
	return cy.encoding.BindQuery(c.Request, v)
}
func (cy *Carry) BindUri(c *gin.Context, v any) error {
	return cy.encoding.BindUri(UrlValues(c.Params), v)
}
func (cy *Carry) ShouldBind(c *gin.Context, v any) error {
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldBindQuery(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldBindUri(c *gin.Context, v any) error {
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldBindBodyUri(c *gin.Context, v any) error {
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldBindQueryUri(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldBindQueryBody(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldBindQueryBodyUri(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if err := cy.Bind(c, v); err != nil {
		return err
	}
	if err := cy.BindUri(c, v); err != nil {
		return err
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) ShouldAutoBind(c *gin.Context, v any) error {
	if err := cy.BindQuery(c, v); err != nil {
		return err
	}
	if method := c.Request.Method; c.Request.ContentLength > 0 && (method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch) {
		if err := cy.Bind(c, v); err != nil {
			return err
		}
	}
	if len(c.Params) > 0 {
		if err := cy.BindUri(c, v); err != nil {
			return err
		}
	}
	return cy.Validate(c.Request.Context(), v)
}
func (cy *Carry) Error(c *gin.Context, err error) {
	var obj any
	var statusCode = http.StatusInternalServerError

	if cy.transformError != nil {
		statusCode, obj = cy.transformError.TransformError(c.Request.Context(), err)
	} else {
		obj = err.Error()
	}
	c.Writer.WriteHeader(statusCode)
	if err := cy.encoding.Render(c.Writer, c.Request, obj); err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (cy *Carry) Render(c *gin.Context, v any) {
	if cy.transformBody != nil {
		v = cy.transformBody.TransformBody(c.Request.Context(), v)
	}
	c.Writer.WriteHeader(http.StatusOK)
	err := cy.encoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (cy *Carry) Validator() *validator.Validate {
	return cy.validation
}
func (cy *Carry) Validate(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *Carry) StructCtx(ctx context.Context, v any) error {
	return cy.validation.StructCtx(ctx, v)
}
func (cy *Carry) Struct(v any) error {
	return cy.validation.Struct(v)
}
func (cy *Carry) VarCtx(ctx context.Context, v any, tag string) error {
	return cy.validation.VarCtx(ctx, v, tag)
}
func (cy *Carry) Var(v any, tag string) error {
	return cy.validation.Var(v, tag)
}
