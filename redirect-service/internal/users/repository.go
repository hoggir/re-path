package users

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetAll() []string {
	return []string{"Alice", "Bob", "Charlie"}
}
