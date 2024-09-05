package user

import "fmt"

type userAlreadyExistsError struct {
	err   error
	login string
}

func NewUserAlreadyExistsError(login string, err error) *userAlreadyExistsError {
	return &userAlreadyExistsError{
		login: login,
		err:   err,
	}
}

func (u *userAlreadyExistsError) Error() string {
	return fmt.Sprintf("login %s already exists: %s", u.login, u.err)
}

func (u *userAlreadyExistsError) Unwrap() error {
	return u.err
}

type UserNotFoundError struct {
	login string
	err   error
}

func NewUserNotFoundError(login string, err error) *UserNotFoundError {
	return &UserNotFoundError{
		login: login,
		err:   err,
	}
}

func (unf *UserNotFoundError) Error() string {
	return fmt.Sprintf("user \"%s\" not found: %s", unf.login, unf.err)
}

func (unf *UserNotFoundError) Unwrap() error {
	return unf.err
}