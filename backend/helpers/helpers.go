package helpers

func Max(vars ...float64) float64 {

	if len(vars) == 0 {
		return 0
	}

	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}
