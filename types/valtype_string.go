// Code generated by "stringer -type=ValType"; DO NOT EDIT.

package types

import "strconv"

const _ValType_name = "NilMetaTypeBoolIntInt8Int16Int32UintUint16Uint32FloatFlt32ImagImag64ByteBytesStringTimeDurationErrorAttrSliceListMuliListAttrListRecordListUniSetMuliSetAttrSetRecordChainedListLinkedListDoubleLinkTupleNodeTreeFunctionPttrNATIVES"

var _ValType_map = map[ValType]string{
	0:            _ValType_name[0:3],
	1:            _ValType_name[3:11],
	4:            _ValType_name[11:15],
	8:            _ValType_name[15:18],
	16:           _ValType_name[18:22],
	32:           _ValType_name[22:27],
	64:           _ValType_name[27:32],
	128:          _ValType_name[32:36],
	256:          _ValType_name[36:42],
	512:          _ValType_name[42:48],
	1024:         _ValType_name[48:53],
	2048:         _ValType_name[53:58],
	4096:         _ValType_name[58:62],
	8192:         _ValType_name[62:68],
	16384:        _ValType_name[68:72],
	32768:        _ValType_name[72:77],
	65536:        _ValType_name[77:83],
	131072:       _ValType_name[83:87],
	262144:       _ValType_name[87:95],
	524288:       _ValType_name[95:100],
	1048576:      _ValType_name[100:104],
	2097152:      _ValType_name[104:109],
	4194304:      _ValType_name[109:113],
	8388608:      _ValType_name[113:121],
	16777216:     _ValType_name[121:129],
	33554432:     _ValType_name[129:139],
	67108864:     _ValType_name[139:145],
	134217728:    _ValType_name[145:152],
	268435456:    _ValType_name[152:159],
	536870912:    _ValType_name[159:165],
	1073741824:   _ValType_name[165:176],
	2147483648:   _ValType_name[176:186],
	4294967296:   _ValType_name[186:196],
	8589934592:   _ValType_name[196:201],
	17179869184:  _ValType_name[201:205],
	34359738368:  _ValType_name[205:209],
	68719476736:  _ValType_name[209:217],
	137438953472: _ValType_name[217:221],
	274877906944: _ValType_name[221:228],
}

func (i ValType) String() string {
	if str, ok := _ValType_map[i]; ok {
		return str
	}
	return "ValType(" + strconv.FormatInt(int64(i), 10) + ")"
}
