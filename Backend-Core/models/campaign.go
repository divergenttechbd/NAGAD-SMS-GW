package models

import "time"

// Campaign represents an SMS campaign
// @Description Represents an SMS campaign associated with a user
type Campaign struct {
	BaseModel
	// UserID is the ID of the user who created the campaign
	UserID uint `gorm:"not null" json:"user_id"`

	// Campaign_Name is the name of the campaign
	Campaign_Name string `gorm:"not null" json:"campaign_name"`

	// Message_Body is the body of the SMS message for the campaign
	Message_Body string `gorm:"not null" json:"message_body"`

	// Start_Date is the start date of the campaign
	Start_Date time.Time `gorm:"not null" json:"start_date"`

	// End_Date is the end date of the campaign
	End_Date time.Time `gorm:"not null" json:"end_date"`

	// Status indicates the status of the campaign (e.g., Active, Inactive)
	Status string `gorm:"not null" json:"status"`

	// CampaignWorkflowProcessings represents the workflows associated with this campaign
	CampaignWorkflowProcessings []CampaignWorkflowProcessing `gorm:"foreignKey:CampaignID" json:"campaign_workflow_processings"`

	// CampaignRecipients represents the recipients associated with this campaign
	CampaignRecipients []CampaignRecipient `gorm:"foreignKey:CampaignID" json:"campaign_recipients"`
}
