package order

import "errors"

var ErrAlreadyExists = errors.New("order repository: order already exists for user")

var ErrCreatedByAnotherUser = errors.New("order repository: order created by another user")