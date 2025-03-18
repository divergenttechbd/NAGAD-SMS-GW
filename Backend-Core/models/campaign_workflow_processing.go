package models

// CampaignWorkflowProcessing represents the association between a campaign and a workflow
// @Description Represents the association between a campaign and a workflow
type CampaignWorkflowProcessing struct {
	BaseModel
	// CampaignID is the ID of the campaign (foreign key)
	CampaignID uint `gorm:"not null" json:"campaign_id"`

	// WorkflowID is the ID of the workflow (foreign key)
	WorkflowID uint `gorm:"not null" json:"workflow_id"`

	// Campaign represents the campaign associated with this workflow
	Campaign Campaign `gorm:"foreignKey:CampaignID" json:"campaign"`

	// CampaignWorkflow represents the workflow associated with this campaign
	CampaignWorkflow CampaignWorkflow `gorm:"foreignKey:WorkflowID" json:"campaign_workflow"`
}
