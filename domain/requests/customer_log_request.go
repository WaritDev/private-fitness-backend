package requests

type ListCustomerLogsRequest struct {
	Page  int32 `form:"page"  json:"page"`
	Limit int32 `form:"limit" json:"limit"`
}

type UpdateCustomerLogRequest struct {
	Timestamp string `json:"timestamp"` // "YYYY-MM-DD HH:MM:SS"
	LogType   string `json:"logType"`   // CHECK_IN|CHECK_OUT|BOOK_SESSION|CANCEL_SESSION
}