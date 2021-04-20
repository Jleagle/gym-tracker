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

func Mean(vars ...float64) float64 {

	if len(vars) == 0 {
		return 0
	}

	var count float64
	var sum float64

	for _, i := range vars {
		count += 1
		sum += i
	}
	return sum / count
}

func ReplaceAtIndex(in string, replacement rune, index int) string {
	out := []rune(in)
	out[index] = replacement
	return string(out)
}
