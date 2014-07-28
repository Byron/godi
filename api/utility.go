package api

// Append elm if it is not yet on dest
func AppendUniqueString(dest []string, elm string) []string {
	for _, d := range dest {
		if d == elm {
			return dest
		}
	}
	return append(dest, elm)
}
