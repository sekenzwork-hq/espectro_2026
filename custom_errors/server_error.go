package customerrors

type ServerError struct {
	OrgError string
}

func (v *ServerError) Error() string {
	return "Something went wrong"
}
