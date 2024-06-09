package employee

import (
	"context"
	"testing"

	"github.com/jainabhishek5986/employee-records/pkg/errs"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/models"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db, %v", err)
	}

	err = db.AutoMigrate(&models.Employee{})
	if err != nil {
		t.Fatalf("failed to migrate schema, %v", err)
	}

	zaplogger.InitLogger(global.TestLogFileName)
	return db
}

func TestGetAllEmployee(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEmployeeRepo(db)

	// Seed data for testing
	employees := []models.Employee{
		{Name: "John Doe", Position: "Manager", Salary: 50000},
		{Name: "Jane Smith", Position: "Developer", Salary: 60000},
	}
	db.Create(&employees)

	// Define the test cases
	testCases := []struct {
		name          string
		queryParams   map[string][]string
		expectedError error
		expectedCount int
	}{
		{
			name:          "Successful fetch with pagination",
			queryParams:   map[string][]string{"page": {"1"}, "per_page": {"10"}},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name:          "Fetch with no records",
			queryParams:   map[string][]string{"page": {"2"}, "per_page": {"10"}},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name:          "Error converting page parameter",
			queryParams:   map[string][]string{"page": {"abc"}, "per_page": {"10"}},
			expectedError: errs.InternalErr(),
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := repo.GetAllEmployee(ctx, tc.queryParams)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(response.Data.([]models.Employee)))
			}
		})
	}
}

func TestCreateEmployee(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEmployeeRepo(db)

	decodeEmployeePOSTRequest := global.DecodeEmployeesPOSTRequest{
		Employees: []global.DecodeEmployee{{
			Name:     "Abhishek",
			Position: "Engineer",
			Salary:   70000,
		}},
	}
	employee := &models.Employee{
		Name:     "Abhishek",
		Position: "Engineer",
		Salary:   70000,
	}

	ctx := context.Background()

	err := repo.CreateEmployee(ctx, decodeEmployeePOSTRequest)
	assert.NoError(t, err)

	var result models.Employee
	err = db.First(&result).Error
	assert.NoError(t, err)
	assert.Equal(t, employee.Name, result.Name)
	assert.Equal(t, employee.Position, result.Position)
	assert.Equal(t, employee.Salary, result.Salary)
}

func TestUpdateEmployee(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEmployeeRepo(db)

	name := "Abhishek Updated"
	position := "Senior Engineer"
	salary := 70000.2

	employee := &models.Employee{
		Name:     "Abhishek",
		Position: "Engineer",
		Salary:   70000,
	}
	db.Create(employee)

	updatedEmployee := global.DecodeEmployeePUTRequest{
		ID:       employee.ID,
		Name:     &name,
		Position: &position,
		Salary:   &salary,
	}

	ctx := context.Background()

	err := repo.UpdateEmployeeByID(ctx, updatedEmployee)
	assert.NoError(t, err)

	var result models.Employee
	err = db.First(&result, employee.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, *updatedEmployee.Name, result.Name)
	assert.Equal(t, *updatedEmployee.Position, result.Position)
	assert.Equal(t, *updatedEmployee.Salary, result.Salary)
}

func TestDeleteEmployee(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEmployeeRepo(db)

	employee := &models.Employee{
		Name:     "Alice",
		Position: "Engineer",
		Salary:   70000,
	}
	db.Create(employee)

	ctx := context.Background()

	err := repo.DeleteEmployeeByID(ctx, employee.ID)
	assert.NoError(t, err)

	var result models.Employee
	err = db.First(&result, employee.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
