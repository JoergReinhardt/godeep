package functions

import (
	d "github.com/joergreinhardt/gatwd/data"
)

type (
	//// NONE VALUE CONSTRUCTOR
	NoneVal func()
	//// TRUTH VALUE CONSTRUCTOR
	TestExpr func(...Expression) Typed
	// OPTION VALUE CONSTRUCTOR
	OptionVal func(...Expression) Expression

	//// CASE & SWITCH TYPE CONSTRUCTORS
	CaseExpr   func(...Expression) (Expression, bool)
	CaseSwitch func(...Expression) (Expression, CaseSwitch, bool)

	//// OPTION TYPE CONSTRUCTOR
	OptionType func(...Expression) OptionVal
)

//// NONE VALUE CONSTRUCTOR
///
// none represens the abscence of a value of any type. implements consumeable,
// key-, index & generic pair interface to be returneable as such.
func NewNone() NoneVal { return func() {} }

func (n NoneVal) Ident() Expression             { return n }
func (n NoneVal) Head() Expression              { return n }
func (n NoneVal) Tail() Consumeable             { return n }
func (n NoneVal) Len() d.IntVal                 { return 0 }
func (n NoneVal) String() string                { return "⊥" }
func (n NoneVal) Eval() d.Native                { return nil }
func (n NoneVal) Call(...Expression) Expression { return nil }
func (n NoneVal) Key() Expression               { return nil }
func (n NoneVal) Index() Expression             { return nil }
func (n NoneVal) Left() Expression              { return nil }
func (n NoneVal) Right() Expression             { return nil }
func (n NoneVal) Both() Expression              { return nil }
func (n NoneVal) Value() Expression             { return nil }
func (n NoneVal) Empty() d.BoolVal              { return true }
func (n NoneVal) Flag() d.BitFlag               { return d.BitFlag(None) }
func (n NoneVal) TypeFnc() TyFnc                { return None }
func (n NoneVal) TypeNat() d.TyNat              { return d.Nil }
func (n NoneVal) TypeElem() Typed               { return None.Type() }
func (n NoneVal) TypeName() string              { return n.String() }
func (n NoneVal) FlagType() d.Uint8Val          { return Flag_Functional.U() }
func (n NoneVal) Type() TyDef                   { return Define(n.TypeName(), None) }
func (n NoneVal) Consume() (Expression, Consumeable) {
	return NewNone(), NewNone()
}

//// TRUTH VALUE CONSTRUCTOR
func NewTestTruth(name string, test func(...Expression) d.BoolVal, paratypes ...Expression) TestExpr {

	if name == "" {
		name = "Truth"
	}
	var params = make([]Typed, 0, len(paratypes))
	if len(paratypes) == 0 {
		paratypes = append(paratypes, Type)
	} else {
		for _, param := range paratypes {
			params = append(params, param.Type())
		}
	}

	return func(args ...Expression) Typed {
		if len(args) > 0 {
			if test(args...) {
				return True

			}
			return False
		}
		return Define(name, Truth, params...)
	}
}

func NewTestTrinary(name string, test func(...Expression) int, paratypes ...Expression) TestExpr {

	if name == "" {
		name = "Trinary"
	}
	var params = make([]Typed, 0, len(paratypes))
	if len(paratypes) == 0 {
		paratypes = append(paratypes, Type)
	} else {
		for _, param := range paratypes {
			params = append(params, param.Type())
		}
	}

	return func(args ...Expression) Typed {
		if len(args) > 0 {
			if test(args...) > 0 {
				return True
			}
			if test(args...) < 0 {
				return False
			}
			return Undecided
		}
		return Define(name, Trinary, params...)
	}
}

func NewTestCompare(name string, test func(...Expression) d.IntVal, paratypes ...Expression) TestExpr {

	if name == "" {
		name = "Compare"
	}
	var params = make([]Typed, 0, len(paratypes))
	if len(paratypes) == 0 {
		paratypes = append(paratypes, Type)
	}
	for _, param := range paratypes {
		params = append(params, param.Type())
	}

	return func(args ...Expression) Typed {
		if len(args) > 0 {
			if test(args...) > 0 {
				return Greater
			}
			if test(args...) < 0 {
				return Lesser
			}
			return Equal
		}
		return Define(name, Compare, params...)
	}
}
func (t TestExpr) Type() TyDef    { return t().(TyDef) }
func (t TestExpr) TypeFnc() TyFnc { return t.Type().Return().(TyFnc) }
func (t TestExpr) TypeNat() d.TyNat {
	if t.TypeFnc() == Compare {
		return d.Int
	}
	return d.Bool
}

func (t TestExpr) TypeName() string {
	return t.Type().Name() + " → " + t.Type().Return().TypeName()
}
func (t TestExpr) String() string       { return t.TypeName() }
func (t TestExpr) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (t TestExpr) Call(args ...Expression) Expression {
	if t.TypeFnc() == Compare {
		return NewData(d.IntVal(t.Compare(args...)))
	}
	return NewData(d.BoolVal(t.Test(args...)))
}

func (t TestExpr) Eval() d.Native { return t }

func (t TestExpr) Test(args ...Expression) d.BoolVal {
	if t.TypeFnc() == Compare {
		if t(args...) == Lesser || t(args...) == Greater {
			return false
		} else {
			return true
		}
	}
	if t.TypeFnc() == Trinary {
		if t(args...) == False || t(args...) == Undecided {
			return false
		} else {
			return true
		}
	}
	if t(args...) != True {
		return false
	}
	return true
}

func (t TestExpr) Compare(args ...Expression) d.IntVal {
	if t.TypeFnc() == Compare {
		switch t(args...) {
		case Lesser:
			return -1
		case Equal:
			return 0
		case Greater:
			return 1
		}
	}
	if t.TypeFnc() == Trinary {
		switch t(args...) {
		case False:
			return -1
		case Undecided:
			return 0
		case True:
			return 1
		}
	}
	if t(args...) != True {
		return -1
	}
	return 0
}

func (t TestExpr) True(arg Expression) d.BoolVal {
	if t.TypeFnc() == Truth || t.TypeFnc() == Trinary {
		if t(arg) == True {
			return true
		}
	}
	return false
}

func (t TestExpr) False(arg Expression) d.BoolVal {
	if t.TypeFnc() == Truth || t.TypeFnc() == Trinary {
		if t(arg) == False {
			return true
		}
	}
	return false
}

func (t TestExpr) Undecided(arg Expression) d.BoolVal {
	if t.TypeFnc() == Trinary {
		if t(arg) == Undecided {
			return true
		}
	}
	return false
}

func (t TestExpr) Lesser(arg Expression) d.BoolVal {
	if t.TypeFnc() == Compare {
		if t(arg) == Lesser {
			return true
		}
	}
	return false
}

func (t TestExpr) Greater(arg Expression) d.BoolVal {
	if t.TypeFnc() == Compare {
		if t(arg) == Greater {
			return true
		}
	}
	return false
}

func (t TestExpr) Equal(arg Expression) d.BoolVal {
	if t.TypeFnc() == Compare {
		if t(arg) == Equal {
			return true
		}
	}
	return false
}

//// CASE EXPRESSION CONSTRUCTOR
///
// takes a test expression and an expression to apply arguments to and return
// result from, if arguments applyed to the test expression returned true.
func NewCase(test TestExpr, expr Expression) CaseExpr {

	// generate expression to return arguments, when none has been passed
	if expr == nil {
		expr = NewGeneric(func(args ...Expression) Expression {
			switch len(args) {
			case 1:
				return args[0]
			case 2:
				return NewPair(
					args[0],
					args[1])
			}
			return NewVector(args...)
		}, "Return", Type)
	}

	// construct case type definition
	var pattern = expr.Type().Pattern()
	if len(pattern) == 0 {
		pattern = []Typed{Type}
	}
	var typed = Define(test.Type().Name(),
		expr.Type().Return(), pattern...)

	// construct case type name
	var ldel, rdel = "(", ")"
	var name = NewData(d.StrVal(
		ldel + expr.Type().PatternName() + " → " +
			typed.Name() + " ⇒ " +
			expr.Type().Name() + " → " +
			expr.Type().ReturnName() + rdel,
	))

	// return constructed case expression
	return func(args ...Expression) (Expression, bool) {

		if len(args) > 0 {
			if test.Test(args...) {
				return expr.Call(args...), true
			}
			return NewVector(args...), false
		}
		// return test, expression, type & name
		return NewVector(test, expr, typed, name), false
	}
}
func (s CaseExpr) Test() TestExpr {
	var vec, _ = s()
	return vec.(VecCol)()[0].(TestExpr)
}
func (s CaseExpr) Expr() Expression {
	var vec, _ = s()
	return vec.(VecCol)()[1]
}
func (s CaseExpr) Type() TyDef {
	var vec, _ = s()
	return vec.(VecCol)()[2].(TyDef)
}
func (s CaseExpr) TypeName() string {
	var vec, _ = s()
	return vec.(VecCol)()[3].String()
}
func (s CaseExpr) String() string       { return s.TypeName() }
func (s CaseExpr) TypeFnc() TyFnc       { return Case }
func (s CaseExpr) TypeNat() d.TyNat     { return d.Function }
func (s CaseExpr) FlagType() d.Uint8Val { return Flag_Functional.U() }
func (s CaseExpr) Eval() d.Native       { return s }
func (s CaseExpr) Call(args ...Expression) Expression {
	if result, ok := s(args...); ok {
		return result
	}
	return NewNone()
}

//// SWITCH CONSTRUCTOR ////
///
// type safe constructor wraps switch constructor that takes case arguments of
// the expression type, loops over case expression arguments once, reallocates
// as expresssion instances to pass on to the untyped private constructor.
func NewSwitch(cases ...CaseExpr) CaseSwitch {
	var exprs = make([]Expression, 0, len(cases))
	for _, cas := range cases {
		exprs = append(exprs, cas)
	}
	return conSwitch(exprs...)
}

// arbitrary typed switch constructor, to eliminate looping and reallocation of
// case expressions intendet to be stored as consumeable vector
func conSwitch(caseset ...Expression) CaseSwitch {

	return func(args ...Expression) (Expression, CaseSwitch, bool) {

		var cases = NewVector(caseset...)
		var head Expression
		var current CaseExpr
		var arguments VecCol
		var completed = NewVector()

		//// CALLED WITH ARGUMENTS ////
		///
		if len(args) > 0 {
			// safe passed arguments
			arguments = NewVector(args...)
			// if cases remain to be tested‥.
			if cases.Len() > 0 {
				// consume head & reassign caselist
				head, cases = cases.ConsumeVec()
				// cast head as case expression
				current = head.(CaseExpr)
				if result, ok := current(args...); ok {
					//// SUCCESSFUL CASE EVALUATION ////
					///
					// replenish case list for switch
					// instance reusal (in case switch gets
					// called empty to retrieve case list)
					cases = NewVector(caseset...)
					// return result, current case and
					// arguments that where passed.
					return result, conSwitch(caseset...), true
				}
				//// FAILED CASE EVALUATION ///
				///
				// add failed case and evaluated arguments to
				// the list of completed cases
				completed = completed.Append(
					NewPair(current, arguments))

				return completed,
					conSwitch(cases()...),
					false
			}
			//// CASES DEPLETION ///
			///
			// replenish cases for switch reusal and return
			// replenished switch instance for eventual reuse.
			cases = NewVector(caseset...)
			return nil, conSwitch(cases()...), false
		}
		//// CALLED WITH EMPTY ARGUMENT SET ////
		///
		// when called without arguments, return list of all defined
		// cases and log of cases completed so far, including arguments
		// that where passed to those cases for evaluation.
		return cases, conSwitch(caseset...), false
	}
}
func (s CaseSwitch) Cases() VecCol {
	var cases, _, _ = s()
	return cases.(VecCol)
}
func (s CaseSwitch) Type() TyDef {
	return Define(s.TypeName(), s.TypeFnc())
}
func (s CaseSwitch) TypeName() string {
	return "[T] → (Case Switch) → (T, [T]) "
}
func (s CaseSwitch) String() string       { return s.TypeName() }
func (s CaseSwitch) TypeFnc() TyFnc       { return Switch }
func (s CaseSwitch) TypeNat() d.TyNat     { return d.Function }
func (s CaseSwitch) FlagType() d.Uint8Val { return Flag_Functional.U() }

// test one set of arguments against all cases until either successful result
// is yielded, or all cases are depleted. gets called by eval & call methods.
// when called without arguments, list of all cases and list of completed
// cases, including former call arguments will be returned.
func (s CaseSwitch) TestAllCases(args ...Expression) (Expression, Expression) {
	var ok bool
	var result, caseargs Expression
	if len(args) > 0 {
		result, caseargs, ok = s(args...)
		for result != nil {
			if ok {
				return result, caseargs
			}
			result, caseargs, ok = caseargs.(CaseSwitch)(args...)
		}
		return nil, caseargs
	}
	return result, caseargs
}
func (s CaseSwitch) Call(args ...Expression) Expression {
	if len(args) > 0 {
		var result, _ = s.TestAllCases(args...)
		if result != nil {
			return result
		}
	}
	return NewNone()
}
func (s CaseSwitch) Eval() d.Native { return s }

///////////////////////////////////////////////////////////////////////////////
/// OPTION TYPE CONSTRUCTOR
func NewOptionType(test CaseSwitch, types ...Expression) OptionType {
	return func(args ...Expression) OptionVal {
		return NewOptionVal(test, NewVector(types...))
	}
}

func (o OptionType) Call(args ...Expression) Expression { return o().Call(args...) }
func (o OptionType) Expr() Expression                   { return o() }
func (o OptionType) FlagType() d.Uint8Val               { return Flag_Def.U() }
func (o OptionType) TypeFnc() TyFnc                     { return Option }
func (o OptionType) ElemType() TyDef                    { return o().Type() }
func (o OptionType) String() string                     { return o().String() }
func (o OptionType) Type() TyDef {
	return Define(o().TypeName(), o.ElemType())
}
func (o OptionType) TypeName() string {
	var name string
	return name
}

/// OPTION VALUE CONSTRUCTOR
func NewOptionVal(test CaseSwitch, exprs ...Expression) OptionVal {
	return func(args ...Expression) Expression {
		var expr, index = test.TestAllCases(args...)
		if !expr.TypeFnc().Match(None) {
			var idx = int(index.(NativeExpr)().(d.IntVal))
			var result = exprs[idx]
			if len(args) > 0 {
				return result.Call(args...)
			}
			return result
		}
		return expr
	}
}
func (o OptionVal) Call(args ...Expression) Expression { return o().Call(args...) }
func (o OptionVal) TypeFnc() TyFnc                     { return o(HigherOrder).TypeFnc() }
func (o OptionVal) FlagType() d.Uint8Val               { return Flag_Functional.U() }
func (o OptionVal) String() string                     { return o(HigherOrder).String() }
func (o OptionVal) TypeName() string                   { return o(HigherOrder).TypeName() }
func (o OptionVal) Type() TyDef                        { return Define(o().TypeName(), o()) }
