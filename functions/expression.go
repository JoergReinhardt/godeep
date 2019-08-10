package functions

import d "github.com/joergreinhardt/gatwd/data"

type (
	// GENERIC EXPRESSIONS
	NoneVal  func()
	ConstVal func() Expression
	FuncVal  func(...Expression) Expression
	ExprVal  func(...Expression) Expression
	ParamVal func(...Expression) Expression
)

//// NONE VALUE CONSTRUCTOR
///
// none represens the abscence of a value of any type. implements countable,
// sliceable, consumeable, testable, compareable, key-, index- and generic pair
// interfaces to be able to stand in as return value for such expressions.
func NewNone() NoneVal { return func() {} }

func (n NoneVal) Head() Expression                   { return n }
func (n NoneVal) Tail() Consumeable                  { return n }
func (n NoneVal) Append(...Expression) Consumeable   { return n }
func (n NoneVal) Len() int                           { return 0 }
func (n NoneVal) Compare(...Expression) int          { return -1 }
func (n NoneVal) String() string                     { return "⊥" }
func (n NoneVal) Call(...Expression) Expression      { return nil }
func (n NoneVal) Key() Expression                    { return nil }
func (n NoneVal) Index() Expression                  { return nil }
func (n NoneVal) Left() Expression                   { return nil }
func (n NoneVal) Right() Expression                  { return nil }
func (n NoneVal) Both() Expression                   { return nil }
func (n NoneVal) Value() Expression                  { return nil }
func (n NoneVal) Empty() d.BoolVal                   { return true }
func (n NoneVal) Test(...Expression) bool            { return false }
func (n NoneVal) TypeFnc() TyFnc                     { return None }
func (n NoneVal) TypeNat() d.TyNat                   { return d.Nil }
func (n NoneVal) Type() TyPattern                    { return Def(None) }
func (n NoneVal) TypeElem() TyPattern                { return Def(None) }
func (n NoneVal) TypeName() string                   { return n.String() }
func (n NoneVal) Slice() []Expression                { return []Expression{} }
func (n NoneVal) Flag() d.BitFlag                    { return d.BitFlag(None) }
func (n NoneVal) FlagType() d.Uint8Val               { return Flag_Function.U() }
func (n NoneVal) Consume() (Expression, Consumeable) { return NewNone(), NewNone() }

//// CONSTANT DECLARATION
///
// declares a constant value
func DecConstant(constant func() Expression) ConstVal { return constant }

func (c ConstVal) Type() TyPattern {
	return Def(None, Def(Constant, c().Type().TypeIdent()), c().Type())
}
func (c ConstVal) TypeFnc() TyFnc                { return Constant }
func (c ConstVal) String() string                { return c().String() }
func (c ConstVal) Call(...Expression) Expression { return c() }

//// EXPRESSION DECLARATION
///
// declares an expression with defined argument-, return- and an optional identity type
func DecFuntion(
	fn func(...Expression) Expression,
	argtype, retype d.Typed,
	identypes ...d.Typed,
) FuncVal {

	var (
		ident TyPattern
	)

	if len(identypes) == 0 {
		ident = Def(Value)
	} else {
		ident = Def(identypes...)
	}

	var pattern = Def(argtype, ident, retype)

	return func(args ...Expression) Expression {
		if len(args) > 0 {
			if pattern.TypeArguments().MatchArgs(args...) {
				return fn(args...)
			}
		}
		return pattern
	}
}

func (g FuncVal) TypeFnc() TyFnc                     { return Value }
func (g FuncVal) Type() TyPattern                    { return g().(TyPattern) }
func (g FuncVal) String() string                     { return g.Type().TypeName() }
func (g FuncVal) Call(args ...Expression) Expression { return g(args...) }

/// PARTIAL APPLYABLE EXPRESSION VALUE
//
// element values yield a subelements of optional, tuple, or enumerable
// expressions with sub-type pattern as second return value
func DecExpression(
	expr Expression,
	argtype, retype d.Typed,
	identypes ...d.Typed,
) ExprVal {

	var (
		arglen         int
		ident, pattern TyPattern
	)
	if len(identypes) == 0 {
		ident = Def(expr.TypeFnc())
	} else {
		ident = Def(identypes...)
	}
	if Flag_Pattern.Match(argtype.FlagType()) {
		arglen = len(argtype.(TyPattern).Patterns())
	} else {
		arglen = 1
	}
	pattern = Def(argtype, ident, retype)

	return func(args ...Expression) Expression {
		var length = len(args)
		if length > 0 {
			if pattern.TypeArguments().MatchArgs(args...) {
				var result Expression
				switch {
				case length == arglen:
					result = expr.Call(args...)
					return result

				case length < arglen:
					var argtypes = make(
						[]d.Typed, 0,
						len(pattern.TypeArguments()[length:]),
					)
					for _, atype := range pattern.TypeArguments()[length:] {
						argtypes = append(argtypes, atype)
					}
					var pattern = Def(Def(argtypes...), ident, retype)
					return DecExpression(FuncVal(
						func(lateargs ...Expression) Expression {
							if len(lateargs) > 0 {
								return expr.Call(append(
									args, lateargs...,
								)...)
							}
							return pattern
						}), Def(argtypes...), ident, retype)

				case length > arglen:
					var vector = NewVector()
					for len(args) > arglen {
						vector = vector.Con(
							expr.Call(args[:arglen]...))
						args = args[arglen:]
					}
					if length > 0 {
						vector = vector.Con(DecExpression(
							expr, argtype, retype, identypes...,
						).Call(args...))
					}
					return vector
				}
			}
			return NewNone()
		}
		return pattern
	}
}
func (e ExprVal) TypeFnc() TyFnc                     { return Value }
func (e ExprVal) Type() TyPattern                    { return e().(TyPattern) }
func (e ExprVal) String() string                     { return e().String() }
func (e ExprVal) Call(args ...Expression) Expression { return e(args...) }

//// PARAMETRIC VALUE
///
// declare parametric value from set of cases. when arguments are passed, they
// are applyed to all enclosed cases in the order they where passed in at
// creation. first result that is not a none instance will be returned, returns
// type pattern, when called empty.
func DecParametric(params ...ExprVal) ParamVal {
	var (
		current Expression
		pattern TyPattern
		exprs   = make([]Expression, 0, len(params))
	)
	for _, param := range params {
		exprs = append(exprs, param)
	}

	return func(args ...Expression) Expression {
		if len(args) > 0 {
			if len(exprs) > 0 {
				current = exprs[0]
				if len(exprs) > 1 {
					exprs = exprs[1:]
				} else {
					exprs = nil
				}
				return current.Call(args...)
			}
			for _, param := range params {
				exprs = append(exprs, param)
			}
			return nil
		}
		return pattern
	}
}

func (p ParamVal) Call(args ...Expression) Expression {
	var result = p(args...)
	return result
}
func (p ParamVal) TypeFnc() TyFnc  { return Parametric }
func (p ParamVal) Type() TyPattern { return p().(TyPattern) }
func (p ParamVal) String() string  { return p().String() }