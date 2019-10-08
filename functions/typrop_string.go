// Code generated by "stringer -type TyProp"; DO NOT EDIT.

package functions

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Default-0]
	_ = x[PostFix-1]
	_ = x[InFix-3]
	_ = x[Atomic-4]
	_ = x[Eager-5]
	_ = x[RightBound-6]
	_ = x[Mutable-7]
	_ = x[SideEffect-8]
	_ = x[Primitive-9]
}

const (
	_TyProp_name_0 = "DefaultPostFix"
	_TyProp_name_1 = "InFixAtomicEagerRightBoundMutableSideEffectPrimitive"
)

var (
	_TyProp_index_0 = [...]uint8{0, 7, 14}
	_TyProp_index_1 = [...]uint8{0, 5, 11, 16, 26, 33, 43, 52}
)

func (i TyProp) String() string {
	switch {
	case 0 <= i && i <= 1:
		return _TyProp_name_0[_TyProp_index_0[i]:_TyProp_index_0[i+1]]
	case 3 <= i && i <= 9:
		i -= 3
		return _TyProp_name_1[_TyProp_index_1[i]:_TyProp_index_1[i+1]]
	default:
		return "TyProp(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
