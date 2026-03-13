package domain

import "errors"

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrUserExists                   = errors.New("username already exists")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrUserAlreadyDeleted           = errors.New("user already deleted")
	ErrUserFrozen                   = errors.New("user is frozen")
	ErrUserAlreadyExists            = errors.New("user already exists") // 用于注册时
	ErrNewPasswordSameAsOldPassword = errors.New("new password cannot be the same as the old password")
)
