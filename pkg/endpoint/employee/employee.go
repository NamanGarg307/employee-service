package employee

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jainabhishek5986/employee-records/pkg/errs"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	service "github.com/jainabhishek5986/employee-records/pkg/services"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
)

// EndPoints : All the Employee endpoints structure
type EndPoints struct {
	CreateEmployee     endpoint.Endpoint
	GetEmployeeByID    endpoint.Endpoint
	UpdateEmployeeByID endpoint.Endpoint
	DeleteEmployeeByID endpoint.Endpoint
	GetAllEmployee     endpoint.Endpoint
}

func NewEndPoint(svc service.EmployeeService) EndPoints {

	return EndPoints{
		CreateEmployee:     makeCreateEmployee(svc),
		GetEmployeeByID:    makeGetEmployeeByID(svc),
		UpdateEmployeeByID: makeUpdateEmployeeByID(svc),
		DeleteEmployeeByID: makeDeleteEmployeeByID(svc),
		GetAllEmployee:     makeGetAllEmployee(svc),
	}
}

func makeCreateEmployee(svc service.EmployeeService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{},
		err error) {
		req, ok := request.(global.DecodeEmployeesPOSTRequest)
		if !ok {
			zaplogger.Error(ctx, errs.DecodeEmployeesStructError)
			return nil, errs.InternalErr()
		}
		err = svc.CreateEmployee(ctx, req)
		if err != nil {

			return nil, err
		}

		return global.SuccessInfo{
			Message: global.EmployeeCreatedSuccessfully,
			Type:    global.Success,
		}, err
	}
}

func makeGetEmployeeByID(svc service.EmployeeService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{},
		err error) {
		req, ok := request.(int)
		if !ok {
			zaplogger.Error(ctx, errs.ConvertToIntError)
			return nil, errs.InternalErr()
		}
		res, err := svc.GetEmployeeByID(ctx, req)
		// Error handling
		if err != nil {
			return nil, err
		}

		return res, err
	}
}

func makeUpdateEmployeeByID(svc service.EmployeeService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{},
		err error) {
		req, ok := request.(global.DecodeEmployeePUTRequest)
		if !ok {
			zaplogger.Error(ctx, errs.DecodeEmployeePUTError)
			// Send Slack Alert Here
			return nil, errs.InternalErr()
		}

		err = svc.UpdateEmployeeByID(ctx, req)
		// Error handling
		if err != nil {
			return nil, err
		}

		return global.SuccessInfo{
			Message: global.EmployeeUpdatedSuccessfully,
			Type:    global.Success,
		}, err
	}
}

func makeDeleteEmployeeByID(svc service.EmployeeService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{},
		err error) {
		req, ok := request.(int)
		if !ok {
			zaplogger.Error(ctx, errs.ConvertToIntError)
			// Send Slack Alert Here
			return nil, errs.InternalErr()
		}

		err = svc.DeleteEmployeeByID(ctx, req)
		// Error handling
		if err != nil {
			return nil, err
		}

		return global.SuccessInfo{
			Message: global.EmployeeDeletedSuccessfully,
			Type:    global.Success,
		}, err
	}
}

func makeGetAllEmployee(svc service.EmployeeService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{},
		err error) {
		req, ok := request.(map[string][]string)
		if !ok {
			zaplogger.Error(ctx, errs.StructDecodeError)
			// Send Slack Alert Here
			return nil, errs.InternalErr()
		}
		res, err := svc.GetAllEmployee(ctx, req)
		// Error handling
		if err != nil {
			return nil, err
		}

		return res, err
	}
}
