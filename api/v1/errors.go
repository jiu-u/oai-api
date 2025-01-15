package v1

import "net/http"

var (
	// common errors
	ErrSuccess             = newError(http.StatusOK, 0, "ok")
	ErrBadRequest          = newError(http.StatusBadRequest, 400, "Bad Request")
	ErrUnauthorized        = newError(http.StatusUnauthorized, 401, "Unauthorized")
	ErrNotFound            = newError(http.StatusNotFound, 404, "Not Found")
	ErrInternalServerError = newError(http.StatusInternalServerError, 500, "Internal Server Error")

	// more biz errors
	//UserNameOrPwdError      = newError(1010001, "Role or password error")
	//ErrUserNotExist         = newError(1010002, "User does not exist")
	//ErrUserAlreadyExist     = newError(1010003, "User already exists")
	//ErrCaptchaIdNotFound    = newError(1010004, "captcha id not found")
	//ErrMenuCodeAlreadyExist = newError(1011001, "Menu Code already exists")
	//ErrEmailAlreadyUse      = newError(1001, "The email is already in use.")
)
