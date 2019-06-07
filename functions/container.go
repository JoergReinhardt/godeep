/*
  FUNCTIONAL CONTAINERS

  containers implement enumeration of functional types, aka lists, vectors sets, pairs, tuples‥.
*/
package functions

import (
	d "github.com/joergreinhardt/gatwd/data"
)

type (
	//// PREDICATE
	PredictArg  func(Callable) bool
	PredictAll  func(...Callable) bool
	PredictAny  func(...Callable) bool
	PredictExpr func(...Callable) bool

	//// CASE-EXPRESSION | CASE-SWITCH
	CaseSwitch func(...Callable) (Callable, bool)

	////  NONE | JUST | MAYBE
	JustVal func(...Callable) Callable
	NoneVal func()

	//// STATIC EXPRESSIONS
	ConstantExpr func() Callable
	UnaryExpr    func(Callable) Callable
	BinaryExpr   func(a, b Callable) Callable
	NaryExpr     func(...Callable) Callable
	VariadicExpr func(...Callable) Callable

	//// DATA VALUE
	Data func() d.Native
)

//// PREDICATE
///
// predict one is an expression that returns either true, or false depending on
// first passed arguement passed. succeeding arguements are ignored
func NewPredictArg(pred func(Callable) bool) PredictArg {
	return func(arg Callable) bool { return pred(arg) }
}
func (p PredictArg) Eval() d.Native { return d.BoolVal(false) }
func (p PredictArg) Call(args ...Callable) Callable {
	if len(args) > 0 {
		return NewData(d.BoolVal(p(args[0])))
	}
	return NewData(d.BoolVal(false))
}

// return single argument predicate as multi argument predicate
func (p PredictArg) ToPredictNarg() PredictExpr {
	return func(args ...Callable) bool {
		if len(args) > 0 {
			return p(args[0])
		}
		return false
	}
}
func (p PredictArg) String() string   { return "Predicate" }
func (p PredictArg) TypeNat() d.TyNat { return d.Functor }
func (p PredictArg) TypeFnc() TyFnc   { return Predicate }

///////////////////////////////////////////////////////////////////////////////
// predict many returns true, or false depending on all arguments that have
// been passed calling it
func NewPredictNarg(pred func(...Callable) bool) PredictExpr {
	return func(args ...Callable) bool {
		return pred(args...)
	}
}
func (p PredictExpr) Call(args ...Callable) Callable {
	if len(args) > 0 {
		return NewData(d.BoolVal(p(args...)))
	}
	return NewData(d.BoolVal(false))
}
func (p PredictExpr) Eval() d.Native   { return d.BoolVal(false) }
func (p PredictExpr) String() string   { return "Nary Predicate" }
func (p PredictExpr) TypeNat() d.TyNat { return d.Functor }
func (p PredictExpr) TypeFnc() TyFnc   { return Predicate }

///////////////////////////////////////////////////////////////////////////////
// all-predicate returns true, if all arguments passed yield true, when applyed
// to predicate one after another
func NewPredictAll(pred func(Callable) bool) PredictExpr {
	return func(args ...Callable) bool {
		var result = true
		for _, arg := range args {
			if !pred(arg) {
				return false
			}
		}
		return result
	}
}

// call passes arguments on to the enclosed all-predicate
func (p PredictAll) Call(args ...Callable) Callable {
	if len(args) > 0 {
		return NewData(d.BoolVal(p(args...)))
	}
	return NewData(d.BoolVal(false))
}

func (p PredictAll) Eval() d.Native { return d.BoolVal(false) }

func (p PredictAll) String() string   { return "All Predicate" }
func (p PredictAll) TypeFnc() TyFnc   { return Predicate }
func (p PredictAll) TypeNat() d.TyNat { return d.Functor }

///////////////////////////////////////////////////////////////////////////////
// will return true, if any of the passed arguments yield true, when applyed to
// predicate one after another
func NewPredictAny(pred func(Callable) bool) PredictExpr {
	return func(args ...Callable) bool {
		var result = false
		for _, arg := range args {
			if pred(arg) {
				return true
			}
		}
		return result
	}
}
func (p PredictAny) Call(args ...Callable) Callable {
	if len(args) > 0 {
		return NewData(d.BoolVal(p(args...)))
	}
	return NewData(d.BoolVal(false))
}
func (p PredictAny) Eval() d.Native   { return d.BoolVal(false) }
func (p PredictAny) String() string   { return "Any Predicate" }
func (p PredictAny) TypeNat() d.TyNat { return d.Functor }
func (p PredictAny) TypeFnc() TyFnc   { return Predicate }

//// CASE SWITCH
///
// applys passed arguments to all enclosed cases in the order passed to the
// switch constructor
func NewSwitch(predicates ...PredictExpr) CaseSwitch {

	// cast predicates as slice of callables
	var exprs = []Callable{}
	for _, predicate := range predicates {
		exprs = append(exprs, predicate)
	}
	// create consumeable list of cases
	var cases = NewList(exprs...)
	// allocate current case
	var current Callable

	return func(args ...Callable) (Callable, bool) {
		// call consumeable to yield current case
		current, cases = cases()
		// if call yielded a case
		if current != nil {
			// scrutinize arguments by applying the case
			if current.(PredictExpr)(args...) {
				// replenish cases before returning
				// successfully scrutinized arguments
				cases = NewList(exprs...)
				// if a single argument has been passed, it is
				// returned without wrapping vector
				if len(args) == 1 {
					return args[0], true
				}
				// return set of arguments and true
				return NewVector(args...), true
			}
			// don't replenish cases, since switch has only been
			// parially applyed yet, return set of arguments and
			// false
			if len(args) == 1 {
				return args[0], false
			}
			return NewVector(args...), false
		}
		// replenish cases before final return due to case depletion
		cases = NewList(exprs...)
		// final call returns instance of none and indicates false
		return NewNone(), false
	}
}

// switches call method iterates over all cases until either boolean indicates
// true and resulting arguments are returned, or all cases are depleted
func (s CaseSwitch) Call(args ...Callable) Callable {
	// if arguments have been passed
	if len(args) > 0 {
		// call switch instance passing the arguments
		var vec, ok = s(args...)
		// if call did not yield an instance of none
		for !vec.TypeFnc().Match(None) {
			// if boolean indicates success
			if ok {
				// return set of arguments
				return vec
			}
			// otherwise call switch to scrutinize next case
			vec, ok = s(args...)
		}
	}
	// return none if all cases are scrutinized, or no arguments where
	// passed
	return NewNone()
}

// eval converts its arguments to callable and evaluates the result to yield a
// return value of native type
func (s CaseSwitch) Eval() d.Native { return d.NilVal{} }

func (s CaseSwitch) String() string   { return "Switch" }
func (s CaseSwitch) TypeFnc() TyFnc   { return Switch }
func (s CaseSwitch) TypeNat() d.TyNat { return d.Functor }

///////////////////////////////////////////////////////////////////////////////
//// MAYVE → JUST | NONE
///
// maybe type has a type constructor to define variations of the maybe type. it
// expects a case switch and an instance of callable. it returns a maybe type
// data constructors enclosing the arguments passed to it.
//
// each data constructor expects arguments to pass on to the expression,
// scrutinizing the result by applying it to the switch case, to either return
// the resulting value wrapped in a just instance, or an instance of none, if
// the result could not be scrutinize
func DefineMaybeType(swi CaseSwitch, expr Callable) MaybeCons {

	return func(args ...Callable) Callable {
		var result = swi.Call(expr.Call(args...))
		if result.TypeFnc().Match(None) {
			return result
		}
		return consJust(result)
	}
}

type MaybeCons func(...Callable) Callable

//// JUST VALUE
///
// just constructor is not exported and called by maybe data constructor
func consJust(expr Callable) JustVal {
	return func(args ...Callable) Callable {
		if len(args) > 0 {
			return expr.Call(args...)
		}
		return expr.Call()
	}
}

// prepend just and concatenate type names of functional, and native type
// divided by dots
func (n JustVal) TypeName() string {
	return "Just·" +
		n().TypeFnc().TypeName() +
		"·" +
		n().TypeNat().TypeName()
}
func (n JustVal) Ident() Callable                { return n }
func (n JustVal) Value() Callable                { return n() }
func (n JustVal) Call(args ...Callable) Callable { return n().Call(args...) }
func (n JustVal) Eval() d.Native                 { return n().Eval() }
func (n JustVal) String() string                 { return n().String() }
func (n JustVal) TypeNat() d.TyNat               { return n().TypeNat() }
func (n JustVal) TypeFnc() TyFnc                 { return Just | n().TypeFnc() }
func (n JustVal) TypeComp() TyComp {
	return DefineComposedType(
		"Just·"+n().TypeFnc().TypeName()+"·"+n().TypeNat().TypeName(),
		n().TypeFnc(), n.TypeNat(),
	)
}

//// NONE VALUE
func NewNone() NoneVal { return func() {} }

func (n NoneVal) Ident() Callable           { return n }
func (n NoneVal) Len() int                  { return 0 }
func (n NoneVal) String() string            { return "⊥" }
func (n NoneVal) Eval() d.Native            { return nil }
func (n NoneVal) Value() Callable           { return nil }
func (n NoneVal) Call(...Callable) Callable { return nil }
func (n NoneVal) Empty() bool               { return true }
func (n NoneVal) TypeFnc() TyFnc            { return None }
func (n NoneVal) TypeNat() d.TyNat          { return d.Nil }
func (n NoneVal) TypeName() string          { return n.String() }

// none implements consumeable interface
func (n NoneVal) Head() Callable                   { return NewNone() }
func (n NoneVal) Tail() Consumeable                { return NewNone() }
func (n NoneVal) Consume() (Callable, Consumeable) { return NewNone(), NewNone() }

///////////////////////////////////////////////////////////////////////////////
//// STATIC FUNCTION EXPRESSIONS OF PREDETERMINED ARITY
///
// static arity functions ignore abundant arguments and return none/nil when
// arguments are passed to call, or eval method
//
/// CONSTANT EXPRESSION
//
// does not take any arguments and ignores all arguments passed
func NewConstant(
	expr Callable,
) ConstantExpr {
	return func() Callable { return expr }
}
func (c ConstantExpr) Ident() Callable  { return c() }
func (c ConstantExpr) Arity() Arity     { return Arity(0) }
func (c ConstantExpr) TypeFnc() TyFnc   { return c().TypeFnc() }
func (c ConstantExpr) TypeNat() d.TyNat { return c().TypeNat() }
func (c ConstantExpr) Eval() d.Native {
	return c().Eval()
}
func (c ConstantExpr) Call(args ...Callable) Callable {
	if len(args) > 0 {
		return NewNone()
	}
	return c()
}

/// UNARY EXPRESSION
//
// expects one argument, ignores abundant arguments and returns nil/none, when
// no arguments are passed.
func NewUnary(
	expr Callable,
) UnaryExpr {
	return func(arg Callable) Callable { return expr.Call(arg) }
}
func (u UnaryExpr) Ident() Callable  { return u }
func (u UnaryExpr) Arity() Arity     { return Arity(1) }
func (u UnaryExpr) TypeFnc() TyFnc   { return u(NewNone()).TypeFnc() }
func (u UnaryExpr) TypeNat() d.TyNat { return d.Functor.TypeNat() }
func (u UnaryExpr) Call(args ...Callable) Callable {
	if len(args) == 1 {
		return u(args[0]).Call()
	}
	return NewNone()
}
func (u UnaryExpr) Eval() d.Native {
	return d.NilVal{}
}

/// BINARY EXPRESSION
//
// expects two arguments, ignores abundant arguments returns nil/none, when
// less than two arguments are passed.
func NewBinary(
	expr Callable,
) BinaryExpr {
	return func(a, b Callable) Callable {
		return expr.Call(a, b)
	}
}

func (b BinaryExpr) Ident() Callable  { return b }
func (b BinaryExpr) Arity() Arity     { return Arity(2) }
func (b BinaryExpr) TypeFnc() TyFnc   { return b(NewNone(), NewNone()).TypeFnc() }
func (b BinaryExpr) TypeNat() d.TyNat { return d.Functor.TypeNat() }
func (b BinaryExpr) Call(args ...Callable) Callable {
	if len(args) == 2 {
		return b(args[0], args[1])
	}
	return NewNone()
}
func (b BinaryExpr) Eval() d.Native {
	return d.NilVal{}
}

/// VARIADIC EXPRESSION
//
// variadic expression has an unknown arity and can take an arbitrary number of
// arguments passed when calling it
func NewVariadic(
	expr Callable,
) VariadicExpr {
	return func(args ...Callable) Callable {
		return expr.Call(args...)
	}
}
func (n VariadicExpr) Ident() Callable             { return n }
func (n VariadicExpr) TypeFnc() TyFnc              { return n().TypeFnc() }
func (n VariadicExpr) TypeNat() d.TyNat            { return d.Functor.TypeNat() }
func (n VariadicExpr) Call(d ...Callable) Callable { return n(d...) }
func (n VariadicExpr) Eval() d.Native              { return n().Eval() }

/// NARY EXPRESSION
//
// nary expression knows it's composed type & arity and returns an expression
// by applying arguments to the enclosed expression, handling partial-, exact-,
// and oversatisfied calls, by returning either:
//
// - a partialy applied function and arity reduced by the number of arguments
//   passed allready,
//
// - the result of applying the exact number of arguments to the expression
//
// - a slice of results of applying abundant arguments repeatedly according to
//   arity until argument depletion. last result may be a partialy applyed nary
func NewNary(
	expr Callable,
	comp TyComp,
	arity int,
) NaryExpr {
	return func(args ...Callable) Callable {
		if len(args) > 0 {
			// argument number satify expression arity exactly
			if len(args) == arity {
				// expression expects one or more arguments
				if arity > 0 {
					// return fully applyed expression with
					// remaining arity set to be zero
					return expr.Call(args...)
				}
			}

			// argument number undersatisfies expression arity
			if len(args) < arity {
				// return a parially applyed expression with
				// reduced arity wrapping a variadic expression
				// that can take succeeding arguments to
				// concatenate to arguments passed in  prior
				// calls.
				return NewNary(VariadicExpr(
					func(succs ...Callable) Callable {
						// return result of calling the
						// nary, passing arguments
						// concatenated to those passed
						// in preceeding calls
						return NewNary(expr, comp, arity).Call(
							append(args, succs...)...)
					}), comp, arity-len(args))
			}

			// argument number oversatisfies expressions arity
			if len(args) > arity {
				// allocate slice of results
				var results = []Callable{}

				// iterate aver arguments & create fully satisfied
				// expressions, wile argument number is higher than
				// expression arity
				for len(args) > arity {
					// apped result of fully satisfiedying
					// expression to results slice
					results = append(
						results,
						NewNary(expr, comp, arity)(
							args[:arity]...),
					)
					// reassign remaining arguments
					args = args[arity:]
				}

				// if any arguments remain, append result of
				// partial application
				if len(args) <= arity && arity > 0 {
					results = append(
						results,
						NewNary(expr, comp, arity)(
							args...))
				}
				// return results slice
				return NewVector(results...)
			}
		}
		// no arguments passed with the call, return arity and composed
		// type instead
		return NewPair(Arity(arity), comp)
	}
}

// eval method converts it's arguments to implement callable to pass on to the
// call method and returns the result as native by evaluating it
func (n NaryExpr) Eval() d.Native { return n().Eval() }

// returns the value returned when calling itself directly, passing along any
// given argument.
func (n NaryExpr) Call(args ...Callable) Callable {
	if len(args) > 0 {
		return n(args...)
	}
	return n()
}
func (n NaryExpr) Arity() Arity     { val := n(); return val.(Paired).Left().(Arity) }
func (n NaryExpr) CompType() TyComp { val := n(); return val.(Paired).Right().(TyComp) }
func (n NaryExpr) TypeName() string { return n.CompType().TypeName() }
func (n NaryExpr) TypeFnc() TyFnc   { return n.CompType().TypeFnc() }
func (n NaryExpr) TypeNat() d.TyNat { return n.CompType().TypeNat() }
func (n NaryExpr) Ident() Callable  { return n }

//// DATA VALUE
///
// data value implements the callable interface but returns an instance of
// data/Value. the eval method of every native can be passed as argument
// instead of the value itself, as in 'DataVal(native.Eval)', to delay, or even
// possibly ommit evaluation of the underlying data value for cases where
// lazynes is paramount.
func New(inf ...interface{}) Callable { return NewData(d.New(inf...)) }

func NewData(args ...d.Native) Data {
	// if any initial arguments have been passed
	if len(args) > 0 {
		return func() d.Native { return d.NewFromData(args...) }
	}
	return func() d.Native { return d.NewNil() }
}

func (n Data) Call(args ...Callable) Callable { return n }
func (n Data) Eval() d.Native                 { return n() }
func (n Data) TypeFnc() TyFnc                 { return Native }
func (n Data) TypeNat() d.TyNat               { return n().TypeNat() }
func (n Data) String() string                 { return n().String() }
func (n Data) TypeName() string               { return n().TypeNat().TypeName() }

//// HELPER FUNCTIONS TO HANDLE ARGUMENTS
///
// since every callable also needs to implement the eval interface and data as
// such allways boils down to native values, conversion between callable-/ &
// native arguments is frequently needed. arguments may also need to be
// reversed when intendet to be passed to certain recursive expressions, or
// returned by those
//
/// REVERSE ARGUMENTS
func revArgs(args ...Callable) []Callable {
	var rev = []Callable{}
	for i := len(args) - 1; i > 0; i-- {
		rev = append(rev, args[i])
	}
	return rev
}

/// CONVERT NATIVE TO FUNCTIONAL
func natToFnc(args ...d.Native) []Callable {
	var result = []Callable{}
	for _, arg := range args {
		result = append(result, NewData(arg))
	}
	return result
}

/// CONVERT FUNCTIONAL TO NATIVE
func fncToNat(args ...Callable) []d.Native {
	var result = []d.Native{}
	for _, arg := range args {
		result = append(result, arg.Eval())
	}
	return result
}

/// GROUP ARGUMENTS PAIRWISE
//
// assumes the arguments to either implement paired, or be alternating pairs of
// key & value. in case the number of passed arguments that are not pairs is
// uneven, last field will be filled up with a value of type none
func argsToPaired(args ...Callable) []Paired {
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
