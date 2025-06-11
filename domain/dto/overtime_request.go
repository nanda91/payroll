package dto

type OvertimeRequest struct {
	Date        string  `json:"date" binding:"required"`
	Hours       float64 `json:"hours" binding:"required,min=0.5,max=3"`
	Description string  `json:"description"`
}
