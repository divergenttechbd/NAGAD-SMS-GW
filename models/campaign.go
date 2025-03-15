package models

// Campaign represents a campaign in the system
// @Description A marketing or messaging campaign associated with users
type Campaign struct {
	BaseModel
	// Name is the name of the campaign
	Name string `gorm:"not null" json:"name"`

	// Status represents the campaign status
	Status string `gorm:"not null" json:"status"`

	// UserID is the ID of the user associated with the campaign
	UserID uint `gorm:"not null" json:"user_id"`
}
