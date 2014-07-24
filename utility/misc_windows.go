package utility

// deviceMap currently only returns one device - how to obtain a device ID on windows ?
func DeviceMap(paths []string) [][]string {
	res := make([][]string, 1)
	res[0] = paths
	return res
}
