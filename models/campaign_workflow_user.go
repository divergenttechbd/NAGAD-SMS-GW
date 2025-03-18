package models

// CampaignWorkflowUser represents a user associated with a workflow
// @Description Represents a user associated with a workflow and their status
type CampaignWorkflowUser struct {
	BaseModel
	// WorkflowID is the ID of the workflow (foreign key)
	WorkflowID uint `gorm:"not null" json:"workflow_id"`

	// UserID is the ID of the user (foreign key)
	UserID uint `gorm:"not null" json:"user_id"`

	// Status indicates whether the user is active or inactive in the workflow
	Status string `gorm:"not null" json:"status"`

	// CampaignWorkflow represents the workflow associated with this user
	CampaignWorkflow CampaignWorkflow `gorm:"foreignKey:WorkflowID" json:"campaign_workflow"`

	// User represents the user associated with this workflow
	User User `gorm:"foreignKey:UserID" json:"user"`
}
