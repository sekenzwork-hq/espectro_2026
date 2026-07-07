package customerrors

type ValidationError struct {
	OrgError string
}

func (v *ValidationError) Error() string {
	return "Validation Error"
}
