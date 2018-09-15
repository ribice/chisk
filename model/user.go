package chisk

import "time"

// User struct
type User struct {
	Base
	Email              string     `json:"email"`
	DisplayName        string     `json:"display_name"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	Password           string     `json:"-"`
	PhoneNumber        string     `json:"phone_number,omitempty"`
	Token              string     `json:"-"`
	IsActive           bool       `json:"is_active"`
	Role               AccessRole `json:"-"`
	LastPasswordChange *time.Time `json:"last_password_change,omitempty"`
}

// ChangePassword changes user's password
func (u *User) ChangePassword(p string) {
	u.Password = p
	t := time.Now()
	u.LastPasswordChange = &t
}

// FullName returns user's full name, firstName + " " + lastName
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// AuthUser converts User struct to AuthUser
func (u *User) AuthUser() *AuthUser {
	return &AuthUser{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Role:        u.Role,
	}
}

// AuthUser represents data stored in session/context for a user
type AuthUser struct {
	ID          string
	DisplayName string
	Email       string
	Role        AccessRole
}
