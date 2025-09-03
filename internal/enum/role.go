package enum

type Role int

const (
	_ Role = iota
	RoleMember
	RoleAdmin
)

var roleState = map[Role]string{
	RoleMember: "member",
	RoleAdmin:  "admin",
}

func (s Role) String() string {
	return roleState[s]
}
