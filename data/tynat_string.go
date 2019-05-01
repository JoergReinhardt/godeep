// Code generated by "stringer -type=TyNat"; DO NOT EDIT.

package data

import "strconv"

const _TyNat_name = "NilBoolInt8Int16Int32IntBigIntUint8Uint16Uint32UintFlt32FloatBigFltRatioImag64ImagTimeDurationByteRuneBytesStringPipeBufferReaderWriterChannelSyncConSyncWaitErrorPairSliceMapLiteralDataExpressionFlagMASK"

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
	2097152:              _TyNat_name[102:107],
	4194304:              _TyNat_name[107:113],
	8388608:              _TyNat_name[113:117],
	16777216:             _TyNat_name[117:123],
	33554432:             _TyNat_name[123:129],
	67108864:             _TyNat_name[129:135],
	134217728:            _TyNat_name[135:142],
	268435456:            _TyNat_name[142:149],
	536870912:            _TyNat_name[149:157],
	1073741824:           _TyNat_name[157:162],
	2147483648:           _TyNat_name[162:166],
	4294967296:           _TyNat_name[166:171],
	8589934592:           _TyNat_name[171:174],
	17179869184:          _TyNat_name[174:181],
	34359738368:          _TyNat_name[181:185],
	68719476736:          _TyNat_name[185:195],
	137438953472:         _TyNat_name[195:199],
	18446744073709551615: _TyNat_name[199:203],
}

func (i TyNat) String() string {
	if str, ok := _TyNat_map[i]; ok {
		return str
	}
	return "TyNat(" + strconv.FormatInt(int64(i), 10) + ")"
}
