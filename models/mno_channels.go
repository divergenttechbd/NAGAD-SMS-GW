package models

// Channel represents a delivery channel for an MNO
// @Description Represents a delivery channel (e.g., HTTP, SMPP) for an MNO
type MnoChannels struct {
	BaseModel
	// ChannelID is the unique identifier for the channel
	ChannelID uint `gorm:"primaryKey;autoIncrement" json:"channel_id"`

	// MNOID is the ID of the MNO associated with this channel
	MNOID uint `gorm:"not null" json:"mno_id"`

	// ChannelType is the type of the channel (HTTP or SMPP)
	ChannelType string `gorm:"not null" json:"channel_type"`

	// Priority is the priority for the SMS type
	Priority int `gorm:"not null" json:"priority"`

	// TPS is the transactions per second limit for the channel
	TPS int `gorm:"not null" json:"tps"`

	// Status indicates whether the channel is active or inactive
	Status string `gorm:"not null" json:"status"`
}
