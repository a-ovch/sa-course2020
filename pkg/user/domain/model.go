package domain

type UserID int

type User struct {
	Id        UserID `json:"-"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UserRepository interface {
	Find(id UserID) (*User, error)
}
