package mysql

import (
	"database/sql"
	"github.com/google/uuid"

	"sa-course-app/pkg/infrastructure/database"
	"sa-course-app/pkg/user/domain"
)

type userRepository struct {
	client database.Client
}

func (ur *userRepository) NextID() domain.UserID {
	return domain.UserID(uuid.New())
}

func (ur *userRepository) Store(u *domain.User) error {
	const query = "INSERT INTO user (id, username, first_name, last_name, phone, email) VALUES (?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE username=VALUES(username), first_name=VALUES(first_name), last_name=VALUES(last_name), phone=VALUES(phone), email=VALUES(email)"

	binaryUUID, err := binaryUserID(u.GetID())
	if err == nil {
		_, err = ur.client.Exec(
			query,
			binaryUUID,
			u.GetUsername(),
			u.GetFirstName(),
			u.GetLastName(),
			u.GetPhone(),
			u.GetEmail(),
		)
	}

	return err
}

func (ur *userRepository) Find(id domain.UserID) (*domain.User, error) {
	const query = "SELECT id, username, first_name, last_name, email, phone FROM user WHERE id = ?"

	binaryUUID, err := binaryUserID(id)
	if err != nil {
		return nil, err
	}

	row := ur.client.QueryRow(query, binaryUUID)
	return hydrateUserFromRow(row)
}

func (ur *userRepository) Delete(id domain.UserID) error {
	const query = "DELETE FROM user WHERE id = ?"

	binaryUUID, err := binaryUserID(id)
	if err == nil {
		_, err = ur.client.Exec(query, binaryUUID)
	}

	return err
}

func NewUserRepository(client database.Client) domain.UserRepository {
	return &userRepository{client: client}
}

func binaryUserID(id domain.UserID) ([]byte, error) {
	return uuid.UUID(id).MarshalBinary()
}

func hydrateUserFromRow(row *sql.Row) (*domain.User, error) {
	var id uuid.UUID
	var username, firstName, lastName, email, phone string

	err := row.Scan(&id, &username, &firstName, &lastName, &email, &phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return domain.NewUser(
		domain.UserID(id),
		username,
		firstName,
		lastName,
		email,
		phone,
	), nil
}
