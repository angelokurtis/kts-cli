package kubectl

func dedupeStr(a []string, b ...string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}

func dedupeInt(a []int, b ...int) []int {

	check := make(map[int]int)
	d := append(a, b...)
	res := make([]int, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}
