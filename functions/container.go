/*
  FUNCTIONAL CONTAINERS

  containers implement enumeration of functional types, aka lists, vectors sets, pairs, tuples‥.
*/
package functions

import (
	d "github.com/joergreinhardt/gatwd/data"
)

type (
	//// DATA
	DataVal func(args ...d.Native) d.Native

	//// EXPRESSIONS
	ConstantExpr func() Callable
	UnaryExpr    func(Callable) Callable
	BinaryExpr   func(a, b Callable) Callable
	NaryExpr     func(...Callable) Callable

	//// COLLECTIONS
	PairVal        func(...Callable) (Callable, Callable)
	AssocPair      func(...Callable) (string, Callable)
	IndexPair      func(...Callable) (int, Callable)
	ListVal        func(...Callable) (Callable, ListVal)
	VecVal         func(...Callable) []Callable
	TupleVal       func(...Callable) []Callable
	AccociativeVal func(...PairVal) []PairVal
	SetVal         func(...PairVal) d.Mapped
)

// reverse arguments
func reverseArguments(args ...Callable) []Callable {

	var rev = []Callable{}

	for i := len(args) - 1; i > 0; i-- {

		rev = append(rev, args[i])
	}

	return rev
}

func functionalizeNatives(args ...d.Native) []Callable {

	var result = []Callable{}

	for _, arg := range args {

		result = append(result, NewFromData(arg))
	}

	return result
}

func evaluateFunctionals(args ...Callable) []d.Native {

	var result = []d.Native{}

	for _, arg := range args {

		result = append(result, arg.Eval())
	}

	return result
}

//// DATA
func New(inf ...interface{}) Callable { return NewFromData(d.New(inf...)) }

func NewFromData(data ...d.Native) DataVal {

	var eval func(...d.Native) d.Native

	for _, val := range data {
		eval = val.Eval
	}

	return func(args ...d.Native) d.Native { return eval(args...) }
}

func (n DataVal) Eval(args ...d.Native) d.Native { return n().Eval(args...) }

func (n DataVal) Call(vals ...Callable) Callable {

	var results = NewVector()

	for _, val := range vals {
		// evaluate arguments to yield contained natives
		results = ConVector(
			results,
			DataVal(func(arguments ...d.Native) d.Native {
				return val.Eval(arguments...)
			}),
		)
	}
	return results
}

func (n DataVal) TypeFnc() TyFnc      { return Data }
func (n DataVal) TypeNat() d.TyNative { return n().TypeNat() }
func (n DataVal) String() string      { return n().String() }

//// STATIC EXPRESSIONS
///
// CONSTANT EXPRESSION
func NewConstant(fnc Callable) Callable          { return fnc }
func (c ConstantExpr) Ident() Callable           { return c() }
func (c ConstantExpr) TypeFnc() TyFnc            { return Expression }
func (c ConstantExpr) TypeNat() d.TyNative       { return c().TypeNat() }
func (c ConstantExpr) Call(...Callable) Callable { return c() }
func (c ConstantExpr) Eval(...d.Native) d.Native { return c().Eval() }

/// UNARY EXPRESSION
func NewUnaryExpr(fnc func(Callable) Callable) UnaryExpr { return fnc }
func (u UnaryExpr) Ident() Callable                      { return u }
func (u UnaryExpr) TypeFnc() TyFnc                       { return Expression }
func (u UnaryExpr) TypeNat() d.TyNative                  { return d.Function.TypeNat() }
func (u UnaryExpr) Call(arg ...Callable) Callable        { return u(arg[0]) }
func (u UnaryExpr) Eval(arg ...d.Native) d.Native        { return u(NewFromData(arg...)) }

/// BINARY EXPRESSION
func NewBinaryExpr(fnc func(l, r Callable) Callable) BinaryExpr { return fnc }

func (b BinaryExpr) Ident() Callable                { return b }
func (b BinaryExpr) TypeFnc() TyFnc                 { return Expression }
func (b BinaryExpr) TypeNat() d.TyNative            { return d.Function.TypeNat() }
func (b BinaryExpr) Call(args ...Callable) Callable { return b(args[0], args[1]) }
func (b BinaryExpr) Eval(args ...d.Native) d.Native {
	return b(NewFromData(args[0]), NewFromData(args[1]))
}

/// NARY EXPRESSION
func NewNaryExpr(fnc func(...Callable) Callable) NaryExpr { return fnc }
func (n NaryExpr) Ident() Callable                        { return n }
func (n NaryExpr) TypeFnc() TyFnc                         { return Expression }
func (n NaryExpr) TypeNat() d.TyNative                    { return d.Function.TypeNat() }
func (n NaryExpr) Call(d ...Callable) Callable            { return n(d...) }
func (n NaryExpr) Eval(args ...d.Native) d.Native {
	var params = []Callable{}
	for _, arg := range args {
		params = append(params, NewFromData(arg))
	}
	return n(params...)
}

/// PAIRS OF VALUES
func NewEmptyPair() PairVal {
	return func(args ...Callable) (a, b Callable) {
		return NewNoOp(), NewNoOp()
	}
}

func NewPair(l, r Callable) PairVal {

	return func(args ...Callable) (Callable, Callable) {

		if len(args) > 0 {

			if len(args) > 1 {

				return args[0], args[1]
			}

			return args[0], nil
		}

		return l, r
	}
}

func NewPairFromData(l, r d.Native) PairVal {

	return func(args ...Callable) (Callable, Callable) {

		if len(args) > 0 {

			if len(args) > 1 {

				// return pointers to natives eval functions
				return DataVal(args[0].Eval), DataVal(args[1].Eval)
			}

			return DataVal(args[0].Eval), nil
		}

		return DataVal(l.Eval), DataVal(r.Eval)
	}
}

func NewPairFromLiteral(l, r interface{}) PairVal {

	return func(args ...Callable) (Callable, Callable) {

		if len(args) > 0 {

			if len(args) > 1 {

				// return values eval methods as continuations
				return DataVal(
						d.New(args[0]).Eval,
					),
					DataVal(
						d.New(args[1]).Eval,
					)
			}

			return DataVal(d.New(args[0]).Eval), nil
		}

		return DataVal(d.New(l).Eval), DataVal(d.New(r).Eval)
	}
}

func (p PairVal) Ident() Callable { return p }
func (p PairVal) Pair() Callable  { return p }

// construct value pairs from any consumeable assuming keys and values alter
func ConPair(list Consumeable) (PairVal, Consumeable) {

	var first, tail = list.DeCap()

	if first != nil {

		var second Callable
		second, tail = tail.DeCap()

		if tail != nil {
			// walk list generate a pair every second step
			// recursively.
			return NewPair(first, second), tail
		}
		// if number of elements in list is not dividable by two, last
		// element will contain an empty list as its right element
		return NewPair(first, tail), nil
	}
	// argument consumeable vanished, return nil for left and right
	return nil, nil
}

// implement consumeable
func (p PairVal) DeCap() (Callable, Consumeable) { l, r := p(); return l, NewList(r) }
func (p PairVal) Head() Callable                 { l, _ := p(); return l }
func (p PairVal) Tail() Consumeable              { _, r := p(); return NewPair(r, NewNoOp()) }

// implement swappable
func (p PairVal) Swap() (Callable, Callable) { l, r := p(); return r, l }
func (p PairVal) SwappedPair() PairVal       { return NewPair(p.Right(), p.Left()) }

// implement associated
func (p PairVal) Left() Callable             { l, _ := p(); return l }
func (p PairVal) Right() Callable            { _, r := p(); return r }
func (p PairVal) Both() (Callable, Callable) { return p() }

// implement sliced
func (p PairVal) Slice() []Callable { return []Callable{p.Left(), p.Right()} }

// associative implementing element access
func (p PairVal) Key() Callable   { return p.Left() }
func (p PairVal) Value() Callable { return p.Right() }

// key and values native and functional types
func (p PairVal) KeyType() TyFnc           { return p.Left().TypeFnc() }
func (p PairVal) KeyNatType() d.TyNative   { return p.Left().TypeNat() }
func (p PairVal) ValueType() TyFnc         { return p.Right().TypeFnc() }
func (p PairVal) ValueNatType() d.TyNative { return p.Right().TypeNat() }

// slightly different element types, since right value is a list now
func (p PairVal) HeadType() TyFnc { return p.Left().TypeFnc() }
func (p PairVal) TailType() TyFnc { return p.Right().TypeFnc() }

// composed functional type of a value pair
func (p PairVal) TypeFnc() TyFnc {
	return Pair | p.Left().TypeFnc() | p.Right().TypeFnc()
}

// composed native type of a value pair
func (p PairVal) TypeNat() d.TyNative {
	return p.Left().TypeNat() | p.Right().TypeNat()
}

// implements compose
func (p PairVal) Empty() bool {

	if (p.Left() == nil ||
		!p.Left().TypeFnc().Flag().Match(None) &&
			!p.Left().TypeNat().Flag().Match(d.Nil)) &&
		(p.Right() == nil ||
			!p.Right().TypeFnc().Flag().Match(None) &&
				!p.Right().TypeNat().Flag().Match(d.Nil)) {

		return true
	}

	return false
}

// call arguments are forwarded to the contained sub elements
func (p PairVal) Call(args ...Callable) Callable {
	return NewPair(p.Left().Call(args...), p.Right().Call(args...))
}

// evaluation arguments are forwarded to the contained sub elements
func (p PairVal) Eval(args ...d.Native) d.Native {
	return d.NewPair(p.Left().Eval(args...), p.Right().Eval(args...))
}

//// ASSOCIATIVE VALUE
func NewAssocPair(key string, val Callable) AssocPair {
	return func(...Callable) (string, Callable) { return key, val }
}

func (a AssocPair) KeyStr() string {
	var key, _ = a()
	return key
}

func (a AssocPair) Ident() Callable { return a }

func (a AssocPair) Pair() PairVal { return NewPair(a.Both()) }

func (a AssocPair) Both() (Callable, Callable) {
	var key, val = a()
	return NewFromData(d.StrVal(key)), val
}

func (a AssocPair) Left() Callable {
	key, _ := a()
	return NewFromData(d.StrVal(key))
}

func (a AssocPair) Right() Callable {
	_, val := a()
	return val
}
func (a AssocPair) Key() Callable   { return a.Left() }
func (a AssocPair) Value() Callable { return a.Right() }

func (a AssocPair) String() string                 { return a.Right().String() }
func (a AssocPair) Call(args ...Callable) Callable { return a.Right().Call(args...) }
func (a AssocPair) Eval(args ...d.Native) d.Native { return a.Right().Eval(args...) }

func (a AssocPair) KeyType() TyFnc           { return Pair | Symbol }
func (a AssocPair) KeyNatType() d.TyNative   { return d.String }
func (a AssocPair) ValueType() TyFnc         { return a.Right().TypeFnc() }
func (a AssocPair) ValueNatType() d.TyNative { return a.Right().TypeNat() }

func (a AssocPair) TypeFnc() TyFnc      { return Pair | Symbol | a.ValueType() }
func (a AssocPair) TypeNat() d.TyNative { return d.Pair | d.String | a.ValueNatType() }

//// ASSOCIATIVE VALUE
func NewIndexPair(idx int, val Callable) IndexPair {
	return func(...Callable) (int, Callable) { return idx, val }
}

func (a IndexPair) Index() int {
	idx, _ := a()
	return idx
}

func (a IndexPair) Ident() Callable { return a }

func (a IndexPair) Both() (Callable, Callable) {
	var idx, val = a()
	return NewFromData(d.IntVal(idx)), val
}

func (a IndexPair) Pair() PairVal { return NewPair(a.Both()) }

func (a IndexPair) Left() Callable {
	idx, _ := a()
	return NewFromData(d.IntVal(idx))
}

func (a IndexPair) Right() Callable {
	_, val := a()
	return val
}

func (a IndexPair) Key() Callable   { return a.Left() }
func (a IndexPair) Value() Callable { return a.Right() }

func (a IndexPair) String() string                 { return a.Right().String() }
func (a IndexPair) Call(args ...Callable) Callable { return a.Right().Call(args...) }
func (a IndexPair) Eval(args ...d.Native) d.Native { return a.Right().Eval(args...) }

func (a IndexPair) KeyType() TyFnc           { return Pair | Index }
func (a IndexPair) KeyNatType() d.TyNative   { return d.Int }
func (a IndexPair) ValueType() TyFnc         { return a.Right().TypeFnc() }
func (a IndexPair) ValueNatType() d.TyNative { return a.Right().TypeNat() }

func (a IndexPair) TypeFnc() TyFnc      { return Pair | Index | a.ValueType() }
func (a IndexPair) TypeNat() d.TyNative { return d.Pair | d.Int | a.ValueNatType() }

//////////////////////////////////////////////////////////////////////////////////////
///// RECURSIVE LIST OF VALUES
////
/// base implementation of recursively linked lists
//
// recursive list function holds a list of values on a late binding call by
// name base. when called without arguments, list function returns the current
// head element and a continuation, that fetches the preceeding one and returns
// it followed by another recursion. this implements the 'unit' operation of
// the list monad
//
// when arguments are passed, list function appends them on to the existing
// list lazyly. arguments stay unevaluated until they are returned. that
// constitutes the 'convert', wrap, or box operation of the  list monad.  list
// monad is the base type of many other monads
func NewList(elems ...Callable) ListVal {

	// function litereal closes over initial list and deals with arguments
	return func(args ...Callable) (Callable, ListVal) {

		var head Callable

		//  arguments are to be returned as head in preceeding calls
		//  until depleted.
		if len(args) > 0 {

			// take first argument as head to return
			head = args[0]

			// remaining arguments are parameters to generate the
			// continuation by calling new list recursively.
			if len(args) > 1 {

				// append preceeding elements from prior list
				// to set of args for list generation.
				return head, NewList(append(args[1:], elems...)...)
			}

			// last argument has been returned as head, return list
			// preceeding the call passing the arguments as tail to
			// smoothly hand over (Krueger industial smoothing‥. we
			// don't care and it shows)
			return head, NewList(elems...)
		}

		// as long as there are elements‥.
		if len(elems) > 0 {

			// ‥.assign first element to head‥.
			head = elems[0]

			// ‥.if there are further elements, return a continuation
			// to contain them‥.
			if len(elems) > 1 {
				return head, NewList(elems[1:]...)
			}

			// ‥.otherwise return last element and replace depleted
			// list, with an empty one for convienience
			return head, NewList()
		}

		// things vanished sulk about it by neither returning head nor
		// tail
		return nil, nil
	}
}

// replicates main call as method to provide type construction by appending
// elements to the list
func (l ListVal) Con(elems ...Callable) ListVal {

	return func(args ...Callable) (Callable, ListVal) {

		if len(args) > 0 {

			return l(append(elems, args...)...)
		}

		return l(elems...)
	}
}

// pushes elements at the front of the list returning the passed arguments as
// list heads until depleted, before it progresses on to return elements of the
// initial list.
func (l ListVal) Push(elems ...Callable) ListVal {

	return func(args ...Callable) (Callable, ListVal) {

		var head Callable
		var tail ListVal
		// concatenate arguments
		args = append(elems, args...)
		var la = len(args)
		var last = la - 1

		// as long as heads and tails are yielded, return them and keep
		// pushing arguments on to another call to push.
		head, tail = l()

		// if tail was yielded‥.
		if tail != nil {

			// ‥.and head was yielded as well‥.
			if head != nil {

				// return both, keep pushing on the arguments
				return head, tail.Push(args...)
			}

			// tail yielded, but no head
			if la > 0 { // one, or more arguments got passed

				if la > 1 { // two, or more arguments got passed

					// assign last argument to head,
					// re-assign remaining arguments to
					// args
					head, args = args[last], args[:last]

					// return new head and keep pushing on
					// new args using the list yielded by
					// prior call
					return head, tail.Push(args...)

				}

				// return last argument as head and return
				// yielded tail as consumable, no more pushing,
				// since arguments are depleted
				return args[last], tail
			}
		}

		// head without tail got yielded‥.
		if head != nil {

			if la > 0 {

				// use yielded head, push args on to new list
				return head, NewList().Push(args...)

			}
		}

		// neither head nor tail got yielded
		if la > 0 {

			head = args[last]

			if la > 1 {

				// use last argument as head and push remaining
				// arguments to a new list
				return head, NewList().Push(args[:last]...)
			}

			// use last argument as head
			return head, NewList()
		}

		// return nil head and a new empty list
		return nil, NewList()
	}
}

func (l ListVal) Ident() Callable { return l }

func (l ListVal) Null() ListVal { return NewList() }

func (l ListVal) Tail() Consumeable { _, t := l(); return t }

func (l ListVal) Head() Callable { h, _ := l(); return h }

func (l ListVal) DeCap() (Callable, Consumeable) { return l() }

func (l ListVal) TypeFnc() TyFnc { return List | Functor }

func (l ListVal) TypeNat() d.TyNative { return d.List.TypeNat() | l.Head().TypeNat() }

// call replicates main function call of list value instances, and either
// returns the head, when called without arguments, or concatenates the
// arguments to the list, when such are passed.
func (l ListVal) Call(d ...Callable) Callable {
	var head Callable
	head, l = l(d...)
	return head
}

// eval applys current heads eval method to passed arguments, or calle it empty
func (l ListVal) Eval(args ...d.Native) d.Native {
	return l.Head().Eval(args...)
}

func (l ListVal) Empty() bool {

	if l.Head() != nil ||
		!l.Head().TypeFnc().Flag().Match(None) &&
			!l.Head().TypeNat().Flag().Match(d.Nil) {

		return false
	}

	return true
}

// to determine the length of a recursive function, it has to be fully unwound,
// so use with care!
func (l ListVal) Len() int {
	var length int
	var head, tail = l()
	if head != nil {
		length += 1 + tail.Len()
	}
	return length
}

//////////////////////////////////////////////////////////////////////////////////////////
//// VECTORS (SLICES) OF VALUES
func NewEmptyVector(init ...Callable) VecVal { return NewVector() }

func NewVector(init ...Callable) VecVal {

	var vector = init

	return func(args ...Callable) []Callable {

		if len(args) > 0 {

			// append args to vector
			vector = append(
				vector,
				args...,
			)
		}

		// return slice vector
		return vector
	}
}

func ConVector(vec Vectorized, args ...Callable) VecVal {

	return ConVecFromCallable(append(reverseArguments(args...), vec.Slice()...)...)
}

func AppendVector(vec Vectorized, args ...Callable) VecVal {

	return ConVecFromCallable(append(vec.Slice(), args...)...)

}

func ConVecFromCallable(init ...Callable) VecVal {

	return func(args ...Callable) []Callable {

		return reverseArguments(append(args, init...)...)
	}
}

func AppendVecFromCallable(init ...Callable) VecVal {

	return func(args ...Callable) []Callable {

		return append(init, args...)
	}
}

func (v VecVal) Ident() Callable { return v }

func (v VecVal) Call(d ...Callable) Callable { return NewVector(v(d...)...) }

func (v VecVal) Eval(args ...d.Native) d.Native {

	var result = []d.Native{}

	for _, arg := range args {
		result = append(result, arg)
	}

	return d.DataSlice(result)
}

func (v VecVal) TypeFnc() TyFnc {
	if len(v()) > 0 {
		return Vector | v.Head().TypeFnc()
	}
	return Vector | None
}

func (v VecVal) TypeNat() d.TyNative {
	if len(v()) > 0 {
		return d.Vector.TypeNat() | v.Head().TypeNat()
	}
	return d.Vector.TypeNat() | d.Nil
}

func (v VecVal) Head() Callable {
	if v.Len() > 0 {
		return v.Slice()[0]
	}
	return nil
}

func (v VecVal) Tail() Consumeable {
	if v.Len() > 1 {
		return NewVector(v.Slice()[1:]...)
	}
	return NewEmptyVector()
}

func (v VecVal) DeCap() (Callable, Consumeable) {
	return v.Head(), v.Tail()
}

func (v VecVal) Empty() bool {

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

func (v VecVal) Len() int { return len(v()) }

func (v VecVal) Vector() VecVal { return v }

func (v VecVal) Slice() []Callable { return v() }

func (v VecVal) Append(args ...Callable) VecVal {
	return NewVector(append(v(), args...)...)
}

func (v VecVal) Con(args ...Callable) VecVal {
	return NewVector(append(reverseArguments(args...), reverseArguments(v()...)...)...)
}

func (v VecVal) Set(i int, val Callable) Vectorized {
	if i < v.Len() {
		var slice = v()
		slice[i] = val
		return VecVal(func(elems ...Callable) []Callable { return slice })

	}
	return v
}

func (v VecVal) Get(i int) Callable {
	if i < v.Len() {
		return v()[i]
	}
	return NewNoOp()
}
func (v VecVal) Search(praed Callable) int { return newDataSorter(v()...).Search(praed) }

func (v VecVal) Sort(flag d.TyNative) {
	var ps = newDataSorter(v()...)
	ps.Sort(flag)
	v = NewVector(ps...)
}

/////////////////////////////////////////////////////////////////////////////////////
//// TUPLE TYPE VALUES
///
// tuples are sequences of values grouped in a distinct sequence of distinct types,
func NewTuple(data ...Callable) TupleVal {

	return TupleVal(func(args ...Callable) []Callable {

		return data
	})
}

func (t TupleVal) Ident() Callable { return t }
func (t TupleVal) Len() int        { return len(t()) }

// pairs prepends annotates member values as pair values carrying this
// instances sub-type signature and tuple position in in the second field
func (t TupleVal) Pairs() []PairVal {

	var pairs = []PairVal{}

	return pairs
}

// implement consumeable
func (t TupleVal) DeCap() (Callable, Consumeable) {

	var head Callable
	var list = NewList()

	for _, pair := range t.Pairs() {
		head, list = list(pair)
	}

	return head, list
}

func (t TupleVal) Head() Callable    { head, _ := t.DeCap(); return head }
func (t TupleVal) Tail() Consumeable { _, tail := t.DeCap(); return tail }

// functional type concatenates the functional types of all the subtypes
func (t TupleVal) TypeFnc() TyFnc {
	var ftype = TyFnc(0)
	for _, typ := range t() {
		ftype = ftype | typ.TypeFnc()
	}
	return ftype
}

// native type concatenates the native types of all the subtypes
func (t TupleVal) TypeNat() d.TyNative {
	var ntype = d.Tuple
	for _, typ := range t() {
		ntype = ntype | typ.TypeNat()
	}
	return ntype
}

// string representation of a tuple generates one row per sub type by
// concatenating each sub types native type, functional type and value.
func (t TupleVal) String() string { return t.Head().String() }

func (t TupleVal) Eval(args ...d.Native) d.Native {
	var result = []d.Native{}
	for _, val := range t() {
		result = append(result, val.Eval(val))
	}
	return d.DataSlice(result)
}

func (t TupleVal) Call(args ...Callable) Callable {
	var result []Callable
	for _, val := range t() {
		result = append(result, val.Call(args...))
	}
	return NewVector(result...)
}

func (t TupleVal) ApplyPartial(args ...Callable) TupleVal {
	return NewTuple(partialApplyTuple(t, args...)()...)
}

func partialApplyTuple(tuple TupleVal, args ...Callable) TupleVal {
	// fetch current tupple
	var result = tuple()
	var l = len(result)

	// range through arguments
	for i := 0; i < l; i++ {

		// pick argument by index
		var arg = args[i]

		// partial arguments can either be given by position, or in
		// pairs that contains the intendet position as integer value
		// in its left and the value itself in its right cell, so‥.
		if pair, ok := arg.(PairVal); ok {
			// ‥.and the left element is an integer‥.
			if pos, ok := pair.Left().(Integer); ok {
				// ‥.and that integer is within the range of indices‥.
				if l < pos.Int() {
					// ‥.and both types of the right element
					// match the corresponding result types
					// of the given index‥.
					if result[i].TypeFnc() == pair.Right().TypeFnc() &&
						result[i].TypeNat() == args[i].TypeNat() {
						// ‥.replace the value in
						// results, with right
						// element of pair.
						result[i] = pair.Right()
					}
				}
			}
		}
		// ‥.otherwise assume arguments are passed one element at a
		// time, altering between position & value and the current
		// index is expected to be the position, so if it's an uneven
		// index (positions)‥.
		if i%2 == 0 {
			var idx = i  // safe current index
			if i+1 < l { // check if next index is out of bounds
				i = i + 1 // advance loop counter by one
				// replace value in results at previous index
				// with value at index of the advanced loop
				// counter
				result[idx] = args[i]
			}
		}
	}
	// return altered result
	return TupleVal(
		func(...Callable) []Callable {
			return result
		})
}

//// ASSOCIATIVE SLICE OF VALUE PAIRS
///
// list of associative values in predefined order.
func ConAssociative(vec Associative, pfnc ...PairVal) AccociativeVal {
	return NewAssociativeFromPairFunction(append(vec.Pairs(), pfnc...)...)
}

func NewAssociativeFromPairFunction(ps ...PairVal) AccociativeVal {
	var pairs = []PairVal{}
	for _, pair := range ps {
		pairs = append(pairs, pair)
	}
	return AccociativeVal(func(pairs ...PairVal) []PairVal { return pairs })
}

func ConAssociativeFromPairs(pp ...PairVal) AccociativeVal {
	return AccociativeVal(func(pairs ...PairVal) []PairVal { return pp })
}

func NewEmptyAssociative() AccociativeVal {
	return AccociativeVal(func(pairs ...PairVal) []PairVal { return []PairVal{} })
}

func NewAssociative(pp ...PairVal) AccociativeVal {

	return func(pairs ...PairVal) []PairVal {
		for _, pair := range pp {
			pairs = append(pairs, pair)
		}
		return pairs
	}
}

func (v AccociativeVal) Call(d ...Callable) Callable {
	if len(d) > 0 {
		for _, val := range d {
			if pair, ok := val.(PairVal); ok {
				v = v.Con(pair)
			}
		}
	}
	return v
}

func (v AccociativeVal) Con(p ...Callable) AccociativeVal {

	var pairs = v.Pairs()

	return ConAssociativeFromPairs(pairs...)
}

func (v AccociativeVal) DeCap() (Callable, Consumeable) {
	return v.Head(), v.Tail()
}

func (v AccociativeVal) Empty() bool {

	if len(v()) > 0 {

		for _, pair := range v() {

			if !pair.Empty() {

				return false
			}
		}
	}
	return true
}

func (v AccociativeVal) KeyFncType() TyFnc {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().TypeFnc()
	}
	return None
}

func (v AccociativeVal) KeyNatType() d.TyNative {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().TypeNat()
	}
	return d.Nil
}

func (v AccociativeVal) ValFncType() TyFnc {
	if v.Len() > 0 {
		return v.Pairs()[0].Right().TypeFnc()
	}
	return None
}

func (v AccociativeVal) ValNatType() d.TyNative {
	if v.Len() > 0 {
		return v.Pairs()[0].Right().TypeNat()
	}
	return d.Nil
}

func (v AccociativeVal) TypeFnc() TyFnc { return Record | Functor }

func (v AccociativeVal) TypeNat() d.TyNative {
	if len(v()) > 0 {
		return d.Vector | v.Head().TypeNat()
	}
	return d.Vector | d.Nil.TypeNat()
}

func (v AccociativeVal) Len() int { return len(v()) }

func (v AccociativeVal) Get(idx int) PairVal {
	if idx < v.Len()-1 {
		return v()[idx]
	}
	return NewPair(NewNoOp(), NewNoOp())
}

func (v AccociativeVal) GetVal(praed Callable) PairVal {
	return newPairSorter(v()...).Get(praed)
}

func (v AccociativeVal) Range(praed Callable) []PairVal {
	return newPairSorter(v()...).Range(praed)
}

func (v AccociativeVal) Search(praed Callable) int {
	return newPairSorter(v()...).Search(praed)
}

func (v AccociativeVal) Pairs() []PairVal { return v() }

func (v AccociativeVal) DeCapPairWise() (PairVal, []PairVal) {

	var pairs = v()

	if len(pairs) > 0 {
		if len(pairs) > 1 {
			return pairs[0], pairs[1:]
		}
		return pairs[0], []PairVal{}
	}
	return nil, []PairVal{}
}

func (v AccociativeVal) SwitchedPairs() []PairVal {
	var switched = []PairVal{}
	for _, pair := range v() {
		switched = append(
			switched,
			NewPair(
				pair.Right(),
				pair.Left()))
	}
	return switched
}

func (v AccociativeVal) SetVal(key, value Callable) Associative {
	if idx := v.Search(key); idx >= 0 {
		var pairs = v.Pairs()
		pairs[idx] = NewPair(key, value)
		return NewAssociative(pairs...)
	}
	return NewAssociative(append(v.Pairs(), NewPair(key, value))...)
}

func (v AccociativeVal) Slice() []Callable {
	var fncs = []Callable{}
	for _, pair := range v() {
		fncs = append(fncs, NewPair(pair.Left(), pair.Right()))
	}
	return fncs
}

func (v AccociativeVal) Head() Callable {
	if v.Len() > 0 {
		return v.Pairs()[0]
	}
	return nil
}

func (v AccociativeVal) Tail() Consumeable {
	if v.Len() > 1 {
		return ConAssociativeFromPairs(v.Pairs()[1:]...)
	}
	return NewEmptyAssociative()
}

func (v AccociativeVal) Eval(p ...d.Native) d.Native {
	var slice = d.DataSlice{}
	for _, pair := range v() {
		d.SliceAppend(slice, d.NewPair(pair.Left(), pair.Right()))
	}
	return slice
}

func (v AccociativeVal) Sort(flag d.TyNative) {
	var ps = newPairSorter(v()...)
	ps.Sort(flag)
	v = NewAssociative(ps...)
}

//////////////////////////////////////////////////////////////////////////////
//// ASSOCIATIVE SET (HASH MAP OF VALUES)
///
// unordered associative set of values
func ConAssocSet(pairs ...PairVal) SetVal {
	var paired = []PairVal{}
	for _, pair := range pairs {
		paired = append(paired, pair)
	}
	return NewAssocSet(paired...)
}

func NewAssocSet(pairs ...PairVal) SetVal {

	var kt d.TyNative
	var set d.Mapped

	// OR concat all accessor types
	for _, pair := range pairs {
		kt = kt | pair.KeyNatType()
	}
	// if accessors are of mixed type‥.
	if kt.Flag().Count() > 1 {
		set = d.SetVal{}
	} else {
		var ktf = kt.Flag()
		switch {
		case ktf.Match(d.Int):
			set = d.SetInt{}
		case ktf.Match(d.Uint):
			set = d.SetUint{}
		case ktf.Match(d.Flag):
			set = d.SetFlag{}
		case ktf.Match(d.Float):
			set = d.SetFloat{}
		case ktf.Match(d.String):
			set = d.SetString{}
		}
	}
	return SetVal(func(pairs ...PairVal) d.Mapped { return set })
}

func (v SetVal) Split() (VecVal, VecVal) {
	var keys, vals = []Callable{}, []Callable{}
	for _, pair := range v.Pairs() {
		keys = append(keys, pair.Left())
		vals = append(vals, pair.Right())
	}
	return NewVector(keys...), NewVector(vals...)
}

func (v SetVal) Pairs() []PairVal {
	var pairs = []PairVal{}
	for _, field := range v().Fields() {
		pairs = append(
			pairs,
			NewPairFromData(
				field.Left(),
				field.Right()))
	}
	return pairs
}

func (v SetVal) Keys() VecVal { k, _ := v.Split(); return k }

func (v SetVal) Data() VecVal { _, d := v.Split(); return d }

func (v SetVal) Len() int { return v().Len() }

func (v SetVal) Empty() bool {

	for _, pair := range v.Pairs() {

		if !pair.Empty() {

			return false
		}
	}

	return true
}

func (v SetVal) GetVal(praed Callable) PairVal {
	var val Callable
	var nat, ok = v().Get(praed)
	if val, ok = nat.(Callable); !ok {
		val = NewFromData(val)
	}
	return NewPair(praed, val)
}

func (v SetVal) SetVal(key, value Callable) Associative {
	var m = v()
	m.Set(key, value)
	return SetVal(func(pairs ...PairVal) d.Mapped { return m })
}

func (v SetVal) Slice() []Callable {
	var pairs = []Callable{}
	for _, pair := range v.Pairs() {
		pairs = append(pairs, pair)
	}
	return pairs
}

func (v SetVal) Call(f ...Callable) Callable { return v }

func (v SetVal) Eval(p ...d.Native) d.Native {
	var slice = d.DataSlice{}
	for _, pair := range v().Fields() {
		d.SliceAppend(slice, d.NewPair(pair.Left(), pair.Right()))
	}
	return slice
}

func (v SetVal) TypeFnc() TyFnc { return Set | Functor }

func (v SetVal) TypeNat() d.TyNative { return d.Set | d.Function }

func (v SetVal) KeyFncType() TyFnc {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().TypeFnc()
	}
	return None
}

func (v SetVal) KeyNatType() d.TyNative {
	if v.Len() > 0 {
		return v.Pairs()[0].Left().TypeNat()
	}
	return d.Nil
}

func (v SetVal) ValFncType() TyFnc {
	if v.Len() > 0 {
		return v.Pairs()[0].Right().TypeFnc()
	}
	return None
}

func (v SetVal) ValNatType() d.TyNative {
	if v.Len() > 0 {
		return v.Pairs()[0].Right().TypeNat()
	}
	return d.Nil
}

func (v SetVal) DeCap() (Callable, Consumeable) {
	return v.Head(), v.Tail()
}

func (v SetVal) Head() Callable {
	if v.Len() > 0 {
		return v.Pairs()[0]
	}
	return nil
}

func (v SetVal) Tail() Consumeable {
	if v.Len() > 1 {
		return ConAssociativeFromPairs(v.Pairs()[1:]...)
	}
	return NewEmptyAssociative()
}