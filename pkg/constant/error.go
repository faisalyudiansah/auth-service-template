package constant

const (
	EntityNotFoundErrorMessage        = "data not found : %s"
	SessionExpiredMessage             = "session expired"
	RequestTimeoutErrorMessage        = "failed to process request in time, please try again"
	UnauthorizedErrorMessage          = "unauthorized to access this resource"
	ExpiredTokenErrorMessage          = "access token expired"
	EOFErrorMessage                   = "missing body request"
	StrConvSyntaxErrorMessage         = "invalid syntax for %s"
	DontHavePermissionErrorMessage    = "you do not have permission to access this resource"
	ForbiddenAccessErrorMessage       = "forbidden access"
	TooManyRequestsErrorMessage       = "the server is experiencing high load, please try again later"
	InternalServerErrorMessage        = "internal server error"
	ResetPasswordErrorMessage         = "please try again later"
	ValidationErrorMessage            = "input validation error"
	InvalidJsonUnmarshallErrorMessage = "invalid JSON format"
	JsonSyntaxErrorMessage            = "invalid JSON syntax"
	InvalidJsonValueTypeErrorMessage  = "invalid value for %s"
	InvalidIDErrorMessage             = "expected a numeric value"
	ClientDefaultErrorMessage         = "bad request"
)
