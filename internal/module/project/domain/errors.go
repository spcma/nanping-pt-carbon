package domain

import "errors"

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
	ErrProjectCodeInvalid   = errors.New("invalid project code")
	ErrProjectStatusInvalid = errors.New("invalid project status")

	// Project Member errors
	ErrProjectMemberNotFound      = errors.New("project member not found")
	ErrProjectMemberAlreadyExists = errors.New("project member already exists")
	ErrProjectMemberRoleInvalid   = errors.New("invalid project member role")
)
