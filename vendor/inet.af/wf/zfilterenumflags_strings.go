// Code generated by "stringer -output=zfilterenumflags_strings.go -type=filterEnumFlags -trimprefix=filterEnumFlags"; DO NOT EDIT.

package wf

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[filterEnumFlagsBestTerminatingMatch-1]
	_ = x[filterEnumFlagsSorted-2]
	_ = x[filterEnumFlagsBootTimeOnly-3]
	_ = x[filterEnumFlagsIncludeBootTime-4]
	_ = x[filterEnumFlagsIncludeDisabled-5]
}

const _filterEnumFlags_name = "BestTerminatingMatchSortedBootTimeOnlyIncludeBootTimeIncludeDisabled"

var _filterEnumFlags_index = [...]uint8{0, 20, 26, 38, 53, 68}

func (i filterEnumFlags) String() string {
	i -= 1
	if i >= filterEnumFlags(len(_filterEnumFlags_index)-1) {
		return "filterEnumFlags(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _filterEnumFlags_name[_filterEnumFlags_index[i]:_filterEnumFlags_index[i+1]]
}
