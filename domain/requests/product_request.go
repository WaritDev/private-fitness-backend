package requests

type ProductFilterRequest struct {
	Type     string `json:"type"`     // DURATION or SESSION
	Category string `json:"category"` // ECONOMIC, BUSINESS, FIRST_CLASS
}
