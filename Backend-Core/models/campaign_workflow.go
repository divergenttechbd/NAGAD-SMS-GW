package models

// CampaignWorkflow represents a workflow for a campaign
// @Description Represents a workflow for a campaign
type CampaignWorkflow struct {
	BaseModel
	// Name is the name of the workflow
	Name string `gorm:"not null" json:"name"`

	// CampaignWorkflowUsers represents the users associated with this workflow
	CampaignWorkflowUsers []CampaignWorkflowUser `gorm:"foreignKey:WorkflowID" json:"campaign_workflow_users"`

	// CampaignWorkflowProcessings represents the campaigns associated with this workflow
	CampaignWorkflowProcessings []CampaignWorkflowProcessing `gorm:"foreignKey:WorkflowID" json:"campaign_workflow_processings"`
}
