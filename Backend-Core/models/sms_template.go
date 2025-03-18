package models

// SMSTemplate represents a predefined SMS template
// @Description Represents a predefined SMS template for different types of messages
type SMSTemplate struct {
	BaseModel
	// Template_Name is the name of the SMS template
	Template_Name string `gorm:"not null" json:"template_name"`

	// Template_Body is the body of the SMS template
	Template_Body string `gorm:"not null" json:"template_body"`

	// Message_Type is the type of message the template is used for (e.g., OTP, Transaction)
	Message_Type string `gorm:"not null" json:"message_type"`

	// Status indicates whether the template is active or inactive
	Status string `gorm:"not null" json:"status"`
}
