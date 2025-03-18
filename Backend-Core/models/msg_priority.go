package models

// PriorityConfiguration represents the priority levels for different types of SMS
// @Description Represents the priority levels for SMS types (e.g., OTP, Transaction, Promotional)
type MsgPriority struct {
	BaseModel
	// Message_Type is the type of SMS (e.g., OTP, Transaction, Promotional)
	Message_Type string `gorm:"not null" json:"message_type"`

	// Priority_Level is the priority level for the SMS type (0 = Highest, 3 = Lowest)
	Priority_Level int `gorm:"not null" json:"priority_level"`

	// Description provides additional details about the priority level
	Description string `gorm:"not null" json:"description"`
}
