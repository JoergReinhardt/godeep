package data

type Chain []Data

func ConChain(val ...Data) Chain {
	l := make([]Data, 0, len(val))
	l = append(l, val...)
	return l
}
func ChainContainedTypes(c []Data) BitFlag {
	var flag = BitFlag(0)
	for _, d := range c {
		if FlagMatch(d.Flag(), Slice.Flag()) {
			ChainContainedTypes(d.(Chain))
			continue
		}
		flag = flag | d.Flag()
	}
	return flag
}
func (c Chain) Flag() BitFlag           { return Slice.Flag() }
func (c Chain) ContainedTypes() BitFlag { return ChainContainedTypes(c.Slice()) }
func (c Chain) Eval() Data              { return c }
func (c Chain) Null() Chain             { return []Data{} }

// SLICE ->
func (v Chain) Slice() []Data { return v }
func (v Chain) Len() int      { return len(v) }

// COLLECTION
func (s Chain) Empty() bool            { return ChainEmpty(s) }
func (s Chain) Head() (h Data)         { return s[0] }
func (s Chain) Tail() (c Consumeable)  { return s[:1] }
func (s Chain) Shift() (c Consumeable) { return s[:1] }

func ChainClear(s Chain) {
	if len(s) > 0 {
		for _, v := range s {
			if d, ok := v.(Destructable); ok {
				d.Clear()
			}
		}
	}
	s = nil
}
func ElemEmpty(d Data) bool {
	// not flagged nil, not a composition either...
	if !FlagMatch(d.Flag(), (Nil.Flag() | Slice.Flag())) {
		if d != nil { // not a nil pointer...
			// --> not empty
			return false
		}
	}
	// since it's a composition, inspect...
	if FlagMatch(d.Flag(), Slice.Flag()) {
		// slice --> call sliceEmpty
		if sl, ok := d.(Chain); ok {
			return ChainEmpty(sl)
		}
		// other sort of collection...
		if col, ok := d.(Collected); ok {
			// --> call it's empty method
			return col.Empty()
		}
	}
	// no idea, what this is, so better call it empty
	return true
}
func ChainEmpty(s Chain) bool {
	if len(s) == 0 { // empty, as in no element...
		return true
	}
	if len(s) > 0 { // empty as in contains empty elements exclusively...
		for _, elem := range ChainSlice(s) { // return at first non empty
			if !ElemEmpty(elem) {
				return false
			}
		}
	} // --> all contained elements are empty
	return true
}

///// CONVERT TO SLICE OF NATIVES ////////
func ChainToNativeSlice(c Chain) NativeVec {
	f := ChainGet(c, 0).Flag()
	if ChainAll(c, func(i int, c Data) bool {
		return FlagMatch(f, c.Flag())
	}) {
		return ConNativeSlice(f, c.Slice()...)
	}
	return c
}
func (c Chain) NativeSlice() []interface{} {
	var s = make([]interface{}, 0, c.Len())
	for _, d := range c.Slice() {
		s = append(s, d.(Evaluable).Eval())
	}
	return s
}

//// LIST OPERATIONS ///////
func ChainFoldL(
	c Chain,
	fn func(i int, data Data, accu Data) Data,
	init Data,
) Data {
	var accu = init
	for i, d := range c.Slice() {
		accu = fn(i, d, accu)
	}
	return accu
}
func ChainMap(c Chain, fn func(i int, d Data) Data) Chain {
	var ch = make([]Data, 0, c.Len())
	for i, d := range c.Slice() {
		ch = append(ch, fn(i, d))
	}
	return ch
}
func ChainFilter(c Chain, fn func(i int, d Data) bool) Chain {
	var ch = []Data{}
	for i, d := range c.Slice() {
		if fn(i, d) {
			ch = append(ch, d)
		}
	}
	return ch
}
func ChainAny(c Chain, fn func(i int, d Data) bool) bool {
	var answ = false
	for i, d := range c.Slice() {
		if fn(i, d) {
			return true
		}
	}
	return answ
}
func ChainAll(c Chain, fn func(i int, d Data) bool) bool {
	var answ = true
	for i, d := range c.Slice() {
		if !fn(i, d) {
			return false
		}
	}
	return answ
}
func ChainReverse(c Chain) Chain {
	var ch = make([]Data, 0, c.Len())
	for i := c.Len() - 1; i > 0; i-- {
		ch = append(ch, ChainGet(c, i))
	}
	return ch
}

// ACCESSABLE SLICE
func ChainGet(s Chain, i int) Data { return s[i] }

// MUTABLE SLICE
func ChainSet(s Chain, i int, v Data) Chain { s[i] = v; return s }

// ITERATOR
func ChainNext(s Chain) (v Data, i Chain) {
	if len(s) > 0 {
		if len(s) > 1 {
			return s[0], s[1:]
		}
		return s[0], Chain([]Data{NilVal{}})
	}
	return NilVal{}, Chain([]Data{NilVal{}})
}

type Iter func() (Data, Iter)

func ConIter(c Chain) Iter {
	data, chain := ChainNext(c)
	return func() (Data, Iter) {
		return data, ConIter(chain)
	}
}

// BOOTOM & TOP
func ChainFirst(s Chain) Data {
	if s.Len() > 0 {
		return s[0]
	}
	return NilVal{}
}
func ChainLast(s Chain) Data {
	if s.Len() > 0 {
		return s[s.Len()-1]
	}
	return NilVal{}
}

// LIFO QUEUE
func ChainPut(s Chain, v Data) Chain {
	return append(s, v)
}
func ChainAppend(s Chain, v ...Data) Chain {
	return append(s, v...)
}
func ChainPull(s Chain) (Data, Chain) {
	if s.Len() > 0 {
		return s[s.Len()-1], s[:s.Len()-1]
	}
	return NilVal{}, s
}

// FIFO STACK
func ChainAdd(s Chain, v ...Data) Chain {
	return append(v, s...)
}
func ChainPush(s Chain, v Data) Chain {
	return append([]Data{v}, s...)
}
func ChainPop(s Chain) (Data, Chain) {
	if len(s) > 0 {
		return s[0], s[1:]
	}
	return NilVal{}, s
}

// TUPLE
func ChainHead(s Chain) (h Data)   { return s[0] }
func ChainTail(s Chain) (c []Data) { return s[:1] }
func ChainDecap(s Chain) (h Data, t Chain) {
	if !ChainEmpty(s) {
		return s[0], t[:1]
	}
	return NilVal{}, ConChain(NilVal{})
}

// SLICE
func ChainSlice(s Chain) []Data { return []Data(s) }
func ChainLen(s Chain) int      { return len(s) }
func ChainSplit(s Chain, i int) (Chain, Chain) {
	h, t := s[:i], s[i:]
	return h, t
}
func ChainCut(s Chain, i, j int) Chain {
	copy(s[i:], s[j:])
	// to prevent a possib. mem leak
	for k, n := len(s)-j+i, len(s); k < n; k++ {
		s[k] = nil
	}
	return s[:len(s)-j+i]
}
func ChainDelete(s Chain, i int) Chain {
	copy(s[i:], s[i+1:])
	s[len(s)-1] = nil
	return s[:len(s)-1]
}
func ChainInsert(s Chain, i int, v Data) Chain {
	s = append(s, NilVal{})
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}
func ChainInsertVector(s Chain, i int, v ...Data) Chain {
	return append(s[:i], append(v, s[i:]...)...)
}
func ChainAttrType(s Chain) BitFlag { return Int.Flag() }

///// TODO: perf test that‥. test sliding window, or similar sophistited shenaegans.
//func ChainAdd(s Chain, v ...Data) Chain {
//	if len(s) >= cap(s)+len(v)/2 {
//		return append(append(make([]Data, 0, len(v)+len(s)), v...), s...)
//	}
//	return append(v, s...)
//}
//func ChainPush(s Chain, v Data) Chain {
//	if len(s) >= cap(s)/2 {
//		return append(append(make([]Data, 0, (len(s))*2), v), s...)
//	}
//	return append([]Data{v}, s...)
//}
