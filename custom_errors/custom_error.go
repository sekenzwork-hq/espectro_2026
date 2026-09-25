package customerrors

var CredentialErr *CredentialsError
var PermissionErr *PermissionError
var ValidationErr *ValidationError
var ServerErr *ServerError
var AuthErr *AuthenticationError
var NotFoundErr *NotFoundError
var InvalidAdminUpdateModeErr *InvalidAdminUpdateModeError
var NotFoundOrLeaderErr *NotFoundOrLeaderError
var SizeErr *SizeError

type PermissionError struct {
	DisplayError string
}

func (p PermissionError) Error() string {
	return p.DisplayError
}

type ValidationError struct {
	DisplayError string
}

func (v *ValidationError) Error() string {
	return v.DisplayError
}

type CredentialsError struct {
	DisplayError string
}

func (c *CredentialsError) Error() string {
	return c.DisplayError
}

type ServerError struct {
	DisplayError string
}

func (v *ServerError) Error() string {
	return v.DisplayError
}

func (v *ServerError) OriginalError() string {
	return v.DisplayError
}

type AuthenticationError struct {
	DisplayError string
}

func (a AuthenticationError) Error() string {
	return a.DisplayError
}

type NotFoundError struct {
	DisplayError string
}

func (n *NotFoundError) Error() string {
	return n.DisplayError
}
func (n *NotFoundError) OriginalError() string {
	return n.DisplayError
}

type InvalidAdminUpdateModeError struct {
	DisplayError string
}

func (i InvalidAdminUpdateModeError) Error() string {
	return i.DisplayError
}

type NotFoundOrLeaderError struct {
	DisplayError string
}

func (i NotFoundOrLeaderError) Error() string {
	return i.DisplayError
}

type SizeError struct {
	DisplayError string
}

func (i SizeError) Error() string {
	return i.DisplayError
}
