package app

import (
	"errors"
	"sa-course-app/pkg/user/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserData struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Service struct {
	repo domain.UserRepository
}

func (s *Service) CreateUser(ud *UserData) (domain.UserID, error) {
	id := s.repo.NextID()
	u := domain.NewUser(id, ud.Username, ud.FirstName, ud.LastName, ud.Email, ud.Phone)
	return id, s.repo.Store(u)
}

func (s *Service) FindUser(id domain.UserID) (*UserData, error) {
	u, err := s.repo.Find(id)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrUserNotFound
	}

	return &UserData{
		Username:  u.GetUsername(),
		FirstName: u.GetFirstName(),
		LastName:  u.GetLastName(),
		Email:     u.GetEmail(),
		Phone:     u.GetPhone(),
	}, nil
}

func (s *Service) UpdateUser(id domain.UserID, data *UserData) error {
	u, err := s.repo.Find(id)
	if err != nil {
		return err
	}

	if u == nil {
		return ErrUserNotFound
	}

	u.SetUsername(data.Username)
	u.SetFirstName(data.FirstName)
	u.SetLastName(data.LastName)
	u.SetEmail(data.Email)
	u.SetPhone(data.Phone)

	return s.repo.Store(u)
}

func (s *Service) DeleteUser(id domain.UserID) error {
	return s.repo.Delete(id)
}

func NewService(r domain.UserRepository) *Service {
	return &Service{repo: r}
}
