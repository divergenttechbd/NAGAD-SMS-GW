package models

// User represents a user in the system
// @Description Represents a user in the system with roles and campaigns
type User struct {
	BaseModel
	// Username is a unique identifier for the user
	Username string `gorm:"unique;not null" json:"username"`

	// Email is the user's email address
	Email string `gorm:"unique;not null" json:"email"`

	// Password is the user's hashed password
	Password string `gorm:"not null" json:"-"`

	// Roles represents the roles assigned to the user
	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`

	// Campaigns represents the campaigns associated with the user
	Campaigns []Campaign `gorm:"foreignKey:UserID" json:"campaigns"`
}
