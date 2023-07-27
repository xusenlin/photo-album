package utils

func SpreadDigit(n int) []int {
	var r []int
	for i := 1; i <= n; i++ {
		r = append(r, i)
	}
	return r
}
