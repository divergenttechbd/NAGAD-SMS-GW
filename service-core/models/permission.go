package models

// Permission represents an action a role can perform
// @Description A permission assigned to roles
type Permission struct {
	BaseModel
	// Name is the unique name of the permission
	// @Property name string true "Permission Name" example("view_users")
	Name string `gorm:"unique;not null" json:"name"`
}
