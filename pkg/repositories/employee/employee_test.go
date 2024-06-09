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
			zaplogger.InitLogger(global.TestLogFileName)
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
