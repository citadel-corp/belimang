package user

var (
	MinUsername = 5
	MaxUsername = 30

	MinPassword = 5
	MaxPassword = 30
)

type UserType string

const (
	Admin UserType = "Admin"
	User  UserType = "User"
)

var UserTypes = []interface{}{Admin, User}

type Users struct {
	ID             uint64
	UID            string
	Username       string
	Email          string
	HashedPassword string
	UserType       UserType
}
