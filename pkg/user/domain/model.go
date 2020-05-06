package domain

import "sa-course-app/pkg/common/domain"

type UserID domain.UUID

type User struct {
	id        UserID
	username  string
	firstName string
	lastName  string
	email     string
	phone     string
}

type UserRepository interface {
	NextID() UserID
	Store(u *User) error
	Find(id UserID) (*User, error)
}

func (u *User) GetID() UserID {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetFirstName() string {
	return u.firstName
}

func (u *User) GetLastName() string {
	return u.lastName
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPhone() string {
	return u.phone
}

func NewUser(id UserID, username string, firstName string, lastName string, email string, phone string) *User {
	return &User{
		id:        id,
		username:  username,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		phone:     phone,
	}
}
