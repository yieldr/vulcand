package auth

type mockUser struct {
	username, fullName, email string
	accounts, roles           []string
}

func NewMockUser(username, fullName, email string, accounts, roles []string) User {
	return mockUser{username, fullName, email, accounts, roles}
}

func (u mockUser) Username() string {
	return u.username
}

func (u mockUser) FullName() string {
	return u.fullName
}

func (u mockUser) Email() string {
	return u.email
}

func (u mockUser) Accounts() []string {
	return u.accounts
}

func (u mockUser) Roles() []string {
	return u.roles
}
