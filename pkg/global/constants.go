package global

const (
	MaxAPIServerStartAttempts = 10
	MaxConnections            = 100
	MaxLifeTime               = 3
)

// Variable Constants
const (
	Success         = "Success"
	LogFileName     = "employee-records-service.log"
	TestLogFileName = "./../employee-records-service-unit-tests.log"
	SQL             = "mysql"
)

// Global Magic numbers
const (
	FiveHundred   = 500
	PointZeroFive = 0.05
	FiveThousand  = 5000
	Hundred       = 100
	One           = 1
	Ten           = 10
	SixtyFour     = 64
)

// API routes
const (
	CreateEmployeeEndpoint  = "POST: /employee"
	UpdateEmployeeEndpoint  = "PUT: /employee"
	GetEmployeeByIDEndpoint = "GET: /employee/:id"
	DeleteEmployeeEndpoint  = "DELETE: /employee/:id"
)
