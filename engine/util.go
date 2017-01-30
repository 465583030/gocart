package engine

type (
	form interface {
		Validate(Validator) error
	}
)

func validateForm(f form, v Validator) error {
	return f.Validate(v)
}

func idFilter(id uint) []*Filter {
	return []*Filter{NewFilter("id", Equal, id)}
}

func boolPtr(v bool) *bool {
	return &v
}
