/*
  type aliases from data package are wrapped by data expressions to implement
  the callable interface.

  there are several types of data expressions:

  - data constant wraps all static data types
  - data slice wraps slices of native instances
  - data go-slice are slices of instances of native go types
  - data pair is the package data implementation of a value pair
  - data sets have one implementation per key type for a variety of keytypes
  - data expression is a generic function with a signature expecting one or
    many instances of package data native instances as its arguments and one
    data/native instance as return value
*/
package functions

import d "github.com/joergreinhardt/gatwd/data"

type (
	//// NATIVE VALUE CONSTRUCTORS
	DataConst   func() d.Native
	DataSlice   func() d.DataSlice
	DataGoSlice func() d.Sliceable
	DataPair    func() d.PairVal
	DataSet     func() d.Mapped
	DataValue   func(...d.Native) d.Native
)

func nativeType(nat d.Native) (typed d.Typed) {
	switch {
	case nat.Type().Match(d.Pair):
		var p = nat.(d.PairVal)
		typed = Def(p.TypeKey(), p.TypeValue())
	case nat.Type().Match(d.Unboxed):
		var u = nat.(d.Sliceable)
		typed = Def(u.Type(), u.TypeElem())
	case nat.Type().Match(d.Slice):
		var s = nat.(d.Sliceable)
		typed = Def(s.Type(), s.TypeElem())
	case nat.Type().Match(d.Map):
		var m = nat.(d.Mapped)
		typed = Def(m.Type(), m.TypeKey(), m.TypeValue())
	default:
		typed = nat.Type()
	}
	return typed
}

//// DATA CONSTRUCTOR
///
// returns an expression with native return type implementing the callable
// interface
func Declare(inf ...interface{}) d.Native {
	return d.New(inf...)
}

func DecNative(inf ...interface{}) Native {
	return DecData(Declare(inf...))
}

func DecData(args ...d.Native) Native {
	var (
		nat   = d.NewData(args...)
		match = nat.Type().Match
	)

	switch {
	case match(d.Function):
		if fn, ok := nat.(d.FuncVal); ok {
			return DataValue(func(args ...d.Native) d.Native {
				if len(args) > 0 {
					return fn(args...)
				}
				return fn()
			})
		}
	case match(d.Slice):
		return DataSlice(func() d.DataSlice {
			return nat.(d.DataSlice)
		})
	case match(d.Unboxed):
		return DataGoSlice(func() d.Sliceable {
			return nat.(d.Sliceable)
		})
	case match(d.Pair):
		return DataPair(func() d.PairVal {
			return nat.(d.PairVal)
		})
	case match(d.Map):
		return DataSet(func() d.Mapped {
			return nat.(d.Mapped)
		})
	}
	return DataConst(func() d.Native {
		return nat
	})
}

// NATIVE FUNCTION VALUE CONSTRUCTOR
func (n DataValue) TypeFnc() TyFnc                 { return Data }
func (n DataValue) TypeNat() d.TyNat               { return n().Type() }
func (n DataValue) String() string                 { return n.Eval().String() }
func (n DataValue) Eval(args ...d.Native) d.Native { return n(args...) }
func (n DataValue) Call(args ...Expression) Expression {
	if len(args) > 0 {
		var nats = make([]d.Native, 0, len(args))
		for _, arg := range args {
			if arg.TypeFnc().Match(Data) {
				if data, ok := arg.(Native); ok {
					var eval = data.Eval()
					nats = append(nats, eval)
				}
			}
		}
		return DecData(n(nats...))
	}
	return DecData(n())
}
func (n DataValue) Type() TyPattern {
	return Def(Data, Def(Data, Value), nativeType(n()))
}

// ATOMIC NATIVE VALUE CONSTRUCTOR
func (n DataConst) Eval(...d.Native) d.Native     { return n() }
func (n DataConst) TypeFnc() TyFnc                { return Data }
func (n DataConst) TypeNat() d.TyNat              { return n().Type() }
func (n DataConst) String() string                { return n().String() }
func (n DataConst) Call(...Expression) Expression { return DecData(n()) }
func (n DataConst) Type() TyPattern {
	return Def(None, Def(Data, Constant), nativeType(n()))
}

// NATIVE SLICE VALUE CONSTRUCTOR
func (n DataSlice) Call(args ...Expression) Expression { return n }
func (n DataSlice) TypeFnc() TyFnc                     { return Data }
func (n DataSlice) TypeNat() d.TyNat                   { return n().Type() }
func (n DataSlice) Len() int                           { return n().Len() }
func (n DataSlice) Head() d.Native                     { return n().Head() }
func (n DataSlice) Tail() d.Sequential                 { return n().Tail() }
func (n DataSlice) Shift() (d.Native, d.DataSlice)     { return n().Shift() }
func (n DataSlice) SliceNat() []d.Native               { return n().Slice() }
func (n DataSlice) Get(key d.Native) d.Native          { return n().Get(key) }
func (n DataSlice) GetInt(idx int) d.Native            { return n().GetInt(idx) }
func (n DataSlice) Range(s, e int) d.Sliceable         { return n().Range(s, e) }
func (n DataSlice) Empty() bool                        { return n().Empty() }
func (n DataSlice) Copy() d.Native                     { return n().Copy() }
func (n DataSlice) ElemType() d.Typed                  { return n().TypeElem() }
func (n DataSlice) String() string                     { return n().String() }
func (n DataSlice) Slice() []d.Native                  { return n().Slice() }
func (n DataSlice) Type() TyPattern {
	return Def(nativeType(n()), Def(Data, Vector), nativeType(n()))
}
func (n DataSlice) Eval(args ...d.Native) d.Native {
	return d.SliceAppend(n(), args...)
}
func (n DataSlice) SliceExpr() []Expression {
	var slice = make([]Expression, 0, n.Len())
	for _, nat := range n.Slice() {
		slice = append(slice, DecData(nat))
	}
	return slice
}

func (n DataGoSlice) Call(args ...Expression) Expression { return n }
func (n DataGoSlice) Eval(args ...d.Native) d.Native {
	return d.NewSlice(append(n.Slice(), args...)...)
}
func (n DataGoSlice) TypeFnc() TyFnc             { return Data }
func (n DataGoSlice) TypeNat() d.TyNat           { return n().Type() }
func (n DataGoSlice) Len() int                   { return n().Len() }
func (n DataGoSlice) Get(key d.Native) d.Native  { return n().Get(key) }
func (n DataGoSlice) GetInt(idx int) d.Native    { return n().GetInt(idx) }
func (n DataGoSlice) Range(s, e int) d.Sliceable { return n().Range(s, e) }
func (n DataGoSlice) Copy() d.Native             { return n().Copy() }
func (n DataGoSlice) Empty() bool                { return n().Empty() }
func (n DataGoSlice) Slice() []d.Native          { return n().Slice() }
func (n DataGoSlice) ElemType() d.Typed          { return n().TypeElem() }
func (n DataGoSlice) String() string             { return n().String() }
func (n DataGoSlice) Type() TyPattern {
	return Def(nativeType(n()), Def(Data, Vector), nativeType(n()))
}
func (n DataGoSlice) SliceExpr() []Expression {
	var slice = make([]Expression, 0, n.Len())
	for _, nat := range n.Slice() {
		slice = append(slice, DecData(nat))
	}
	return slice
}

// NATIVE PAIR VALUE CONSTRUCTOR
func (n DataPair) Call(args ...Expression) Expression { return n }
func (n DataPair) Eval(...d.Native) d.Native          { return n() }
func (n DataPair) TypeFnc() TyFnc                     { return Data }
func (n DataPair) TypeNat() d.TyNat                   { return n().Type() }
func (n DataPair) Left() d.Native                     { return n().Left() }
func (n DataPair) Right() d.Native                    { return n().Right() }
func (n DataPair) Both() (l, r d.Native)              { return n().Both() }
func (n DataPair) LeftType() d.TyNat                  { return n().TypeKey() }
func (n DataPair) RightType() d.TyNat                 { return n().TypeValue() }
func (n DataPair) SubType() d.Typed                   { return n().Type() }
func (n DataPair) String() string                     { return n().String() }
func (n DataPair) LeftExpr() Expression               { return DecData(n().Left()) }
func (n DataPair) RightExpr() Expression              { return DecData(n().Right()) }
func (n DataPair) BothExpr() (l, r Expression) {
	return DecData(n().Left()),
		DecData(n().Right())
}
func (n DataPair) Pair() Paired {
	return NewPair(
		DecData(n().Left()),
		DecData(n().Right()))
}
func (n DataPair) Type() TyPattern {
	return Def(
		Def(
			nativeType(n().Left()),
			nativeType(n().Right()),
		),
		Def(Data, Pair),
		Def(
			nativeType(n().Left()),
			nativeType(n().Right()),
		),
	)
}

// NATIVE SET VALUE CONSTRUCTOR
func (n DataSet) Call(args ...Expression) Expression   { return n }
func (n DataSet) Eval(...d.Native) d.Native            { return n() }
func (n DataSet) TypeFnc() TyFnc                       { return Data }
func (n DataSet) TypeNat() d.TyNat                     { return n().Type() }
func (n DataSet) Len() int                             { return n().Len() }
func (n DataSet) Slice() []d.Native                    { return n().Slice() }
func (n DataSet) GetNat(acc d.Native) (d.Native, bool) { return n().Get(acc) }
func (n DataSet) SetNat(acc, val d.Native) d.Mapped    { return n().Set(acc, val) }
func (n DataSet) Delete(acc d.Native) bool             { return n().Delete(acc) }
func (n DataSet) Get(acc d.Native) (d.Native, bool)    { return n().Get(acc) }
func (n DataSet) Set(acc, val d.Native) d.Mapped       { return n().Set(acc, val) }
func (n DataSet) Keys() []d.Native                     { return n().Keys() }
func (n DataSet) Data() []d.Native                     { return n().Data() }
func (n DataSet) Fields() []d.Paired                   { return n().Fields() }
func (n DataSet) KeyType() d.Typed                     { return n().TypeKey() }
func (n DataSet) ValType() d.Typed                     { return n().TypeValue() }
func (n DataSet) SubType() d.Typed                     { return n().Type() }
func (n DataSet) String() string                       { return n().String() }
func (n DataSet) KeysExpr() []Expression {
	var exprs = make([]Expression, 0, n.Len())
	for _, key := range n().Keys() {
		exprs = append(exprs, DecData(key))
	}
	return exprs
}
func (n DataSet) DataExpr() []Expression {
	var exprs = make([]Expression, 0, n.Len())
	for _, val := range n().Data() {
		exprs = append(exprs, DecData(val))
	}
	return exprs
}
func (n DataSet) SliceExpr() []Expression {
	var slice = make([]Expression, 0, n.Len())
	for _, nat := range n.Fields() {
		slice = append(slice, DecData(nat))
	}
	return slice
}
func (n DataSet) Pairs() []Paired {
	var pairs = []Paired{}
	for _, field := range n.Fields() {
		pairs = append(
			pairs, NewPair(
				DecData(field.Left()),
				DecData(field.Right())))
	}
	return pairs
}
func (n DataSet) Type() TyPattern {
	if n().Len() > 0 {
		return Def(
			Def(
				nativeType(n().First().Left()),
				nativeType(n().First().Right()),
			),
			Def(Data, Pair),
			Def(
				nativeType(n().First().Left()),
				nativeType(n().First().Right()),
			),
		)
	}
	return Def(Data, Pair)
}
