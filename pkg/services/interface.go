package service

import (
	"context"

	"github.com/jainabhishek5986/employee-records/pkg/global"
)

/*
EmployeeService : Interface for Employee Service
*/
type EmployeeService interface {
	CreateEmployee(ctx context.Context, request global.DecodeEmployeesPOSTRequest) error
	GetEmployeeByID(ctx context.Context, id int) (global.SuccessGETInfo, error)
	UpdateEmployeeByID(ctx context.Context, request global.DecodeEmployeePUTRequest) error
	DeleteEmployeeByID(ctx context.Context, id int) error
	GetAllEmployee(ctx context.Context, queryParams map[string][]string) (global.SuccessGETInfo, error)
}
