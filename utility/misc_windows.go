package utility

// deviceMap currently only returns one device - how to obtain a device ID on windows ?
func DeviceMap(paths []string) map[uint64][]string {
	res := make(map[uint64][]string, 1)
	res[0] = paths
	return res
}
