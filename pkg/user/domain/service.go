package domain

type UserService struct {
	r UserRepository
}

func (us *UserService) FindUser(id UserID) (*User, error) {
	return us.r.Find(id)
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{r: r}
}
