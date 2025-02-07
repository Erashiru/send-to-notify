// Code generated by "stringer -type=Event"; DO NOT EDIT.

package models

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DeliveryOrderUpdate-1]
	_ = x[DeliveryOrderError-2]
	_ = x[ReserveUpdate-3]
	_ = x[ReserveError-4]
	_ = x[TableOrderUpdate-5]
	_ = x[TableOrderError-6]
	_ = x[StopListUpdate-7]
}

const _Event_name = "DeliveryOrderUpdateDeliveryOrderErrorReserveUpdateReserveErrorTableOrderUpdateTableOrderErrorStopListUpdate"

var _Event_index = [...]uint8{0, 19, 37, 50, 62, 78, 93, 107}

func (i Event) String() string {
	i -= 1
	if i < 0 || i >= Event(len(_Event_index)-1) {
		return "Event(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Event_name[_Event_index[i]:_Event_index[i+1]]
}
