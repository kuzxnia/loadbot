// Code generated by "stringer -type=CommandType -trimprefix=CommandType"; DO NOT EDIT.

package database

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CommandTypeStartWorkload-0]
	_ = x[CommandTypeStopWorkload-1]
}

const _CommandType_name = "StartWorkloadStopWorkload"

var _CommandType_index = [...]uint8{0, 13, 25}

func (i CommandType) String() string {
	if i < 0 || i >= CommandType(len(_CommandType_index)-1) {
		return "CommandType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CommandType_name[_CommandType_index[i]:_CommandType_index[i+1]]
}
