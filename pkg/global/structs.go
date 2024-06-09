package global

type DecodeEmployeesPOSTRequest struct {
	Employees []DecodeEmployee `json:"employees" validate:"required"`
}

type DecodeEmployeePUTRequest struct {
	ID       int      `json:"id"`
	Name     *string  `json:"name"`
	Position *string  `json:"position"`
	Salary   *float64 `json:"salary"`
}

type DecodeEmployee struct {
	Name     string  `json:"name" validate:"required"`
	Position string  `json:"position" validate:"required"`
	Salary   float64 `json:"salary" validate:"required"`
}
