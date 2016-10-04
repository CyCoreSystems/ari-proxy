package session

func stringSliceEq(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for idx := range left {
		if left[idx] != right[idx] {
			return false
		}
	}

	return true
}
