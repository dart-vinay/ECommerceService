package utils

import "errors"

var (
	ErrBadRequest = errors.New("Bad Request")
	ErrRequestDataInvalid = errors.New("Request Data is Invalid")
	ErrServerError = errors.New("Internal Server Error")
	ErrInvalidRequest = errors.New("Invalid Requests")
	ErrUserAlreadyExists = errors.New("User Already Exists")
	ErrUserDoesNotExist = errors.New("User Doesn't Exist")
	ErrUnauthorized = errors.New("Unauthorized User")
	ErrUserNotRegistered = errors.New("User Doesn't Exist")
	ErrBrandHandleAlreadyExist = errors.New("Brand Handle Already in Use.")
)

type ErrorResponse struct{
	StatusCode 			int	`json:"status_code"`
	ErrorMessage 		string `json:"error_message"`
}