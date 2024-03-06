// Code generated by "stringer -type=CommandState -trimprefix=CommandState"; DO NOT EDIT.

package database

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CommandStateCreated-0]
	_ = x[CommandStateRunning-1]
	_ = x[CommandStateDone-2]
	_ = x[CommandStateError-3]
}

const _CommandState_name = "CreatedRunningDoneError"

var _CommandState_index = [...]uint8{0, 7, 14, 18, 23}

func (i CommandState) String() string {
	if i < 0 || i >= CommandState(len(_CommandState_index)-1) {
		return "CommandState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CommandState_name[_CommandState_index[i]:_CommandState_index[i+1]]
}