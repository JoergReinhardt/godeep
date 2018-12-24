package types

// DESTRUCTABLE SLICE

func newSlice(val ...Evaluable) slice {
	l := make([]Evaluable, 0, len(val))
	l = append(l, val...)
	return l
}

func sliceClear(s slice) {
	if len(s) > 0 {
		for i, v := range s {
			if !fmatch(v.Type(), Nullable) {
				if d, ok := v.(Destructable); ok {
					d.Clear()
				}
			}
			s[i] = nil
		}
	}
	s = nil
}

// ACCESSABLE SLICE
func sliceGetInt(s slice, i int) Evaluable { return s[i] }

// MUTABLE SLICE
func sliceSetInt(s slice, i int, v Evaluable) slice { s[i] = v; return s }

// ITERATOR
func sliceNext(s slice) (v Evaluable, i slice) {
	if len(s) > 0 {
		if len(s) > 1 {
			return s[0], s[1:]
		}
		return s[0], slice([]Evaluable{nilVal{}})
	}
	return nilVal{}, slice([]Evaluable{nilVal{}})
}

// BOOTOM & TOP
func sliceFirst(s slice) Evaluable {
	if s.Len() > 0 {
		return s[0]
	}
	return s
}
func sliceLast(s slice) Evaluable {
	if s.Len() > 0 {
		return s[s.Len()-1]
	}
	return s
}

// LIFO QUEUE
func slicePut(s slice, v Evaluable) slice {
	if len(s) == cap(s) {
		return append(append(make([]Evaluable, 0, len(s)*2), s...), v)
	}
	return append(s, v)
}
func sliceAppend(s slice, v ...Evaluable) slice {
	if len(s) == cap(s) {
		return append(append(make([]Evaluable, 0, (len(s)+len(v))), s...), v...)
	}
	return append(s, v...)
}
func slicePull(s slice) (Evaluable, slice) {
	if s.Len() > 0 {
		return s[s.Len()-1], s[:s.Len()-1]
	}
	return nilVal{}, s
}

// FIFO STACK
func sliceAdd(s slice, v ...Evaluable) slice {
	if len(s) == cap(s)+len(v) {
		return append(append(make([]Evaluable, 0, len(v)+len(s)), v...), s...)
	}
	return append(v, s...)
}
func slicePush(s slice, v Evaluable) slice {
	if len(s) == cap(s) {
		return append(append(make([]Evaluable, 0, (len(s))*2), v), s...)
	}
	return append([]Evaluable{v}, s...)
}
func slicePop(s slice) (Evaluable, slice) {
	if s.Len() > 0 {
		return s[0], s[1:]
	}
	return nilVal{}, s
}

// ARITY

// TUPLE
func sliceHead(s slice) (h Evaluable) { return s[0] }
func sliceTail(s slice) (c Evaluable) { return s[:1] }
func sliceDecap(s slice) (h Evaluable, t Nested) {
	return h, t
}

// N-TUPLE
func sliceHeadNary(s slice, arity int) (h Evaluable) { return s[:arity] }
func sliceTailNary(s slice, arity int) (c Evaluable) { return s[arity:] }
func sliceDecapNary(s slice, arity int) (h Evaluable, t slice) {
	if s.Len()+1 > arity {
		return s[:arity], s[arity:]
	}
	return h, t
}

// SLICE
func sliceSlice(s slice) []Evaluable { return []Evaluable(s) }
func sliceLen(s slice) int           { return len(s) }
func sliceSplit(s slice, i int) (slice, slice) {
	h, t := s[:i], s[i:]
	return h, t
}
func sliceCut(s slice, i, j int) slice {
	copy(s[i:], s[j:])
	// to prevent a possib. mem leak
	for k, n := len(s)-j+i, len(s); k < n; k++ {
		s[k] = nil
	}
	return s[:len(s)-j+i]
}
func sliceDelete(s slice, i int) slice {
	copy(s[i:], s[i+1:])
	s[len(s)-1] = nil
	return s[:len(s)-1]
}
func sliceInsert(s slice, i int, v Evaluable) slice {
	s = append(s, nilVal{})
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}
func sliceInsertVari(s slice, i int, v ...Evaluable) slice {
	return append(s[:i], append(v, s[i:]...)...)
}
func sliceAttrType(s slice) flag { return Int.Type() }