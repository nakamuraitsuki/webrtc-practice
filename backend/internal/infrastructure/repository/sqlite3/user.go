package sqlite3

import (
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.IUserRepository {
	return &UserRepository{db: db}
}

type UserModel struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Email      string `db:"email"`
	PasswdHash string `db:"passwd_hash"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (ur *UserRepository) GetAllUsers() ([]*entity.User, error) {
	var users []UserModel
	err := ur.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var res []*entity.User
	for _, user := range users {
		res = append(res, entity.NewUser(
			user.ID,
			user.Name,
			user.Email,
			user.PasswdHash,
			user.CreatedAt,
			user.UpdatedAt,
		))
	}

	return res, nil
}

func (ur *UserRepository) GetUserByID(id int) (*entity.User, error) {
	var user UserModel
	err := ur.db.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return entity.NewUser(
		user.ID,
		user.Name,
		user.Email,
		user.PasswdHash,
		user.CreatedAt,
		user.UpdatedAt,
	), nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user UserModel
	err := ur.db.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return entity.NewUser(user.ID, user.Name, user.Email, user.PasswdHash, user.CreatedAt, user.UpdatedAt), nil
}

func (ur *UserRepository) CreateUser(params repository.CreateUserParams) (*entity.User, error) {
	created_at := time.Now()
	_, err := ur.db.Exec("INSERT INTO users (name, email, passwd_hash, created_at) VALUES (?, ?, ?, ?)", params.Name, params.Email, params.PasswdHash, created_at)
	if err != nil {
		return nil, err
	}

	res, err := ur.GetUserByEmail(params.Email)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ur *UserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	_, err := ur.db.Exec(`UPDATE users SET name = ?, email = ?, passwd_hash = ?, updated_at = ? WHERE id = ?`,
	 user.GetName(), 
	 user.GetEmail(), 
	 user.GetPasswdHash(), 
	 time.Now(), 
	 user.GetID(),
	)
	if err != nil {
		return nil, err
	}

	res, err := ur.GetUserByID(user.GetID())

	return res, nil
}

func (ur *UserRepository) DeleteUser(id int) error {
	_, err := ur.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
