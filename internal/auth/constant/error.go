package constant

const (
	InvalidLoginCredentialsErrorMessage  = "invalid email or password"
	InvalidResetPasswordCodeErrorMessage = "invalid reset password code"
	InvalidTokenCredentialsErrorMessage  = "invalid combination of token or email"
	UnverifiedErrorMessage               = "your account is not verified"
	EmailNotExistsErrorMessage           = "email doesn't exists"
	InvalidEmailAlreadyExists            = "email already exists"
	TokenAlreadyExistsErrorMessage       = "a token was recently generated, please wait before requesting a new one"
	AlreadyVerifiedErrorMessage          = "your account has been verified"
	ExpiredTokenErrorMessage             = "the token has expired, please request a new one"
	InvalidQueryLimit                    = "invalid query limit"
	InvalidQueryPage                     = "invalid query page"
	InvalidQueryRole                     = "invalid query role"
	InvalidQueryIsAssign                 = "invalid query is_assign"
	InvalidUserId                        = "invalid user id or user id does not exists"
	BirthDateErrorMessage                = "birth date cannot be in the future"
	InvalidPhoneNumberAlreadyExists      = "phone number already exists"
	InvalidQueryisAssignAndRole          = "you cannot choose role 1 or 3 with is_assign 1 or 2"
	UserDetailNotFoundErrorMessage       = "user detail not found"
	AccountIsNotValidErrorMessage        = "account is not valid"
)
