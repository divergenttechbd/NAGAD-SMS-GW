package models

// DND represents a Do Not Disturb entry
// @Description Represents a phone number that should not receive promotional SMS
type DND struct {
	BaseModel
	// Phone_Number is the number that should not receive promotional SMS
	Phone_Number string `gorm:"not null" json:"phone_number"`

	// Reason provides the reason for adding the number to the DND list
	Reason string `gorm:"not null" json:"reason"`

	// Status indicates whether the DND entry is active or inactive
	Status string `gorm:"not null" json:"status"`
}
