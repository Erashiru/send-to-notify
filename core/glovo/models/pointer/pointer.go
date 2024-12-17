package pointer

func OfBool(b bool) *bool {
	return &b
}

func OfFloat64(f float64) *float64 {
	return &f
}

func OfInt(d int) *int {
	return &d
}
