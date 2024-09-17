package order

import "errors"

var ErrAlreadyExists = errors.New("order already exists for user")

var ErrCreatedByAnotherUser = errors.New("order created by another user")