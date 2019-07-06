package functions

import d "github.com/joergreinhardt/gatwd/data"

type (
	//// COLLECTION
	ListCol func(...Expression) (Expression, ListCol)
	VecCol  func(...Expression) []Expression
	SetCol  func(...Paired) d.Mapped

	PairVal   func(...Expression) (Expression, Expression)
	KeyPair   func(...Expression) (Expression, string)
	TypedPair func(...Expression) (Expression, Typed)
	IndexPair func(...Expression) (Expression, int)

	PairList func(...Paired) (Paired, PairList)
	PairVec  func(...Paired) []Paired
)

//// RECURSIVE LIST OF VALUES
///
// base implementation of recursively linked lists
func ConList(list ListCol, elems ...Expression) ListCol {
	return list.Con(elems...)
}

func ConcatLists(a, b ListCol) ListCol {
	return ListCol(func(args ...Expression) (Expression, ListCol) {
		if len(args) > 0 {
			b = b.Con(args...)
		}
		var head Expression
		if head, a = a(); head != nil {
			return head, ConcatLists(a, b)
		}
		return b()
	})
}

func NewList(elems ...Expression) ListCol {
	return func(args ...Expression) (Expression, ListCol) {
		if len(args) > 0 {
			elems = append(elems, args...)
		}
		if len(elems) > 0 {
			var head = elems[0]
			if len(elems) > 1 {
				return head, NewList(
					elems[1:]...,
				)
			}
			return head, NewList()
		}
		return nil, NewList()
	}
}

func (l ListCol) Con(elems ...Expression) ListCol {
	return ListCol(func(args ...Expression) (Expression, ListCol) {
		return l(append(elems, args...)...)
	})
}

func (l ListCol) Push(elems ...Expression) ListCol {
	return ConcatLists(NewList(elems...), l)
}

func (l ListCol) Slice() []Expression {
	var vec = NewVector()
	var head Expression
	var tail Consumeable
	head, tail = l.Head(), l.Tail()
	for head != nil {
		vec = vec.Append(head)
		head, tail = tail.Consume()
	}
	return vec.Slice()
}
func (l ListCol) Eval() d.Native { return native(l.Slice()...) }
func (l ListCol) Call(args ...Expression) Expression {
	if len(args) > 0 {
		return l.Con(args...)
	}
	return l.Head()
}

// get n'st element in list
func (l ListCol) GetIdx(n int) Expression {
	var head, list = l()
	for i := 0; i < n; i++ {
		head, list = list()
		if head == nil {
			return NewNone()
		}
	}
	return head
}

// eval applys current heads eval method to passed arguments, or calle it empty

func (l ListCol) Empty() bool {
	if l.Head() != nil {
		if !None.Flag().Match(l.Head().TypeFnc()) ||
			!d.Nil.Flag().Match(l.Head().TypeNat()) {
			return false
		}
	}

	return true
}

// to determine the length of a recursive function, it has to be fully unwound,
// so use with care! (and ask yourself, what went wrong to make the length of a
// list be of importance)
func (l ListCol) Len() int {
	var length int
	var head, tail = l()
	if head != nil {
		length += 1 + tail.Len()
	}
	return length
}

func (l ListCol) Ident() Expression                  { return l }
func (l ListCol) Null() ListCol                      { return NewList() }
func (l ListCol) Tail() Consumeable                  { _, t := l(); return t }
func (l ListCol) Head() Expression                   { h, _ := l(); return h }
func (l ListCol) Consume() (Expression, Consumeable) { return l() }
func (l ListCol) TypeFnc() TyFnc                     { return List }
func (l ListCol) TypeNat() d.TyNat                   { return d.Function }
func (l ListCol) TailList() ListCol                  { _, t := l(); return t }
func (l ListCol) ConsumeList() (Expression, ListCol) {
	return l.Head(), l.TailList()
}

func (l ListCol) TypeElem() Typed {
	if l.Len() > 0 {
		return l.Head().Type()
	}
	return None.Type()
}
func (l ListCol) TypeName() string {
	if !l.TypeElem().TypeFnc().Match(None) {
		return "[" + l.TypeElem().Type().Name() + "]"
	}
	return "[]"
}
func (l ListCol) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (l ListCol) Type() TyDef {
	return Define(l.TypeName(), l.TypeElem().Type())
}

///////////////////////////////////////////////////////////////////////////////
//// PAIRS OF VALUES
///
// pairs can be created empty, key & value may be constructed later
func NewEmptyPair() PairVal {
	return func(args ...Expression) (a, b Expression) {
		if len(args) > 0 {
			if len(args) > 1 {
				return args[0], args[1]
			}
			return args[0], NewNone()
		}
		return NewNone(), NewNone()
	}
}

// new pair from two callable instances
func NewPair(l, r Expression) PairVal {
	return func(args ...Expression) (Expression, Expression) {
		if len(args) > 0 {
			if len(args) > 1 {
				return args[0], args[1]
			}
			return args[0], r
		}
		return l, r
	}
}

// pairs identity is a pair
func (p PairVal) Ident() Expression { return p }

// pair implements associative collection
func (p PairVal) Pair() Paired { return p }

// implement swappable
func (p PairVal) Swap() (Expression, Expression) { l, r := p(); return r, l }
func (p PairVal) SwappedPair() Paired            { return NewPair(p.Right(), p.Left()) }

// implement associated
func (p PairVal) Left() Expression               { l, _ := p(); return l }
func (p PairVal) Right() Expression              { _, r := p(); return r }
func (p PairVal) Both() (Expression, Expression) { return p() }

// implement sliced
func (p PairVal) Slice() []Expression { return []Expression{p.Left(), p.Right()} }

// associative implementing element access
func (p PairVal) Key() Expression   { return p.Left() }
func (p PairVal) Value() Expression { return p.Right() }

func (p PairVal) TypeFnc() TyFnc       { return Pair }
func (p PairVal) TypeNat() d.TyNat     { return d.Function }
func (p PairVal) KeyType() TyDef       { return p.Left().Type() }
func (p PairVal) ValType() TyDef       { return p.Right().Type() }
func (p PairVal) KeyNatType() d.TyNat  { return p.Left().TypeNat() }
func (p PairVal) ValNatType() d.TyNat  { return p.Right().TypeNat() }
func (p PairVal) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (p PairVal) TypeName() string {
	return "(" + p.Key().Type().Name() + ", " + p.Value().Type().Name() + ")"
}
func (p PairVal) Type() TyDef {
	return Define(p.TypeName(), NewPair(
		p.KeyType(), p.ValType()))
}

// implements compose
func (p PairVal) Empty() bool {
	if (p.Left() == nil ||
		(!p.Left().TypeFnc().Flag().Match(None) ||
			!p.Left().TypeNat().Flag().Match(d.Nil))) ||
		(p.Right() == nil ||
			(!p.Right().TypeFnc().Flag().Match(None) ||
				!p.Right().TypeNat().Flag().Match(d.Nil))) {
		return true
	}
	return false
}

// call calls the value, arguments are forwarded when calling right element
func (p PairVal) Eval() d.Native { return native(p) }
func (p PairVal) Call(args ...Expression) Expression {
	return NewPair(p.Left().Call(args...), p.Right().Call(args...))
}

///////////////////////////////////////////////////////////////////////////////
//// ASSOCIATIVE PAIRS
///
// pair composed of a string key and a functional value
func NewKeyPair(key string, val Expression) KeyPair {
	return func(...Expression) (Expression, string) { return val, key }
}

func (a KeyPair) KeyStr() string                     { _, key := a(); return key }
func (a KeyPair) Value() Expression                  { val, _ := a(); return val }
func (a KeyPair) Ident() Expression                  { return a }
func (a KeyPair) Left() Expression                   { return a.Value() }
func (a KeyPair) Right() Expression                  { return NewData(d.StrVal(a.KeyStr())) }
func (a KeyPair) Both() (Expression, Expression)     { return a.Left(), a.Right() }
func (a KeyPair) Pair() Paired                       { return NewPair(a.Both()) }
func (a KeyPair) Pairs() []Paired                    { return []Paired{NewPair(a.Both())} }
func (a KeyPair) Key() Expression                    { return a.Right() }
func (a KeyPair) Call(args ...Expression) Expression { return a.Value().Call(args...) }
func (a KeyPair) Eval() d.Native                     { return native(a) }
func (a KeyPair) KeyNatType() d.TyNat                { return d.String }
func (a KeyPair) ValNatType() d.TyNat                { return a.Value().TypeNat() }
func (a KeyPair) ValType() TyDef                     { return a.Value().Type() }
func (a KeyPair) KeyType() TyDef                     { return Key.Type() }
func (a KeyPair) TypeFnc() TyFnc                     { return Key }
func (a KeyPair) TypeNat() d.TyNat                   { return d.Function }
func (a KeyPair) FlagType() d.Uint8Val               { return Flag_Functional.U() }
func (p KeyPair) TypeName() string {
	return "(String, " + p.Value().Type().Name() + ")"
}
func (p KeyPair) Type() TyDef {
	return Define(p.TypeName(), NewPair(
		p.KeyType(), p.ValType()))
}

// implement swappable
func (p KeyPair) Swap() (Expression, Expression) {
	l, r := p()
	return NewData(d.StrVal(r)), l
}
func (p KeyPair) SwappedPair() Paired { return NewPair(p.Right(), p.Left()) }

func (a KeyPair) Empty() bool {
	if a.Key() != nil && a.Value() != nil && a.Value().TypeFnc() != None {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////
//// INDEX PAIR
///
// pair composed of an integer and a functional value
func NewIndexPair(idx int, val Expression) IndexPair {
	return func(...Expression) (Expression, int) { return val, idx }
}
func (a IndexPair) Ident() Expression                  { return a }
func (a IndexPair) Index() int                         { _, idx := a(); return idx }
func (a IndexPair) Value() Expression                  { val, _ := a(); return val }
func (a IndexPair) Left() Expression                   { return a.Value() }
func (a IndexPair) Right() Expression                  { return NewData(d.IntVal(a.Index())) }
func (a IndexPair) Both() (Expression, Expression)     { return a.Left(), a.Right() }
func (a IndexPair) Pair() Paired                       { return a }
func (a IndexPair) Pairs() []Paired                    { return []Paired{NewPair(a.Both())} }
func (a IndexPair) Key() Expression                    { return a.Right() }
func (a IndexPair) Call(args ...Expression) Expression { return a.Value().Call(args...) }
func (a IndexPair) Eval() d.Native                     { return native(a) }
func (a IndexPair) ValNatType() d.TyNat                { return a.Value().TypeNat() }
func (a IndexPair) KeyNatType() d.TyNat                { return d.Int }
func (a IndexPair) TypeFnc() TyFnc                     { return Index }
func (a IndexPair) TypeNat() d.TyNat                   { return d.Function }
func (a IndexPair) KeyType() TyDef                     { return Index.Type() }
func (a IndexPair) ValType() TyDef                     { return a.Value().Type() }
func (a IndexPair) FlagType() d.Uint8Val               { return Flag_Functional.U() }
func (a IndexPair) TypeName() string {
	return "(Index, " + a.Value().Type().Name() + ")"
}
func (a IndexPair) Type() TyDef {
	return Define(a.TypeName(), NewPair(a.KeyType(), a.ValType()))
}

// implement swappable
func (p IndexPair) Swap() (Expression, Expression) {
	l, r := p()
	return NewData(d.StrVal(r)), l
}
func (p IndexPair) SwappedPair() Paired { return NewPair(p.Right(), p.Left()) }
func (a IndexPair) Empty() bool {
	if a.Index() >= 0 && a.Value() != nil && a.Value().TypeFnc() != None {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////
//// LIST OF PAIRS
func ConPairList(list PairList, pairs ...Paired) PairList {
	return list.Con(pairs...)
}
func ConcatPairLists(a, b PairList) PairList {
	return PairList(func(args ...Paired) (Paired, PairList) {
		if len(args) > 0 {
			b = b.Con(args...)
		}
		var pair Paired
		if pair, a = a(); pair != nil {
			return pair, ConcatPairLists(a, b)
		}
		return b()
	})
}
func NewPairList(elems ...Paired) PairList {
	return func(pairs ...Paired) (Paired, PairList) {
		if len(pairs) > 0 {
			elems = append(elems, pairs...)
		}
		if len(elems) > 0 {
			var pair = elems[0]
			if len(elems) > 1 {
				return pair, NewPairList(
					elems[1:]...,
				)
			}
			return pair, NewPairList()
		}
		return nil, NewPairList()
	}
}
func (l PairList) Con(elems ...Paired) PairList {
	return PairList(func(args ...Paired) (Paired, PairList) {
		return l(append(elems, args...)...)
	})
}
func (l PairList) Push(elems ...Paired) PairList {
	return ConcatPairLists(NewPairList(elems...), l)
}
func (l PairList) Call(args ...Expression) Expression {
	var pairs = []Paired{}
	if len(args) > 0 {
		pairs = append(pairs, argsToPaired(args...)...)
	}
	var head Expression
	head, l = l(pairs...)
	return head
}

// eval applys current heads eval method to passed arguments, or calle it empty
func (l PairList) Eval() d.Native { return native(l) }

func (l PairList) Empty() bool {
	if pair := l.HeadPair(); pair != nil {
		return pair.Empty()
	}
	return true
}

// to determine the length of a recursive function, it has to be fully unwound,
// so use with care! (and ask yourself, what went wrong to make the length of a
// list be of importance)
func (l PairList) Len() int {
	var length int
	var head, tail = l()
	if head != nil {
		length += 1 + tail.Len()
	}
	return length
}

func (l PairList) Ident() Expression                        { return l }
func (l PairList) TypeFnc() TyFnc                           { return List }
func (l PairList) TypeNat() d.TyNat                         { return d.Function }
func (l PairList) FlagType() d.Uint8Val                     { return Flag_Functional.U() }
func (l PairList) Null() PairList                           { return NewPairList() }
func (l PairList) Consume() (Expression, Consumeable)       { return l() }
func (l PairList) ConsumePair() (Paired, ConsumeablePaired) { return l() }
func (l PairList) ConsumePairList() (Paired, PairList)      { return l() }
func (l PairList) Tail() Consumeable                        { _, t := l(); return t }
func (l PairList) TailPairs() ConsumeablePaired             { _, t := l(); return t }
func (l PairList) TailPairList() PairList                   { _, t := l(); return t }
func (l PairList) Head() Expression                         { h, _ := l(); return h }
func (l PairList) HeadPair() Paired                         { p, _ := l(); return p }
func (l PairList) TypeElem() Typed {
	if l.Len() > 0 {
		return l.Head().Type()
	}
	return None.Type()
}
func (l PairList) KeyNatType() d.TyNat {
	return l.Head().(Paired).KeyNatType()
}
func (l PairList) ValNatType() d.TyNat {
	return l.Head().(Paired).ValNatType()
}
func (l PairList) KeyType() TyDef {
	return l.Head().(PairVal).KeyType()
}
func (l PairList) ValType() TyDef {
	return l.Head().(Paired).ValType()
}
func (l PairList) TypeName() string {
	if l.Len() > 0 {
		return "[" + l.HeadPair().Type().Name() + "]"
	}
	return "[]"
}
func (l PairList) Type() TyDef {
	return Define(l.TypeName(), NewPair(
		l.KeyType(), l.ValType()))
}

// helper function to group arguments pair wise. assumes the arguments to
// either implement paired, or be alternating pairs of key & value. in case the
// number of passed arguments that are not pairs is uneven, last field will be
// filled up with a value of type none
func argsToPaired(args ...Expression) []Paired {
	var pairs = []Paired{}
	var alen = len(args)
	for i, arg := range args {
		if arg.TypeFnc().Match(Pair) {
			pairs = append(pairs, arg.(Paired))
		}
		if i < alen-2 {
			i = i + 1
			pairs = append(pairs, NewPair(arg, args[i]))
		}
		pairs = append(pairs, NewPair(arg, NewNone()))
	}
	return pairs
}

//////////////////////////////////////////////////////////////////////////////////////////
//// VECTORS (SLICES) OF VALUES
func NewEmptyVector(init ...Expression) VecCol { return NewVector() }

func NewVector(init ...Expression) VecCol {
	var vector = init
	return func(args ...Expression) []Expression {
		if len(args) > 0 {
			vector = append(
				vector,
				args...,
			)
		}
		return vector
	}
}

func ConVector(vec Vectorized, args ...Expression) VecCol {
	return NewVector(append(vec.Slice(), args...)...)
}

func AppendVectors(vec Vectorized, args ...Expression) VecCol {
	return NewVector(append(vec.Slice(), args...)...)
}

func AppendArgToVector(init ...Expression) VecCol {
	return func(args ...Expression) []Expression {
		return append(init, args...)
	}
}

func (v VecCol) Append(args ...Expression) VecCol {
	return NewVector(append(v(), args...)...)
}

func (v VecCol) Con(args ...Expression) VecCol {
	return ConVector(v, args...)
}

func (v VecCol) Ident() Expression { return v }

func (v VecCol) Call(d ...Expression) Expression {
	return NewVector(v(d...)...)
}

func (v VecCol) Eval() d.Native { return native(v()...) }

func (v VecCol) Head() Expression {
	if v.Len() > 0 {
		return v.Slice()[0]
	}
	return nil
}

func (v VecCol) Tail() Consumeable {
	if v.Len() > 1 {
		return NewVector(v.Slice()[1:]...)
	}
	return NewEmptyVector()
}

func (v VecCol) Consume() (Expression, Consumeable) {
	return v.Head(), v.Tail()
}

func (v VecCol) TailVec() VecCol {
	if v.Len() > 1 {
		return NewVector(v.Slice()[1:]...)
	}
	return NewEmptyVector()
}

func (v VecCol) ConsumeVec() (Expression, VecCol) {
	return v.Head(), v.TailVec()
}

func (v VecCol) Empty() bool {
	if len(v()) > 0 {
		for _, val := range v() {
			if !val.TypeNat().Flag().Match(d.Nil) &&
				!val.TypeFnc().Flag().Match(None) {
				return false
			}
		}
	}
	return true
}

func (v VecCol) Len() int            { return len(v()) }
func (v VecCol) Vector() VecCol      { return v }
func (v VecCol) Slice() []Expression { return v() }

func (v VecCol) Get(i int) (Expression, bool) {
	if i < v.Len() {
		return v()[i], true
	}
	return NewNone(), false
}

func (v VecCol) Set(i int, val Expression) (Vectorized, bool) {
	if i < v.Len() {
		var slice = v()
		slice[i] = val
		return VecCol(
			func(elems ...Expression) []Expression {
				return slice
			}), true

	}
	return v, false
}

func (v VecCol) Sort(flag d.TyNat) {
	var ps = Sort(v()...)
	ps.Sort(flag)
	v = NewVector(ps...)
}

func (v VecCol) Search(praed Expression) int {
	return Sort(v()...).Search(praed)
}

func (v VecCol) TypeFnc() TyFnc       { return Vector }
func (v VecCol) TypeNat() d.TyNat     { return d.Function }
func (v VecCol) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (v VecCol) TypeElem() Typed {
	if v.Len() > 0 {
		return v.Head().Type()
	}
	return None.Type()
}
func (v VecCol) Type() TyDef {
	return Define(v.TypeName(), v.TypeFnc())
}
func (v VecCol) TypeName() string {
	if v.Len() > 0 {
		return "[" + v.TypeElem().Type().Name() + "]"
	}
	return "[]"
}

///////////////////////////////////////////////////////////////////////////////
//// ASSOCIATIVE SLICE OF VALUE PAIRS
///
// list of associative pairs in sequential order associated, sorted and
// searched by left value of the pairs
func NewEmptyPairVec() PairVec {
	return PairVec(func(args ...Paired) []Paired {
		var pairs = []Paired{}
		if len(args) > 0 {
			pairs = append(pairs, args...)
		}
		return pairs
	})
}

func NewPairVectorFromPairs(pairs ...Paired) PairVec {
	return PairVec(func(args ...Paired) []Paired {
		if len(args) > 0 {
			return append(pairs, args...)
		}
		return pairs
	})
}

func ConPairListFromArgs(rec PairVec, args ...Expression) PairVec {
	var pairs = []Paired{}
	for _, arg := range args {
		if pair, ok := arg.(Paired); ok {
			pairs = append(pairs, pair)
		}
	}
	return NewPairVectorFromPairs(append(rec(), pairs...)...)
}

func NewPairVec(args ...Paired) PairVec {
	return NewPairVectorFromPairs(args...)
}

func ConPairVec(rec PairVec, pairs ...Paired) PairVec {
	return NewPairVectorFromPairs(append(rec(), pairs...)...)
}

func ConPairVecFromArgs(pvec PairVec, args ...Expression) PairVec {
	var pairs = pvec.Pairs()
	for _, arg := range args {
		if pair, ok := arg.(Paired); ok {
			pairs = append(pairs, pair)
		}
	}
	return PairVec(func(args ...Paired) []Paired {
		if len(args) > 0 {
			return ConPairVec(pvec, args...)()
		}
		return append(pvec(), pairs...)
	})
}
func (v PairVec) Con(args ...Expression) PairVec {
	return ConPairVecFromArgs(v, args...)
}

func (v PairVec) Consume() (Expression, Consumeable) {
	return v.Head(), v.Tail()
}

func (v PairVec) ConsumePairVec() (Paired, PairVec) {
	return v.HeadPair(), v.Tail().(PairVec)
}

func (v PairVec) Empty() bool {
	if len(v()) > 0 {
		for _, pair := range v() {
			if !pair.Empty() {
				return false
			}
		}
	}
	return true
}

func (v PairVec) TypeFnc() TyFnc   { return Vector }
func (v PairVec) TypeNat() d.TyNat { return d.Function }
func (v PairVec) TypeElem() Typed {
	if v.Len() > 0 {
		v()[0].Type()
	}
	return None.Type()
}
func (v PairVec) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (v PairVec) KeyNatType() d.TyNat {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().TypeNat()
	}
	return d.Nil
}
func (v PairVec) ValNatType() d.TyNat {
	if v.Len() > 0 {
		return v.Pairs()[0].Right().TypeNat()
	}
	return d.Nil
}
func (v PairVec) KeyType() TyDef {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().Type()
	}
	return None.Type()
}
func (v PairVec) ValType() TyDef {
	if v.Len() > 0 {
		return v.Pairs()[0].Right().Type()
	}
	return None.Type()
}
func (v PairVec) TypeName() string {
	if v.Len() > 0 {
		return "[" + v.Type().Type().Name() + "]"
	}
	return "[]"
}
func (v PairVec) Type() TyDef {
	return Define(v.TypeName(), NewPair(
		v.KeyType(), v.ValType()))
}

func (v PairVec) Len() int { return len(v()) }

func (v PairVec) Sort(flag d.TyNat) {
	var ps = SortPairs(v.Pairs()...)
	ps.Sort(flag)
	v = NewPairVectorFromPairs(ps...)
}

func (v PairVec) Search(praed Expression) int {
	return SortPairs(v.Pairs()...).Search(praed)
}

func (v PairVec) Get(idx int) (Paired, bool) {
	if idx < v.Len()-1 {
		return v()[idx], true
	}
	return NewKeyPair("None", NewNone()), false
}

func (v PairVec) GetVal(praed Expression) (Expression, bool) {
	return NewPairVectorFromPairs(SortPairs(v.Pairs()...).Get(praed)), true
}

func (v PairVec) Range(praed Expression) []Paired {
	return SortPairs(v.Pairs()...).Range(praed)
}

func (v PairVec) Pairs() []Paired {
	var pairs = []Paired{}
	for _, pair := range v() {
		pairs = append(pairs, pair)
	}
	return pairs
}

func (v PairVec) ConsumePair() (Paired, ConsumeablePaired) {
	var pairs = v()
	if len(pairs) > 0 {
		if len(pairs) > 1 {
			return pairs[0], NewPairVec(pairs[1:]...)
		}
		return pairs[0], NewPairVec()
	}
	return nil, NewPairVec()
}

func (v PairVec) SwitchedPairs() []Paired {
	var switched = []Paired{}
	for _, pair := range v() {
		switched = append(
			switched,
			pair,
		)
	}
	return switched
}

func (v PairVec) SetVal(key, value Expression) (AssociativeCollected, bool) {
	if idx := v.Search(key); idx >= 0 {
		var pairs = v()
		pairs[idx] = NewKeyPair(key.String(), value)
		return NewPairVec(pairs...), true
	}
	return NewPairVec(append(v(), NewKeyPair(key.String(), value))...), false
}

func (v PairVec) Slice() []Expression {
	var fncs = []Expression{}
	for _, pair := range v() {
		fncs = append(fncs, NewPair(pair.Left(), pair.Right()))
	}
	return fncs
}

func (v PairVec) HeadPair() Paired {
	if v.Len() > 0 {
		return v()[0].(Paired)
	}
	return NewPair(NewNone(), NewNone())
}
func (v PairVec) Head() Expression {
	if v.Len() > 0 {
		return v.Pairs()[0]
	}
	return nil
}

func (v PairVec) TailPairs() ConsumeablePaired {
	if v.Len() > 1 {
		return NewPairVectorFromPairs(v.Pairs()[1:]...)
	}
	return NewEmptyPairVec()
}
func (v PairVec) Tail() Consumeable {
	if v.Len() > 1 {
		return NewPairVectorFromPairs(v.Pairs()[1:]...)
	}
	return NewEmptyPairVec()
}

func (v PairVec) Call(args ...Expression) Expression {
	return v.Con(args...)
}

func (v PairVec) Eval() d.Native { return native(v) }

///////////////////////////////////////////////////////////////////////////////
//// ASSOCIATIVE SET (HASH MAP OF VALUES)
///
// unordered associative set of key/value pairs that can be sorted, accessed
// and searched by the left (key) value of the pair
func ConSet(set SetCol, pairs ...Paired) SetCol {
	var knat = set.KeyNatType()
	var vnat = set.ValNatType()
	var m = set()
	for _, arg := range pairs {
		if pair, ok := arg.(Paired); ok {
			if pair.Left().TypeNat() == knat &&
				pair.Right().TypeNat() == vnat {
				m.Set(pair.Left(), pair.Right())
			}
		}
	}
	return SetCol(func(pairs ...Paired) d.Mapped { return m })
}

// new set discriminates between sets where all members have identical keys and
// such with mixed keys and chooses the appropriate native set accordingly.
func NewSet(pairs ...Paired) SetCol {
	var set d.Mapped
	var knat d.BitFlag
	if len(pairs) > 0 {
		// first passed pair determines initial key type
		knat = pairs[0].Left().TypeNat().Flag()
		// OR concat all the keys types, to see if arguments are of
		// mixed type
		for _, pair := range pairs {
			knat = knat | pair.Left().TypeNat().Flag()
		}
		// for sets with pure key type, choose the appropriate native
		// set type
		if knat.Count() == 1 {
			switch {
			case knat.Match(d.Int):
				set = d.SetInt{}
			case knat.Match(d.Uint):
				set = d.SetUint{}
			case knat.Match(d.Type):
				set = d.SetFlag{}
			case knat.Match(d.Float):
				set = d.SetFloat{}
			case knat.Match(d.String):
				set = d.SetString{}
			}
		} else {
			// otherwise choose a set keyed by interface type to
			// keep every possible kind of value
			set = d.SetVal{}
		}
	}
	return SetCol(func(pairs ...Paired) d.Mapped { return set })
}

// splits set into two lists, one containing all keys and the other all values
func (v SetCol) Split() (VecCol, VecCol) {
	var keys, vals = []Expression{}, []Expression{}
	for _, pair := range v.Pairs() {
		keys = append(keys, pair.Left())
		vals = append(vals, pair.Right())
	}
	return NewVector(keys...), NewVector(vals...)
}

func (v SetCol) Pairs() []Paired {
	var pairs = []Paired{}
	for _, field := range v().Fields() {
		pairs = append(
			pairs,
			NewPair(
				NewData(field.Left()),
				NewData(field.Right())))
	}
	return pairs
}

// return all members keys
func (v SetCol) Keys() VecCol { k, _ := v.Split(); return k }

// return all members values
func (v SetCol) Data() VecCol { _, d := v.Split(); return d }

func (v SetCol) Len() int { return v().Len() }

func (v SetCol) Empty() bool {
	for _, pair := range v.Pairs() {
		if !pair.Empty() {
			return false
		}
	}
	return true
}

func (v SetCol) GetVal(key Expression) (Expression, bool) {
	var m = v()
	if value, ok := m.Get(key); ok {
		return NewData(value), ok
	}
	return NewNone(), false
}

func (v SetCol) SetVal(key, value Expression) (AssociativeCollected, bool) {
	var m = v()
	return SetCol(func(pairs ...Paired) d.Mapped { return m.Set(key, value) }), true
}

func (v SetCol) Slice() []Expression {
	var pairs = []Expression{}
	for _, pair := range v.Pairs() {
		pairs = append(pairs, pair)
	}
	return pairs
}

// call method performs a value lookup
func (v SetCol) Call(args ...Expression) Expression {
	var results = []Expression{}
	for _, arg := range args {
		if val, ok := v.GetVal(arg); ok {
			results = append(results, val)
		}
	}
	if len(results) > 0 {
		if len(results) > 1 {
			return NewVector(results...)
		}
		return results[0]
	}
	return NewNone()
}

// eval method performs a value lookup and returns contained value as native
// without any conversion
func (v SetCol) Eval() d.Native { return native(v) }

func (v SetCol) TypeNat() d.TyNat { return d.Function }
func (v SetCol) TypeFnc() TyFnc   { return Set }
func (v SetCol) TypeElem() Typed {
	if v.Len() > 0 {
		return v.Head().Type()
	}
	return None.Type()
}
func (v SetCol) TypeName() string {
	if v.Len() > 0 {
		return "{" + v.Pairs()[0].Left().Type().Name() +
			":: " + v.Pairs()[0].Right().Type().Name() + "}"
	}
	return "{}"
}
func (v SetCol) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (v SetCol) Type() TyDef {
	return Define(v.TypeName(), NewPair(
		v.KeyType(), v.ValType()))
}

func (v SetCol) KeyNatType() d.TyNat {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().TypeNat()
	}
	return d.Nil
}

func (v SetCol) KeyType() TyDef {
	if v.Len() > 0 {
		return v.Pairs()[0].KeyType()
	}
	return None.Type()
}
func (s SetCol) ValType() TyDef {
	if s.Len() > 0 {
		return s.Pairs()[0].ValType()
	}
	return None.Type()
}

func (s SetCol) ValNatType() d.TyNat {
	if s.Len() > 0 {
		return s.Pairs()[0].Right().TypeNat()
	}
	return d.Nil
}

func (v SetCol) Head() Expression {
	if v.Len() > 0 {
		var vec = NewPairVectorFromPairs(
			v.Pairs()...,
		)
		vec.Sort(v.KeyNatType())
		return vec()[0]
	}
	return nil
}

func (v SetCol) Tail() Consumeable {
	if v.Len() > 1 {
		var vec = NewPairVectorFromPairs(
			v.Pairs()...,
		)
		vec.Sort(v.KeyNatType())
		return NewPairVec(vec()[:1]...)
	}
	return nil
}

func (v SetCol) Consume() (Expression, Consumeable) {
	return v.Head(), v.Tail()
}

func (v SetCol) TailPairVec() PairVec {
	if v.Len() > 1 {
		var vec = NewPairVectorFromPairs(
			v.Pairs()...,
		)
		vec.Sort(v.KeyNatType())
		return NewPairVec(vec()[:1]...)
	}
	return nil
}

func (v SetCol) ConsumeSet() (Expression, PairVec) {
	return v.Head(), v.TailPairVec()
}
