// Code generated by "stringer -type=Type"; DO NOT EDIT.

package types

import "strconv"

const _Type_name = "NilBoolIntInt8Int16Int32BigIntUintUint8Uint16Uint32FloatFlt32BigFltRatioImagImag64ByteRuneBytesStringTimeDurationAttrErrorTupleListChainAtListUniSetAtSetRecordLinkDLinkNodeTreeFunctionFlagDataTypeNodeTypeTokenTypeMetaType"

var _Type_map = map[Type]string{
	0:             _Type_name[0:3],
	1:             _Type_name[3:7],
	4:             _Type_name[7:10],
	8:             _Type_name[10:14],
	16:            _Type_name[14:19],
	32:            _Type_name[19:24],
	64:            _Type_name[24:30],
	128:           _Type_name[30:34],
	256:           _Type_name[34:39],
	512:           _Type_name[39:45],
	1024:          _Type_name[45:51],
	2048:          _Type_name[51:56],
	4096:          _Type_name[56:61],
	8192:          _Type_name[61:67],
	16384:         _Type_name[67:72],
	32768:         _Type_name[72:76],
	65536:         _Type_name[76:82],
	131072:        _Type_name[82:86],
	262144:        _Type_name[86:90],
	524288:        _Type_name[90:95],
	1048576:       _Type_name[95:101],
	2097152:       _Type_name[101:105],
	4194304:       _Type_name[105:113],
	8388608:       _Type_name[113:117],
	16777216:      _Type_name[117:122],
	33554432:      _Type_name[122:127],
	67108864:      _Type_name[127:131],
	134217728:     _Type_name[131:136],
	268435456:     _Type_name[136:142],
	536870912:     _Type_name[142:148],
	1073741824:    _Type_name[148:153],
	2147483648:    _Type_name[153:159],
	4294967296:    _Type_name[159:163],
	8589934592:    _Type_name[163:168],
	17179869184:   _Type_name[168:172],
	34359738368:   _Type_name[172:176],
	68719476736:   _Type_name[176:184],
	137438953472:  _Type_name[184:188],
	274877906944:  _Type_name[188:196],
	549755813888:  _Type_name[196:204],
	1099511627776: _Type_name[204:213],
	2199023255552: _Type_name[213:221],
}

func (i Type) String() string {
	if str, ok := _Type_map[i]; ok {
		return str
	}
	return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
}
