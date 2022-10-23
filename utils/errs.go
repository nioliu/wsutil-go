package utils

import "errors"

var (
	IdNotFound = errors.New("id not found in group")

	InvalidOptions = errors.New("current options is invalid")

	InvalidArgs = errors.New("invalid argument")

	DuplicatedId = errors.New("id is duplicated")

	OutOfMaxCnt = errors.New("numbers out of limit")
)
