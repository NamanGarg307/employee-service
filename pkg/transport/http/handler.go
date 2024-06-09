package http

import (
	"context"
	"encoding/json"
	"fmt"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	gohttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	"github.com/jainabhishek5986/employee-records/pkg/errs"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"go.uber.org/zap"
)

type DecodeRequestFunc func(context.Context, *gin.Context) (request interface{}, err error)

type EncodeResponseFunc func(context.Context, *gin.Context,
	interface{}) error

func NewHTTPHandler(ep endpoint.Endpoint, dec DecodeRequestFunc,
	enc EncodeResponseFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorEncoder := gohttp.DefaultErrorEncoder
		request, err := dec(c, c)
		if err != nil {
			errorEncoder(c, err, c.Writer)
			return
		}
		response, err := ep(c, request)
		if err != nil {
			errorEncoder(c, err, c.Writer)
			return
		}

		if err := enc(c, c, response); err != nil {
			zaplogger.Error(c, "err", zap.Error(err))
			errorEncoder(c, err, c.Writer)
			return
		}
	}
}

func EncodeJSONResponse(_ context.Context, c *gin.Context,
	response interface{}) error {

	// adding a header content-type
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	if headerer, ok := response.(gohttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				c.Writer.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(gohttp.StatusCoder); ok {
		code = sc.StatusCode()
	}

	c.Writer.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	return json.NewEncoder(c.Writer).Encode(response)
}

func translateError(ctx context.Context, err error, validate *validator.Validate,
	decodeStruct reflect.Value) (errMessage map[string]string, internalError error) {

	if err == nil {
		return nil, nil
	}

	// Converting error into string format
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		zaplogger.Error(ctx, "Converting error into string translator error", zap.Any("found", found))
		return nil, errs.InternalErr()
	}

	translationError := en_translations.RegisterDefaultTranslations(validate,
		trans)
	if translationError != nil {
		zaplogger.Error(ctx, "Converting error into string translator error", zap.Any("error", translationError.Error()))
		return nil, errs.InternalErr()
	}

	errMessage = make(map[string]string)

	// Checking validation error
	validatorErrs, _ := err.(validator.ValidationErrors)

	for _, e := range validatorErrs {
		// Extract the JSON tag from the struct field, if available
		field, ok := decodeStruct.Type().FieldByName(e.Field())
		jsonFieldName := e.Field()
		if ok && field.Tag.Get("json") != "" {
			jsonFieldName = field.Tag.Get("json")
		}

		// Customize the error message
		var errMsg string
		switch e.Tag() {
		case "required":
			errMsg = fmt.Sprintf("%s is a required field", jsonFieldName)
		case "trimspace":
			errMsg = fmt.Sprintf("%s cannot be just spaces", jsonFieldName)
		// Add more cases here for other validation tags if needed
		default:
			errMsg = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", jsonFieldName, e.Tag())
		}

		errMessage[jsonFieldName] = errMsg
	}

	return errMessage, nil
}

func DecodeEmployeesPOSTRequest(c context.Context, g *gin.Context) (request interface{}, err error) {

	// Checking body payload is empty or not
	ErrMsg := make([]interface{}, 0)
	queryParams := g.Request.URL.Query()

	if len(queryParams) > 0 {
		// Todo : change error message format
		ErrMsg = append(ErrMsg, errs.ErrMessage{
			Key:    "BadPayload",
			Detail: errs.BadQueryParams})
		return nil, errs.ErrResponse(errs.BadRequestTitle,
			http.StatusBadRequest, ErrMsg)
	}

	var decodeEmployeesPOSTRequest global.DecodeEmployeesPOSTRequest
	err = g.ShouldBindJSON(&decodeEmployeesPOSTRequest)
	if err != nil {
		zaplogger.Error(c, errs.DecodeEmployeesPOSTError, zap.Error(err))
		err = errs.ErrorReqHandler(err)
		return nil, err
	}
	err = Validate.Struct(decodeEmployeesPOSTRequest)
	if err != nil {
		zaplogger.Error(c, errs.DecodeEmployeesPOSTError, zap.Error(err))
		val := reflect.ValueOf(global.DecodeEmployeesPOSTRequest{})
		payloadErrorMessages, internalError := translateError(c, err,
			Validate, val)
		if internalError != nil {

			return nil, errs.InternalErr()
		}
		return nil, errs.RequestNotProcessed(payloadErrorMessages)
	}
	for _, decodeEnv := range decodeEmployeesPOSTRequest.Employees {
		err = Validate.Struct(decodeEnv)
		if err != nil {
			zaplogger.Error(c, errs.DecodeEmployeesPOSTError, zap.Error(err))
			val := reflect.ValueOf(global.DecodeEmployee{})
			payloadErrorMessages, internalError := translateError(c, err,
				Validate, val)
			if internalError != nil {

				return nil, errs.InternalErr()
			}
			return nil, errs.RequestNotProcessed(payloadErrorMessages)
		}
	}

	return decodeEmployeesPOSTRequest, nil
}

func DecodeEmployeePUTRequest(c context.Context, g *gin.Context) (request interface{}, err error) {

	// Checking body payload is empty or not
	ErrMsg := make([]interface{}, 0)
	queryParams := g.Request.URL.Query()

	if len(queryParams) > 0 {
		// Todo : change error message format
		ErrMsg = append(ErrMsg, errs.ErrMessage{
			Key:    "BadPayload",
			Detail: errs.BadQueryParams})
		return nil, errs.ErrResponse(errs.BadRequestTitle,
			http.StatusBadRequest, ErrMsg)
	}

	var decodeEmployeePUTRequest global.DecodeEmployeePUTRequest
	err = g.ShouldBindJSON(&decodeEmployeePUTRequest)
	if err != nil {
		zaplogger.Error(c, errs.DecodeEmployeesPOSTError, zap.Error(err))
		err = errs.ErrorReqHandler(err)
		return nil, err
	}
	err = Validate.Struct(decodeEmployeePUTRequest)
	if err != nil {
		zaplogger.Error(c, errs.DecodeEmployeePUTError, zap.Error(err))
		val := reflect.ValueOf(global.DecodeEmployeePUTRequest{})
		payloadErrorMessages, internalError := translateError(c, err,
			Validate, val)
		if internalError != nil {

			return nil, errs.InternalErr()
		}
		return nil, errs.RequestNotProcessed(payloadErrorMessages)
	}

	return decodeEmployeePUTRequest, nil
}

func DecodeByIDRequest(ctx context.Context, g *gin.Context) (request interface{}, err error) {

	// Checking body payload is empty or not
	ErrMsg := make([]interface{}, 0)
	queryParams := g.Request.URL.Query()

	// Query params
	if len(queryParams) > 0 {
		ErrMsg = append(ErrMsg, errs.ErrMessage{
			Key:    "BadPayload",
			Detail: errs.BadQueryParams})
		return nil, errs.ErrResponse(errs.BadRequestTitle,
			http.StatusBadRequest, ErrMsg)
	}

	// Empty body payload
	if g.Request.Body != http.NoBody {
		ErrMsg = append(ErrMsg, errs.ErrMessage{
			Key:    "BadPayload",
			Detail: errs.PayloadShouldBeEmpty})
		return nil, errs.ErrResponse(errs.BadRequestTitle,
			http.StatusBadRequest, ErrMsg)
	}

	id := g.Param("id")
	integerID, err := strconv.Atoi(id)
	if err != nil {
		zaplogger.Error(ctx, errs.ConvertToIntError)
		return nil, errs.BadRequest(errs.ConvertToIntError)
	}

	return integerID, err
}

func DecodeAllRequest(ctx context.Context, g *gin.Context) (request interface{}, err error) {

	// Checking body payload is empty or not
	ErrMsg := make([]interface{}, 0)
	queryParams := g.Request.URL.Query()

	// Empty body payload
	if g.Request.Body != http.NoBody {
		ErrMsg = append(ErrMsg, errs.ErrMessage{
			Key:    "BadPayload",
			Detail: errs.PayloadShouldBeEmpty})
		return nil, errs.ErrResponse(errs.BadRequestTitle,
			http.StatusBadRequest, ErrMsg)
	}
	paramsMap := make(map[string][]string, 0)

	// Query params
	if len(queryParams) > 0 {
		for key, values := range queryParams {
			paramsMap[key] = values
		}

	}

	return paramsMap, err
}
