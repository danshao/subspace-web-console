package business

import "errors"

type UserModel interface {
	GetUserId() int
	GetRole() string
	IsEnabled() bool
}

func CanSignIn(user UserModel) (bool, error) {
	if !user.IsEnabled() {
		return false, errors.New("User is disabled. Please contact an administrator for help.")
	}

	if "admin" != user.GetRole() {
		return false, errors.New("User is not an administrator. Please contact an administrator for help.")
	}

	return true, nil
}