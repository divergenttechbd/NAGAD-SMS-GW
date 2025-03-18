package models

// Role represents a user role
// @Description Role assigned to users with permissions
type Role struct {
	BaseModel
	// Name is the unique name of the role
	// @Property name string true "Role Name" example("admin")
	Name string `gorm:"unique;not null" json:"name"`

	// Permissions represents the permissions assigned to the role
	// @Property permissions array "Permissions" items={Permission} example([{"id":1,"name":"view_users"}])
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
