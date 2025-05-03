package rest

import (
	"contacts-list/internal/domain/errs"
	"contacts-list/pkg/sl"
)

func NewUnmarshalError(err error) error {
	return &errs.AppError{
		Underlying: err,
		Message:    err.Error(),
		Code:       errs.CodeBadInput,
		Level:      sl.LevelInfo,
		Additional: nil,
	}
}
