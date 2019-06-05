// Code generated by "stringer -type Arity"; DO NOT EDIT.

package functions

import "strconv"

const _Arity_name = "NullaryUnaryBinaryTernaryQuaternaryQuinarySenarySeptenaryOctonaryNovenaryDenary"

var _Arity_index = [...]uint8{0, 7, 12, 18, 25, 35, 42, 48, 57, 65, 73, 79}

func (i Arity) String() string {
	if i < 0 || i >= Arity(len(_Arity_index)-1) {
		return "Arity(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Arity_name[_Arity_index[i]:_Arity_index[i+1]]
}
