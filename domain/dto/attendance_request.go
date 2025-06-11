package dto

type AttendanceRequest struct {
	Date     string `json:"date" binding:"required"`
	CheckIn  string `json:"check_in" binding:"required"`
	CheckOut string `json:"check_out,omitempty"`
}
