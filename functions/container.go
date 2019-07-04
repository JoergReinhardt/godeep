/*
  FUNCTIONAL CONTAINERS

  containers implement enumeration of functional types, aka lists, vectors
  sets, pairs, tuples‥.
*/
package functions

import (
	d "github.com/joergreinhardt/gatwd/data"
)

type (
	//// CONSTANT & GENERIC EXPRESSIONS
	ConstantExpr func() Expression
	GenericExpr  func(...Expression) Expression

	//// NATIVE VALUE CONSTRUCTORS
	NativeExpr func(...d.Native) d.Native
	NativePair func(...d.Native) d.Paired
	NativeUbox func(...d.Native) d.Sliceable
	NativeCol  func(...d.Native) d.Sliceable
	NativeSet  func(...d.Native) d.Mapped

	//// EXPRESSION VALUE CONSTRUCTOR
	PartialExpr func(...Expression) Expression
)

/// CONSTANT VALUE CONSTRUCTOR
//
func NewConstant(constant func() Expression) ConstantExpr { return constant }

func (c ConstantExpr) Ident() Expression                  { return c }
func (c ConstantExpr) Call(args ...Expression) Expression { return c() }
func (c ConstantExpr) Arity() Arity                       { return Arity(0) }
func (c ConstantExpr) TypeFnc() TyFnc                     { return Constant }
func (c ConstantExpr) TypeNat() d.TyNat                   { return c().TypeNat() }
func (c ConstantExpr) String() string                     { return c().String() }
func (c ConstantExpr) Eval(args ...d.Native) d.Native     { return c().Eval(args...) }
func (c ConstantExpr) FlagType() d.Uint8Val               { return Flag_Functional.U() }
func (c ConstantExpr) TypeName() string                   { return c().TypeName() }
func (c ConstantExpr) Type() TyDef {
	return Define("ϝ → "+c().TypeName(), c().TypeFnc())
}

func NewGeneric(expr func(args ...Expression) Expression) GenericExpr { return expr }

func (c GenericExpr) Ident() Expression                  { return c }
func (c GenericExpr) Type() TyDef                        { return c().Type() }
func (c GenericExpr) TypeFnc() TyFnc                     { return c().TypeFnc() }
func (c GenericExpr) TypeNat() d.TyNat                   { return c().TypeNat() }
func (c GenericExpr) String() string                     { return c().String() }
func (c GenericExpr) TypeName() string                   { return c().TypeName() }
func (c GenericExpr) FlagType() d.Uint8Val               { return Flag_Functional.U() }
func (c GenericExpr) Call(args ...Expression) Expression { return c(args...) }
func (c GenericExpr) Eval(args ...d.Native) d.Native {
	var exprs = make([]Expression, 0, len(args))
	for _, arg := range args {
		exprs = append(exprs, NewNative(arg))
	}
	return c(exprs...).Eval()
}

//// NATIVE EXPRESSION CONSTRUCTOR
///
// returns an expression with native return type implementing the callable
// interface
func New(inf ...interface{}) Expression {
	return NewNative(d.New(inf...))
}

// TODO replace expression wrappers by d.Native.Eval
func NewNative(args ...d.Native) Expression {
	// if any initial arguments have been passed

	var nat = d.NewFromData(args...)
	var tnat = nat.TypeNat()

	switch {
	case tnat.Match(d.Slice):
		if slice, ok := nat.(d.DataSlice); ok {
			return NativeCol(func(args ...d.Native) d.Sliceable {
				if len(args) > 0 {
					return slice.Eval(args...).(d.Sliceable)
				}
				return slice
			})
		}
	case tnat.Match(d.Unboxed):
		return NativeUbox(func(args ...d.Native) d.Sliceable {
			if len(args) > 0 {
				return nat.Eval(args...).(d.Sliceable)
			}
			return nat.(d.Sliceable)
		})
	case tnat.Match(d.Pair):
		if pair, ok := nat.(d.Paired); ok {
			return NativePair(func(args ...d.Native) d.Paired {
				if len(args) > 0 {
					return pair.Eval(args...).(d.Paired)
				}
				return pair
			})
		}
	case tnat.Match(d.Map):
		if set, ok := nat.(d.Mapped); ok {
			return NativeSet(func(args ...d.Native) d.Mapped {
				if len(args) > 0 {
					return set.Eval(args...).(d.Mapped)
				}
				return set
			})
		}
	default:
		return NativeExpr(nat.Eval)
	}
	return NativeExpr(func(...d.Native) d.Native { return d.NewNil() })
}

// ATOMIC NATIVE VALUE CONSTRUCTOR
func (n NativeExpr) Call(args ...Expression) Expression {
	var nats = make([]d.Native, 0, len(args))
	for _, arg := range args {
		nats = append(nats, arg.Eval())
	}
	return NewNative(n(nats...))
}
func (n NativeExpr) TypeFnc() TyFnc                 { return Data }
func (n NativeExpr) Eval(args ...d.Native) d.Native { return n(args...) }
func (n NativeExpr) TypeNat() d.TyNat               { return n().TypeNat() }
func (n NativeExpr) FlagType() d.Uint8Val           { return Flag_Functional.U() }
func (n NativeExpr) String() string                 { return n().String() }
func (n NativeExpr) TypeName() string               { return n().TypeName() }
func (n NativeExpr) Type() TyDef {
	return Define(n().TypeName(), New(n().TypeNat()))
}

// NATIVE SLICE VALUE CONSTRUCTOR
func (n NativeCol) Call(args ...Expression) Expression {
	if len(args) > 0 {
		var nats = make([]d.Native, 0, len(args))
		for _, arg := range args {
			nats = append(nats, arg.Eval())
		}
		return NewNative(n(nats...))
	}
	return NewNative(n())
}
func (n NativeCol) TypeFnc() TyFnc                 { return Data }
func (n NativeCol) Eval(args ...d.Native) d.Native { return n(args...) }
func (n NativeCol) Len() int                       { return n().Len() }
func (n NativeCol) SliceNat() []d.Native           { return n().Slice() }
func (n NativeCol) Get(key d.Native) d.Native      { return n().Get(key) }
func (n NativeCol) GetInt(idx int) d.Native        { return n().GetInt(idx) }
func (n NativeCol) Range(s, e int) d.Native        { return n().Range(s, e) }
func (n NativeCol) Copy() d.Native                 { return n().Copy() }
func (n NativeCol) TypeNat() d.TyNat               { return n().TypeNat() }
func (n NativeCol) String() string                 { return n().String() }
func (n NativeCol) TypeName() string               { return n().TypeName() }
func (n NativeCol) Vector() VecCol                 { return NewVector(n.Slice()...) }
func (n NativeCol) FlagType() d.Uint8Val           { return Flag_Functional.U() }
func (n NativeCol) Type() TyDef {
	return Define(n().TypeName(), New(n().TypeNat()))
}
func (n NativeCol) Slice() []Expression {
	var slice = []Expression{}
	for _, val := range n.SliceNat() {
		slice = append(slice, NewNative(val))
	}
	return slice
}

func (n NativeUbox) Call(args ...Expression) Expression {
	if len(args) > 0 {
		var nats = make([]d.Native, 0, len(args))
		for _, arg := range args {
			nats = append(nats, arg.Eval())
		}
		return NewNative(n(nats...))
	}
	return NewNative(n())
}
func (n NativeUbox) TypeFnc() TyFnc                 { return Data }
func (n NativeUbox) Eval(args ...d.Native) d.Native { return n(args...) }
func (n NativeUbox) Len() int                       { return n().Len() }
func (n NativeUbox) SliceNat() []d.Native           { return n().Slice() }
func (n NativeUbox) Get(key d.Native) d.Native      { return n().Get(key) }
func (n NativeUbox) GetInt(idx int) d.Native        { return n().GetInt(idx) }
func (n NativeUbox) Range(s, e int) d.Native        { return n().Range(s, e) }
func (n NativeUbox) Copy() d.Native                 { return n().Copy() }
func (n NativeUbox) TypeNat() d.TyNat               { return n().TypeNat() }
func (n NativeUbox) String() string                 { return n().String() }
func (n NativeUbox) TypeName() string               { return n().TypeName() }
func (n NativeUbox) Vector() VecCol                 { return NewVector(n.Slice()...) }
func (n NativeUbox) FlagType() d.Uint8Val           { return Flag_Functional.U() }
func (n NativeUbox) Type() TyDef {
	return Define(n().TypeName(), New(n().TypeNat()))
}
func (n NativeUbox) Slice() []Expression {
	var slice = []Expression{}
	for _, val := range n.Slice() {
		slice = append(slice, NewNative(val))
	}
	return slice
}

// NATIVE PAIR VALUE CONSTRUCTOR
func (n NativePair) Call(args ...Expression) Expression {
	var nats = make([]d.Native, 0, len(args))
	for _, arg := range args {
		nats = append(nats, arg.Eval())
	}
	return NewNative(n(nats...))
}
func (n NativePair) TypeFnc() TyFnc                 { return Data }
func (n NativePair) TypeNat() d.TyNat               { return n().TypeNat() }
func (n NativePair) Eval(args ...d.Native) d.Native { return n(args...) }
func (n NativePair) LeftNat() d.Native              { return n().Left() }
func (n NativePair) RightNat() d.Native             { return n().Right() }
func (n NativePair) BothNat() (l, r d.Native)       { return n().Both() }
func (n NativePair) Left() Expression               { return NewNative(n().Left()) }
func (n NativePair) Right() Expression              { return NewNative(n().Right()) }
func (n NativePair) KeyType() d.TyNat               { return n().LeftType() }
func (n NativePair) ValType() d.TyNat               { return n().RightType() }
func (n NativePair) SubType() d.Typed               { return n().TypeNat() }
func (n NativePair) TypeName() string               { return n().TypeName() }
func (n NativePair) FlagType() d.Uint8Val           { return Flag_Functional.U() }
func (n NativePair) String() string                 { return n().String() }
func (n NativePair) Type() TyDef {
	return Define(n().TypeName(), New(n().TypeNat()))
}
func (n NativePair) Pair() Paired {
	return NewPair(
		NewNative(n().Left()),
		NewNative(n().Right()))
}
func (n NativePair) Both() (l, r Expression) {
	return NewNative(n().Left()),
		NewNative(n().Right())
}

// NATIVE SET VALUE CONSTRUCTOR

func (n NativeSet) Call(args ...Expression) Expression {
	var nats = make([]d.Native, 0, len(args))
	for _, arg := range args {
		nats = append(nats, arg.Eval())
	}
	return NewNative(n(nats...))
}
func (n NativeSet) Ident() Expression                    { return n }
func (n NativeSet) Eval(args ...d.Native) d.Native       { return n(args...) }
func (n NativeSet) TypeFnc() TyFnc                       { return Data }
func (n NativeSet) TypeNat() d.TyNat                     { return n().TypeNat() }
func (n NativeSet) GetNat(acc d.Native) (d.Native, bool) { return n().Get(acc) }
func (n NativeSet) SetNat(acc, val d.Native) d.Mapped    { return n().Set(acc, val) }
func (n NativeSet) Delete(acc d.Native) bool             { return n().Delete(acc) }
func (n NativeSet) KeysNat() []d.Native                  { return n().Keys() }
func (n NativeSet) DataNat() []d.Native                  { return n().Data() }
func (n NativeSet) Fields() []d.Paired                   { return n().Fields() }
func (n NativeSet) KeyTypeNat() d.TyNat                  { return n().KeyType() }
func (n NativeSet) ValTypeNat() d.TyNat                  { return n().ValType() }
func (n NativeSet) SubType() d.Typed                     { return n().TypeNat() }
func (n NativeSet) TypeName() string                     { return n().TypeName() }
func (n NativeSet) FlagType() d.Uint8Val                 { return Flag_Functional.U() }
func (n NativeSet) String() string                       { return n().String() }
func (n NativeSet) Set() SetCol                          { return NewSet(n.Pairs()...) }
func (n NativeSet) Type() TyDef {
	return Define(n().TypeName(), New(n().TypeNat()))
}
func (n NativeSet) Pairs() []Paired {
	var pairs = []Paired{}
	for _, field := range n.Fields() {
		pairs = append(
			pairs, NewPair(
				NewNative(field.Left()),
				NewNative(field.Right())))
	}
	return pairs
}

//// EXPRESSION TYPE CONSTRUCTOR
///
// TODO: make nary type safe by deriving type switch from signature and
// exposing it in a match method
//
// expression type definition takes an optional name, an expression and a
// number of expressions, or typed definitions to declare the expression
// signature. last signature expression is assumed to be the return type. all
// signature arguments before that are assumed to be the arguments types.
//
// if no signature is passed, return type is derived from expression. if no
// signature, or only return type are passed, argument types are assumed to be
// parametric matching all types.
//
// defined expressions can are enumerated and partialy applyable.
func DefinePartial(
	name string,
	expr Expression,
	signature ...Expression,
) PartialExpr {

	var arity = len(signature)
	var typed = Define(name, expr, signature...)

	// create and return nary expression
	return func(args ...Expression) Expression {
		var arglen = len(args) // count arguments
		if arglen > 0 {        // if arguments where passed
			// argument number satisfies expression arity exactly
			if arglen == arity {
				return expr.Call(args...)
			}
			// argument number undersatisfies expression arity
			if arglen < arity {
				return DefinePartial(name, NewGeneric(
					func(lateargs ...Expression) Expression {
						return expr.Call(append(
							args,
							lateargs...)...)
					},
				),
					signature[arglen:]...)
			}
			// argument number oversatisfies expressions arity
			if arglen > arity {
				var remain []Expression
				args, remain = args[:arity], args[arity:]
				var vec = NewVector(expr.Call(args...))
				for len(remain) > arity {
					args, remain = remain[:arity], remain[arity:]
					vec = vec.Append(expr.Call(args...))
				}
				if len(args) == 0 {
					return vec
				}
				return vec.Append(
					DefinePartial(
						name,
						expr,
						signature...,
					).Call(remain...))
			}
		}
		// if no arguments are passed, return definition
		return typed
	}
}

// returns the value returned when calling itself directly, passing arguments
func (n PartialExpr) Ident() Expression                  { return n }
func (n PartialExpr) Type() TyDef                        { return n().(TyDef) }
func (n PartialExpr) FlagType() d.Uint8Val               { return Flag_DataCons.U() }
func (n PartialExpr) Arity() Arity                       { return n.Type().Arity() }
func (n PartialExpr) Return() Expression                 { return n.Type().Return() }
func (n PartialExpr) Pattern() []Expression              { return n.Type().Pattern() }
func (n PartialExpr) TypeName() string                   { return n.Type().TypeName() }
func (n PartialExpr) PatternName() string                { return n.Type().PatternName() }
func (n PartialExpr) ReturnName() string                 { return n.Type().ReturnName() }
func (n PartialExpr) String() string                     { return n.Return().String() }
func (n PartialExpr) TypeFnc() TyFnc                     { return n.Return().TypeFnc() }
func (n PartialExpr) TypeNat() d.TyNat                   { return n.Return().TypeNat() }
func (n PartialExpr) Eval(args ...d.Native) d.Native     { return n.Return().Eval(args...) }
func (n PartialExpr) Call(args ...Expression) Expression { return n.Return().Call(args...) }
