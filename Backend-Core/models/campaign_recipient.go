package models

// CampaignRecipient represents a recipient of a campaign and their delivery status
// @Description Represents a recipient of a campaign and their delivery status
type CampaignRecipient struct {
	BaseModel
	// CampaignID is the ID of the campaign (foreign key)
	CampaignID uint `gorm:"not null" json:"campaign_id"`

	// RecipientID is the ID of the recipient (foreign key)
	Recipient uint `gorm:"not null" json:"recipient"`

	// Status indicates the delivery status of the campaign to the recipient
	Status string `gorm:"not null" json:"status"`

	// Campaign represents the campaign associated with this recipient
	Campaign Campaign `gorm:"foreignKey:CampaignID" json:"campaign"`
}
