package user

type users map[string]User

type Repository struct {
	users users
}

func NewRepository(users *map[string]User) *Repository {
	return &Repository{users: *users}
}

func (r *Repository) Save(user User) User {
	r.users[user.ID] = user

	return user
}

func (r *Repository) Update(user User) User {
	r.users[user.ID] = user

	return user
}

func (r *Repository) FindOne(id string) User {
	return r.users[id]
}

func (r *Repository) All() map[string]User {
	return r.users
}

//func (r *Repository) find(query string) []User {
//	return r.users
//}
