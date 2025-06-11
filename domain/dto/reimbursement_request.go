package dto

type ReimbursementRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0"`
	Description string  `json:"description" binding:"required"`
}
