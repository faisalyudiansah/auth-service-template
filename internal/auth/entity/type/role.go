package custom_type

type Role uint

const (
	RoleUnknown Role = iota
	RoleAdmin
	RoleUser
)

func (s Role) String() string {
	return [...]string{"Unknown", "Admin", "User"}[s]
}

func (s Role) EnumIndex() uint {
	return uint(s)
}

func (s Role) IsRoleUnknown() bool {
	return s == RoleUnknown
}

func (s Role) IsRoleAdmin() bool {
	return s == RoleAdmin
}

func (s Role) IsRoleUser() bool {
	return s == RoleUser
}
