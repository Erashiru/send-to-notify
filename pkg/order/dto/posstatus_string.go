// Code generated by "stringer -type=PosStatus"; DO NOT EDIT.

package dto

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ACCEPTED-1]
	_ = x[NEW-2]
	_ = x[WAIT_COOKING-3]
	_ = x[READY_FOR_COOKING-4]
	_ = x[COOKING_COMPLETE-5]
	_ = x[COOKING_STARTED-6]
	_ = x[CLOSED-7]
	_ = x[READY_FOR_PICKUP-8]
	_ = x[PICKED_UP_BY_CUSTOMER-9]
	_ = x[ON_WAY-10]
	_ = x[OUT_FOR_DELIVERY-11]
	_ = x[DELIVERED-12]
	_ = x[CANCELLED_BY_POS_SYSTEM-13]
	_ = x[PAYMENT_NEW-14]
	_ = x[PAYMENT_IN_PROGRESS-15]
	_ = x[PAYMENT_SUCCESS-16]
	_ = x[PAYMENT_CANCELLED-17]
	_ = x[PAYMENT_WAITING-18]
	_ = x[PAYMENT_DELETED-19]
	_ = x[FAILED-20]
	_ = x[WAIT_SENDING-21]
}

const _PosStatus_name = "ACCEPTEDNEWWAIT_COOKINGREADY_FOR_COOKINGCOOKING_COMPLETECOOKING_STARTEDCLOSEDREADY_FOR_PICKUPPICKED_UP_BY_CUSTOMERON_WAYOUT_FOR_DELIVERYDELIVEREDCANCELLED_BY_POS_SYSTEMPAYMENT_NEWPAYMENT_IN_PROGRESSPAYMENT_SUCCESSPAYMENT_CANCELLEDPAYMENT_WAITINGPAYMENT_DELETEDFAILEDWAIT_SENDING"

var _PosStatus_index = [...]uint16{0, 8, 11, 23, 40, 56, 71, 77, 93, 114, 120, 136, 145, 168, 179, 198, 213, 230, 245, 260, 266, 278}

func (i PosStatus) String() string {
	i -= 1
	if i < 0 || i >= PosStatus(len(_PosStatus_index)-1) {
		return "PosStatus(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _PosStatus_name[_PosStatus_index[i]:_PosStatus_index[i+1]]
}
