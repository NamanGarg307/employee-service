package errs

// Titles
const (
	BadRequestTitle          = "Bad Request Error"
	InternalServerErrorTitle = "Internal Server Error"
	UnathorizedErrorTitle    = "Unauthorized Error"
)

// Error Message
const (
	InternalServerErrorMessage = "Sorry! Something went wrong"
	PayloadShouldBeEmpty       = "Sorry! Body payload should be empty"
	BadQueryParams             = "Bad query params"
)

// Error message details
const (
	BadRequestErrorMessageDetail       = "Please, give correct input in the body payload"
	SyntaxErrorMessageDetatil          = "Please, Body Payload format is not correct"
	InputErrorMessageDetatil           = "Please, give correct input value for the %q field"
	MissingFieldErrorMessageDetail     = "Required params are missing"
	BodyPayloadLimitErrorMessageDetail = "Request body must not be larger than 1MB"
	UnprocessableEntityMessage         = "Request Not Processed"
	TooManyRequests                    = "Too Many Requests"
)

// DB Errors
const (
	CommitTransactionError = "Transaction Commit Error"
)

// General Errors
const (
	StartServerError    = "Start Server Error"
	APIServerStartError = "Error while Starting API Server"
	InitiateLoggerError = "Initiate Logger Error"

	// Decoder Errors
	StructDecodeError = "Error while decoding struct"
	ConvertToIntError = "Error while converting to Int Error"
)

// Employees
const (
	DecodeEmployeesPOSTError   = "Error while decoding Employee POST request"
	DecodeEmployeePUTError     = "Error while decoding Employee PUT request"
	EmployeeNewRecordError     = "Error while creating record for employee"
	EmployeeNoRecordFoundError = "Invalid Employee ID"
	EmployeeFetchRecordsError  = "Error while fetching employee Records"
	EmployeeUpdateError        = "Error while updating employee from db"
	DeleteEmployeeError        = "Error while deleting employee from db"
	DecodeEmployeesStructError = "Error while decoding employees struct"
)
