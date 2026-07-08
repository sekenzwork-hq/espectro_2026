package customerrors

type CredentialsError struct {
	OrgError string
}

func (c *CredentialsError) Error() string {
	return c.OrgError
}
