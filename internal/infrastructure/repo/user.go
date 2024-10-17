package repo

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
)

type UserRepository interface {
	RefillBalance(ctx context.Context, userID int, amount int64) error
	GetBalance(ctx context.Context, userID int) (*int64, error)
	GetAllUsers(ctx context.Context) ([]User, error)
}

type UserRepo struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) RefillBalance(ctx context.Context, userID int, amount int64) error {
	user := &User{}
	err := r.db.Model(user).Where("id = ?", userID).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	_, err = r.db.Model(user).Set("balance = balance + ?", amount).Where("id = ?", userID).Update()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetBalance(ctx context.Context, userID int) (*int64, error) {
	var user User
	err := r.db.Model(&user).
		Column("balance").
		Where("id = ?", userID).
		Select()
	if err != nil {
		return nil, err
	}
	return user.Balance, nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}
