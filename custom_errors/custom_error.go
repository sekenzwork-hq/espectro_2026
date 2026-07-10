package customerrors

var CredentialErr *CredentialsError
var PermissionErr *PermissionError
var ValidationErr *ValidationError
var ServerErr *ServerError
var AuthErr *AuthenticationError
var NotFoundErr *NotFoundError
var InvalidAdminUpdateModeErr *InvalidAdminUpdateModeError
var NotFoundOrLeaderErr *NotFoundOrLeaderError
var SpaceErr *SpaceError

type PermissionError struct {
	OrgError string
}

func (p PermissionError) Error() string {
	return p.OrgError
}

type ValidationError struct {
	OrgError string
}

func (v *ValidationError) Error() string {
	return v.OrgError
}

type CredentialsError struct {
	OrgError string
}

func (c *CredentialsError) Error() string {
	return c.OrgError
}

type ServerError struct {
	OrgError string
}

func (v *ServerError) Error() string {
	return v.OrgError
}

type AuthenticationError struct {
	OrgError string
}

func (a AuthenticationError) Error() string {
	return a.OrgError
}

type NotFoundError struct {
	OrgError string
}

func (n *NotFoundError) Error() string {
	return n.OrgError
}

type InvalidAdminUpdateModeError struct {
	OrgError string
}

func (i InvalidAdminUpdateModeError) Error() string {
	return i.OrgError
}

type NotFoundOrLeaderError struct {
	OrgError string
}

func (i NotFoundOrLeaderError) Error() string {
	return i.OrgError
}

type SpaceError struct {
	OrgError string
}

func (i SpaceError) Error() string {
	return i.OrgError
}
