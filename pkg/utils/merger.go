package utils

func MergeMaps(maps ...map[string]string) map[string]string {
	mergedMap := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			mergedMap[k] = v
		}
	}
	return mergedMap
}

func MergeSlices(slices ...[]string) []string {
	var mergedSlice []string
	for _, s := range slices {
		mergedSlice = append(mergedSlice, s...)
	}
	return mergedSlice
}
