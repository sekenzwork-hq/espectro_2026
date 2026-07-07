package customerrors

type ServerError struct {
	OrgError string
}

func (v *ServerError) Error() string {
	return v.OrgError
}
