package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"gorm.io/gorm"

	ep "github.com/jainabhishek5986/employee-records/pkg/endpoint/employee"
	svc "github.com/jainabhishek5986/employee-records/pkg/services/employee"
)

func RegisterAPIRoutes(v1RoutesGroup *gin.RouterGroup, db *gorm.DB) {

	var (
		service  = svc.NewService(db)
		endpoint = ep.NewEndPoint(service)
	)

	// Employee Endpoints
	v1RoutesGroup.GET("/employee/:id", NewHTTPHandler(
		endpoint.GetEmployeeByID, DecodeByIDRequest,
		EncodeJSONResponse))

	v1RoutesGroup.GET("/employee", NewHTTPHandler(
		endpoint.GetAllEmployee, DecodeAllRequest,
		EncodeJSONResponse))

	v1RoutesGroup.POST("/employee", NewHTTPHandler(
		endpoint.CreateEmployee, DecodeEmployeesPOSTRequest,
		EncodeJSONResponse))

	v1RoutesGroup.PUT("/employee", NewHTTPHandler(
		endpoint.UpdateEmployeeByID, DecodeEmployeePUTRequest,
		EncodeJSONResponse))

	v1RoutesGroup.DELETE("/employee/:id", NewHTTPHandler(
		endpoint.DeleteEmployeeByID, DecodeByIDRequest,
		EncodeJSONResponse))

	zaplogger.Info(context.Background(), "v1.0 routes injected")
}
