package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	UserNameOrPwdError      = newError(1010001, "Role or password error")
	ErrUserNotExist         = newError(1010002, "User does not exist")
	ErrUserAlreadyExist     = newError(1010003, "User already exists")
	ErrCaptchaIdNotFound    = newError(1010004, "captcha id not found")
	ErrMenuCodeAlreadyExist = newError(1011001, "Menu Code already exists")
	ErrEmailAlreadyUse      = newError(1001, "The email is already in use.")
)
