// Code generated by "stringer -type TyFlag"; DO NOT EDIT.

package functions

import "strconv"

const _TyFlag_name = "Flag_BitFlag"

var _TyFlag_index = [...]uint8{0, 12}

func (i TyFlag) String() string {
	if i >= TyFlag(len(_TyFlag_index)-1) {
		return "TyFlag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TyFlag_name[_TyFlag_index[i]:_TyFlag_index[i+1]]
}
