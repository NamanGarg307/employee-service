package employee

import (
	"context"
	"errors"
	"github.com/jainabhishek5986/employee-records/pkg/repositories"
	"strconv"

	"github.com/jainabhishek5986/employee-records/pkg/errs"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/models"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) repositories.EmployeeRepository {
	return &Repository{db: db}
}

// CreateEmployee
func (repo *Repository) CreateEmployee(ctx context.Context, req global.DecodeEmployeesPOSTRequest) error {
	employees := make([]models.Employee, 0)
	var employee models.Employee

	for _, emp := range req.Employees {
		employee := models.Employee{
			Name:     emp.Name,
			Position: emp.Position,
			Salary:   emp.Salary,
		}

		employees = append(employees, employee)
	}

	err := repo.db.Table(employee.GetTableName()).Create(&employees).Error
	if err != nil {
		zaplogger.Error(ctx, errs.EmployeeNewRecordError, zap.Error(err))
		return errs.InternalErr()
	}
	zaplogger.Info(ctx, global.EmployeeCreatedSuccessfully)

	return nil
}

// GetEmployeeByID
func (repo *Repository) GetEmployeeByID(ctx context.Context, id int) (response global.SuccessGETInfo, err error) {
	var employee models.Employee
	tx := repo.db.Begin()
	res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table(employee.GetTableName()).Where("id = ?", id).Find(&employee)

	if res.Error != nil {
		tx.Rollback()
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return response, nil
		}
		zaplogger.Error(ctx, errs.EmployeeNoRecordFoundError, zap.Error(res.Error))
		return response, errs.InternalErr()
	}
	err = tx.Commit().Error
	if err != nil {
		zaplogger.Error(ctx, errs.CommitTransactionError, zap.Error(err))
		return response, err
	}
	if employee.ID == 0 {
		return response, errs.RequestNotProcessed(errs.EmployeeNoRecordFoundError)
	}

	response = global.SuccessGETInfo{
		Data: employee,
	}

	return response, nil

}

// UpdateEmployeeByID
func (repo *Repository) UpdateEmployeeByID(ctx context.Context, request global.DecodeEmployeePUTRequest) error {
	var employee models.Employee
	if request.Name != nil {
		employee.Name = *request.Name
	}
	if request.Position != nil {
		employee.Position = *request.Position
	}
	if request.Salary != nil {
		employee.Salary = *request.Salary
	}

	tx := repo.db.Begin()
	res := tx.Table(employee.GetTableName()).Where("id = ?", request.ID).Updates(employee)
	if res.Error != nil {
		tx.Rollback()
		zaplogger.Error(ctx, errs.EmployeeUpdateError, zap.Error(res.Error),
			zap.Int("employee_id", request.ID),
		)
		return errs.InternalErr()
	}

	// Check if any rows were affected
	if res.RowsAffected == 0 {
		tx.Rollback()
		zaplogger.Error(ctx, errs.EmployeeNoRecordFoundError, zap.Int("employee_id", request.ID))
		return errs.RequestNotProcessed(errs.EmployeeNoRecordFoundError)
	}

	err := tx.Commit().Error
	if err != nil {
		zaplogger.Error(ctx, errs.CommitTransactionError, zap.Error(err))
		return err
	}
	zaplogger.Info(ctx, global.EmployeeUpdatedSuccessfully,
		zap.Int("employee_id", request.ID),
	)

	return nil
}

// DeleteEmployeeByID
func (repo *Repository) DeleteEmployeeByID(ctx context.Context, id int) error {
	var employee models.Employee
	employee.ID = id

	tx := repo.db.Begin()

	res := tx.Table(employee.GetTableName()).Where(&employee).Delete(&employee)
	if res.Error != nil {
		tx.Rollback()
		zaplogger.Error(ctx, errs.DeleteEmployeeError, zap.Error(res.Error),
			zap.Int("employee_id", id),
		)
		return errs.InternalErr()
	}

	// Check if any rows were affected
	if res.RowsAffected == 0 {
		tx.Rollback()
		zaplogger.Error(ctx, errs.EmployeeNoRecordFoundError, zap.Int("employee_id", id))
		return errs.RequestNotProcessed(errs.EmployeeNoRecordFoundError)
	}
	err := tx.Commit().Error
	if err != nil {
		zaplogger.Error(ctx, errs.CommitTransactionError, zap.Error(err))
		return err
	}
	zaplogger.Info(ctx, global.EmployeeDeletedSuccessfully,
		zap.Int("employee_id", id),
	)

	return nil
}

// GetAllEmployee
func (repo *Repository) GetAllEmployee(ctx context.Context, queryParams map[string][]string) (response global.SuccessGETInfo, err error) {
	var employees []models.Employee
	var employee models.Employee
	var totalCount int64

	// Calculate offset based on the current page and page size
	currentPage := 1
	if pageQuery, isExist := queryParams["page"]; isExist {
		currentPage, err = strconv.Atoi(pageQuery[0])
		if err != nil {
			zaplogger.Error(ctx, errs.ConvertToIntError)
			return response, errs.InternalErr()
		}
	}
	pageSize := 10
	if pageSizeQuery, isExist := queryParams["per_page"]; isExist {
		pageSize, err = strconv.Atoi(pageSizeQuery[0])
		if err != nil {
			zaplogger.Error(ctx, errs.ConvertToIntError)
			return response, errs.InternalErr()
		}
	}
	offset := (currentPage - 1) * pageSize
	tx := repo.db.Begin()

	// Get the total count of employees
	if err := tx.Table(employee.GetTableName()).Count(&totalCount).Error; err != nil {
		tx.Rollback()
		zaplogger.Error(ctx, errs.EmployeeFetchRecordsError, zap.Error(err))
		return response, errs.InternalErr()
	}
	res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Table(employee.GetTableName()).
		Limit(pageSize).
		Offset(offset).
		Find(&employees)

	if res.Error != nil {
		tx.Rollback()
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return response, nil
		}
		zaplogger.Error(ctx, errs.EmployeeFetchRecordsError, zap.Error(res.Error))
		return response, errs.InternalErr()
	}
	err = tx.Commit().Error
	if err != nil {
		zaplogger.Error(ctx, errs.CommitTransactionError, zap.Error(err))
		return response, err
	}

	lastPage := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	paginationResponse := map[string]int{
		"total":        int(totalCount),
		"current_page": currentPage,
		"last_page":    lastPage,
	}
	response = global.SuccessGETInfo{
		Data:       employees,
		Pagination: paginationResponse,
	}

	return response, nil

}
