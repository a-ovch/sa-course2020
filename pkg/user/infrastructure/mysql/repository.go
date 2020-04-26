package mysql

import (
	"database/sql"
	"sa-course-app/pkg/infrastructure/database"
	"sa-course-app/pkg/user/domain"
)

type userRepository struct {
	client database.Client
}

func (ur *userRepository) Find(id domain.UserID) (*domain.User, error) {
	const query = "SELECT id, username, first_name, last_name, email, phone FROM user WHERE id = ?"
	row := ur.client.QueryRow(query, int(id))
	return hydrateUserFromRow(row)
}

func NewUserRepository(client database.Client) domain.UserRepository {
	return &userRepository{client: client}
}

func hydrateUserFromRow(row *sql.Row) (*domain.User, error) {
	var (
		id                                          int
		username, firstName, lastName, email, phone string
	)

	err := row.Scan(&id, &username, &firstName, &lastName, &email, &phone)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        domain.UserID(id),
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}, nil
}
