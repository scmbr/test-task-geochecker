package v1

func maskString(s string) string {
	r := []rune(s)
	n := len(r)
	const mask = "****"
	return string(r[0]) + mask + string(r[n-1])
}
