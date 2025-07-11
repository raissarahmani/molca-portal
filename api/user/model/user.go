package model

import "github.com/go-playground/validator/v10"

type User struct {
	ID string `json:"id"`
}

func (user *User) GetValue() *User {
	return user
}

func (user *User) GetID() string {
	return user.ID
}

func NewUser(id string) (*User, error) {

	u := User{
		ID: id,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (user *User) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}
