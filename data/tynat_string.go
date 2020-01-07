// Code generated by "stringer -type=TyNat"; DO NOT EDIT.

package data

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nil-1]
	_ = x[Bool-2]
	_ = x[Int8-4]
	_ = x[Int16-8]
	_ = x[Int32-16]
	_ = x[Int-32]
	_ = x[BigInt-64]
	_ = x[Uint8-128]
	_ = x[Uint16-256]
	_ = x[Uint32-512]
	_ = x[Uint-1024]
	_ = x[Flt32-2048]
	_ = x[Float-4096]
	_ = x[BigFlt-8192]
	_ = x[Ratio-16384]
	_ = x[Imag64-32768]
	_ = x[Imag-65536]
	_ = x[Time-131072]
	_ = x[Duration-262144]
	_ = x[Byte-524288]
	_ = x[Rune-1048576]
	_ = x[Flag-2097152]
	_ = x[String-4194304]
	_ = x[Bytes-8388608]
	_ = x[Error-16777216]
	_ = x[Pair-33554432]
	_ = x[Slice-67108864]
	_ = x[Unboxed-134217728]
	_ = x[Map-268435456]
	_ = x[Function-536870912]
	_ = x[Literal-1073741824]
	_ = x[Type-2147483648]
	_ = x[MASK-18446744073709551615]
}

const _TyNat_name = "NilBoolInt8Int16Int32IntBigIntUint8Uint16Uint32UintFlt32FloatBigFltRatioImag64ImagTimeDurationByteRuneFlagStringBytesErrorPairSliceUnboxedMapFunctionLiteralTypeMASK"

var _TyNat_map = map[TyNat]string{
	1:                    _TyNat_name[0:3],
	2:                    _TyNat_name[3:7],
	4:                    _TyNat_name[7:11],
	8:                    _TyNat_name[11:16],
	16:                   _TyNat_name[16:21],
	32:                   _TyNat_name[21:24],
	64:                   _TyNat_name[24:30],
	128:                  _TyNat_name[30:35],
	256:                  _TyNat_name[35:41],
	512:                  _TyNat_name[41:47],
	1024:                 _TyNat_name[47:51],
	2048:                 _TyNat_name[51:56],
	4096:                 _TyNat_name[56:61],
	8192:                 _TyNat_name[61:67],
	16384:                _TyNat_name[67:72],
	32768:                _TyNat_name[72:78],
	65536:                _TyNat_name[78:82],
	131072:               _TyNat_name[82:86],
	262144:               _TyNat_name[86:94],
	524288:               _TyNat_name[94:98],
	1048576:              _TyNat_name[98:102],
	2097152:              _TyNat_name[102:106],
	4194304:              _TyNat_name[106:112],
	8388608:              _TyNat_name[112:117],
	16777216:             _TyNat_name[117:122],
	33554432:             _TyNat_name[122:126],
	67108864:             _TyNat_name[126:131],
	134217728:            _TyNat_name[131:138],
	268435456:            _TyNat_name[138:141],
	536870912:            _TyNat_name[141:149],
	1073741824:           _TyNat_name[149:156],
	2147483648:           _TyNat_name[156:160],
	18446744073709551615: _TyNat_name[160:164],
}

func (i TyNat) String() string {
	if str, ok := _TyNat_map[i]; ok {
		return str
	}
	return "TyNat(" + strconv.FormatInt(int64(i), 10) + ")"
}
