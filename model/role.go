package chisk

// AccessRole represents access role type
type AccessRole int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = iota + 1

	// AdminRole is an admin role
	AdminRole

	// UserRole is a standard user role
	UserRole
)
