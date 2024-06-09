package models

import "time"

// Attachments - It stores all the attachements.
type Employee struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Position  string     `json:"position"`
	Salary    float64    `json:"salary"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (m *Employee) GetTableName() string {
	return "employees"
}
