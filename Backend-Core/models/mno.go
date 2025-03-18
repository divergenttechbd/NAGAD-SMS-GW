package models

// MNOConfiguration represents the configuration for a Mobile Network Operator (MNO)
// @Description Represents the configuration for an MNO, including prefix, channels, and TPS limits
type MNO struct {
	BaseModel
	// MNO_Name is the name of the Mobile Network Operator (e.g., GP, BL, RB)
	MNO_Name string `gorm:"not null" json:"mno_name"`

	// Prefix is the phone number prefix associated with the MNO (e.g., 017, 019)
	Prefix string `gorm:"not null" json:"prefix"`

	// Channels represents the delivery channels associated with the MNO
	Channels []MnoChannels `gorm:"foreignKey:MNOID" json:"channels"`

	// Status indicates whether the MNO configuration is active or inactive
	Status string `gorm:"not null" json:"status"`
}
