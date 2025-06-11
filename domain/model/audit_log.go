package model

type AuditLog struct {
	BaseModel
	UserID    *uint  `json:"user_id,omitempty"`
	Action    string `json:"action"`
	TableName string `json:"table_name"`
	RecordID  *uint  `json:"record_id,omitempty"`
	OldData   string `json:"old_data,omitempty"`
	NewData   string `json:"new_data,omitempty"`

	// Relationships
	User *User `json:"user,omitempty"`
}
