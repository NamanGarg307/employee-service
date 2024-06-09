package global

/*
SuccessInfo : Success message
*/
type SuccessInfo struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type SuccessGETInfo struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
}
