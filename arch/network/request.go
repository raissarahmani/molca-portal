package network

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ReqBodyMultipart[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindWith(dto, binding.FormMultipart); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqBody[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindJSON(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqQuery[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindQuery(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqParams[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindUri(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqHeaders[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindHeader(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func processErrors[T any](dto Dto[T], err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		msgs, e := dto.ValidateErrors(validationErrors)
		if e != nil {
			return e
		}
		var msg strings.Builder
		br := ", "
		for _, m := range msgs {
			msg.WriteString(m + br)
		}
		// Remove the trailing separator
		errorMsg := msg.String()
		if len(errorMsg) > 0 {
			errorMsg = errorMsg[:len(errorMsg)-len(br)]
		}
		return errors.New(errorMsg)
	}
	return err
}
