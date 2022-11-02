package utils

import "errors"

var (
	IdNotFoundErr = errors.New("id not found in group")

	InvalidOptionsErr = errors.New("current options is invalid")

	InvalidArgsErr = errors.New("invalid argument")

	DuplicatedIdErr = errors.New("id is duplicated")

	OutOfMaxCntErr = errors.New("numbers out of limit")

	TimeOutErr = errors.New("timeout")

	DeleteObjectFailed = errors.New("delete object in map failed")
)
