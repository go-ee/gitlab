package cmd

func SliceToMap(values []string) (ret map[string]bool) {
	ret = make(map[string]bool)
	if values != nil {
		for _, name := range values {
			ret[name] = true
		}
	}
	return
}
