package employee

import (
	"context"
	"github.com/jainabhishek5986/employee-records/pkg/repositories"

	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/repositories/employee"
	services "github.com/jainabhishek5986/employee-records/pkg/services"
	"gorm.io/gorm"
)

// Employee Service Structure
type service struct {
	db   *gorm.DB
	repo repositories.EmployeeRepository
}

func NewService(db *gorm.DB) services.EmployeeService {

	repo := employee.NewEmployeeRepo(db)
	return &service{db: db, repo: repo}
}

func (envSvc *service) CreateEmployee(ctx context.Context, req global.DecodeEmployeesPOSTRequest) error {
	return envSvc.repo.CreateEmployee(ctx, req)
}

func (envSvc *service) GetEmployeeByID(ctx context.Context, id int) (global.SuccessGETInfo, error) {
	return envSvc.repo.GetEmployeeByID(ctx, id)
}

func (envSvc *service) UpdateEmployeeByID(ctx context.Context, request global.DecodeEmployeePUTRequest) error {
	return envSvc.repo.UpdateEmployeeByID(ctx, request)
}

func (envSvc *service) DeleteEmployeeByID(ctx context.Context, id int) error {
	return envSvc.repo.DeleteEmployeeByID(ctx, id)
}

func (envSvc *service) GetAllEmployee(ctx context.Context, queryParams map[string][]string) (global.SuccessGETInfo, error) {
	return envSvc.repo.GetAllEmployee(ctx, queryParams)
}
