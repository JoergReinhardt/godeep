// Code generated by "stringer -type=TyFnc"; DO NOT EDIT.

package functions

import "strconv"

const _TyFnc_name = "TypeDataDefinitionExpressionVariableFunctionClosureArgumentParameterAccessorAttributPredicateAggregatorGeneratorConstructorFunctorMonadConditionFalseTrueJustNoneCaseSwitchEitherOrIfElseErrorPairListTupleUniSetMuliSetAssocVecRecordVectorDLinkLinkNodeTreeIOHigherOrder"

var _TyFnc_map = map[TyFnc]string{
	1:             _TyFnc_name[0:4],
	2:             _TyFnc_name[4:8],
	4:             _TyFnc_name[8:18],
	8:             _TyFnc_name[18:28],
	16:            _TyFnc_name[28:36],
	32:            _TyFnc_name[36:44],
	64:            _TyFnc_name[44:51],
	128:           _TyFnc_name[51:59],
	256:           _TyFnc_name[59:68],
	512:           _TyFnc_name[68:76],
	1024:          _TyFnc_name[76:84],
	2048:          _TyFnc_name[84:93],
	4096:          _TyFnc_name[93:103],
	8192:          _TyFnc_name[103:112],
	16384:         _TyFnc_name[112:123],
	32768:         _TyFnc_name[123:130],
	65536:         _TyFnc_name[130:135],
	131072:        _TyFnc_name[135:144],
	262144:        _TyFnc_name[144:149],
	524288:        _TyFnc_name[149:153],
	1048576:       _TyFnc_name[153:157],
	2097152:       _TyFnc_name[157:161],
	4194304:       _TyFnc_name[161:165],
	8388608:       _TyFnc_name[165:171],
	16777216:      _TyFnc_name[171:177],
	33554432:      _TyFnc_name[177:179],
	67108864:      _TyFnc_name[179:181],
	134217728:     _TyFnc_name[181:185],
	268435456:     _TyFnc_name[185:190],
	536870912:     _TyFnc_name[190:194],
	1073741824:    _TyFnc_name[194:198],
	2147483648:    _TyFnc_name[198:203],
	4294967296:    _TyFnc_name[203:209],
	8589934592:    _TyFnc_name[209:216],
	17179869184:   _TyFnc_name[216:224],
	34359738368:   _TyFnc_name[224:230],
	68719476736:   _TyFnc_name[230:236],
	137438953472:  _TyFnc_name[236:241],
	274877906944:  _TyFnc_name[241:245],
	549755813888:  _TyFnc_name[245:249],
	1099511627776: _TyFnc_name[249:253],
	2199023255552: _TyFnc_name[253:255],
	4398046511104: _TyFnc_name[255:266],
}

func (i TyFnc) String() string {
	if str, ok := _TyFnc_map[i]; ok {
		return str
	}
	return "TyFnc(" + strconv.FormatInt(int64(i), 10) + ")"
}
