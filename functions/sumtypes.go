package functions

import (
	"sort"
	"strings"

	d "github.com/joergreinhardt/gatwd/data"
)

type (
	// GENERIC EXPRESSIONS
	NoneVal      func()
	GenericConst func() Expression
	GenericFunc  func(...Expression) Expression

	//// DECLARED EXPRESSION
	FuncDef func(...Expression) Expression

	// TUPLE (TYPE[0]...TYPE[N])
	TupDef func(...Expression) TupVal
	TupVal []Expression

	// RECORD (PAIR(KEY, VAL)[0]...PAIR(KEY, VAL)[N])
	RecDef func(...Expression) RecVal
	RecVal []KeyPair
)

//// NONE VALUE CONSTRUCTOR
///
// none represens the abscence of a value of any type. implements countable,
// sliceable, consumeable, testable, compareable, key-, index- and generic pair
// interfaces to be able to stand in as return value for such expressions.
func NewNone() NoneVal { return func() {} }

func (n NoneVal) Head() Expression                   { return n }
func (n NoneVal) Tail() Consumeable                  { return n }
func (n NoneVal) Cons(...Expression) Sequential      { return n }
func (n NoneVal) Prepend(...Expression) Sequential   { return n }
func (n NoneVal) Append(...Expression) Sequential    { return n }
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
func (n NoneVal) Type() TyComp                       { return Def(None) }
func (n NoneVal) TypeElem() TyComp                   { return Def(None) }
func (n NoneVal) TypeName() string                   { return n.String() }
func (n NoneVal) Slice() []Expression                { return []Expression{} }
func (n NoneVal) Flag() d.BitFlag                    { return d.BitFlag(None) }
func (n NoneVal) FlagType() d.Uint8Val               { return Kind_Fnc.U() }
func (n NoneVal) Consume() (Expression, Consumeable) { return NewNone(), NewNone() }

//// GENERIC CONSTANT DEFINITION
///
// declares a constant value
func NewConstant(constant func() Expression) GenericConst { return constant }

func (c GenericConst) Type() TyComp                  { return Def(Constant, c().Type(), None) }
func (c GenericConst) TypeIdent() TyComp             { return c().Type().TypeIdent() }
func (c GenericConst) TypeReturn() TyComp            { return c().Type().TypeReturn() }
func (c GenericConst) TypeArguments() TyComp         { return Def(None) }
func (c GenericConst) TypeFnc() TyFnc                { return Constant }
func (c GenericConst) String() string                { return c().String() }
func (c GenericConst) Call(...Expression) Expression { return c() }

//// GENERIC FUNCTION DEFINITION
///
// declares a constant value
func NewFunction(fnc func(...Expression) Expression) GenericFunc {
	return func(args ...Expression) Expression {
		if len(args) > 0 {
			return fnc(args...)
		}
		return fnc()
	}
}

func (c GenericFunc) Call(args ...Expression) Expression {
	if len(args) > 0 {
		return c(args...)
	}
	return c()
}
func (c GenericFunc) String() string        { return c().String() }
func (c GenericFunc) TypeFnc() TyFnc        { return c().TypeFnc() }
func (c GenericFunc) Type() TyComp          { return c().Type() }
func (c GenericFunc) TypeIdent() TyComp     { return c().Type().TypeIdent() }
func (c GenericFunc) TypeReturn() TyComp    { return c().Type().TypeReturn() }
func (c GenericFunc) TypeArguments() TyComp { return c().Type().TypeArguments() }

/// PARTIAL APPLYABLE EXPRESSION VALUE
//
// defines typesafe partialy applicable expression. if the set of optional type
// argument(s) starts with a symbol, that will be assumed to be the types
// identity. otherwise the identity is derived from the passed expression,
// types first field will be the return type, its second field the (set of)
// argument type(s), additional arguments are considered propertys.
func createFuncType(expr Expression, types ...d.Typed) TyComp {
	// if type arguments have been passed, build the type based on them‥.
	if len(types) > 0 {
		// if the first element in pattern is a symbol to be used as
		// ident, just define type from type arguments‥.
		if Kind_Sym.Match(types[0].Kind()) {
			return Def(types...)
		} else { // ‥.otherwise use the expressions ident type
			return Def(append([]d.Typed{expr.Type().TypeIdent()}, types...)...)
		}
	}
	// ‥.otherwise define by expressions identity entirely in terms of the
	// passed expression type
	return Def(expr.Type().TypeIdent(),
		expr.Type().TypeReturn(),
		expr.Type().TypeArguments())

}
func Define(
	expr Expression,
	types ...d.Typed,
) FuncDef {
	var (
		ct     = createFuncType(expr, types...)
		arglen = ct.TypeArguments().Len()
	)
	// return partialy applicable function
	return func(args ...Expression) Expression {
		var length = len(args)
		if length > 0 {
			if ct.TypeArguments().MatchArgs(args...) {
				switch {
				// NUMBER OF PASSED ARGUMENTS MATCHES EXACTLY →
				case length == arglen:
					return expr.Call(args...)

				// NUMBER OF PASSED ARGUMENTS IS INSUFFICIENT →
				case length < arglen:
					// safe types of arguments remaining to be filled
					var (
						remains = ct.TypeArguments().Types()[length:]
						newpat  = Def(
							ct.TypeIdent(),
							ct.TypeReturn(),
							Def(remains...))
					)
					// define new function from remaining
					// set of argument types, enclosing the
					// current arguments & appending its
					// own aruments to them, when called.
					return Define(GenericFunc(func(lateargs ...Expression) Expression {
						// will return result, or
						// another partial, when called
						// with arguments
						if len(lateargs) > 0 {
							return expr.Call(append(
								args, lateargs...,
							)...)
						}
						// if no arguments where
						// passed, return the reduced
						// type ct
						return newpat
					}), newpat.Types()...)

				// NUMBER OF PASSED ARGUMENTS OVERSATISFYING →
				case length > arglen:
					// allocate vector to hold multiple instances
					var vector = NewVector()
					// iterate over arguments, allocate an instance per satisfying set
					for len(args) > arglen {
						vector = vector.Cons(
							expr.Call(args[:arglen]...)).(VecVal)
						args = args[arglen:]
					}
					if length > 0 { // number of leftover arguments is insufficient
						// add a partial expression as vectors last element
						vector = vector.Cons(Define(
							expr, ct.Types()...,
						).Call(args...)).(VecVal)
					}
					// return vector of instances
					return vector
				}
			}
			// passed argument(s) didn't match the expected type(s)
			return None
		}
		// no arguments where passed, return the expression type
		return ct
	}
}
func (e FuncDef) TypeFnc() TyFnc                     { return Constructor | Value }
func (e FuncDef) Type() TyComp                       { return e().(TyComp) }
func (e FuncDef) TypeIdent() TyComp                  { return e.Type().TypeIdent() }
func (e FuncDef) TypeArguments() TyComp              { return e.Type().TypeArguments() }
func (e FuncDef) TypeReturn() TyComp                 { return e.Type().TypeReturn() }
func (e FuncDef) ArgCount() int                      { return e.Type().TypeArguments().Count() }
func (e FuncDef) String() string                     { return e().String() }
func (e FuncDef) Call(args ...Expression) Expression { return e(args...) }

//// TUPLE TYPE
///
// tuple type constructor expects a slice of field types and possibly a symbol
// type flag, to define the types name, otherwise 'tuple' is the type name and
// the sequence of field types is shown instead
func NewTuple(types ...d.Typed) TupDef {
	return func(args ...Expression) TupVal {
		var tup = make(TupVal, 0, len(args))
		if Def(types...).MatchArgs(args...) {
			for _, arg := range args {
				tup = append(tup, arg)
			}
		}
		return tup
	}
}

func (t TupDef) Call(args ...Expression) Expression { return t(args...) }
func (t TupDef) TypeFnc() TyFnc                     { return Tuple | Constructor }
func (t TupDef) String() string                     { return t.Type().String() }
func (t TupDef) Type() TyComp {
	var types = make([]d.Typed, 0, len(t()))
	for _, tup := range t() {
		types = append(types, tup.Type())
	}
	return Def(Tuple, Def(types...))
}

/// TUPLE VALUE
// tuple value is a slice of expressions, constructed by a tuple type
// constructor validated according to its type pattern.
func (t TupVal) Len() int { return len(t) }
func (t TupVal) String() string {
	var strs = make([]string, 0, t.Len())
	for _, val := range t {
		strs = append(strs, val.String())
	}
	return "[" + strings.Join(strs, ", ") + "]"
}
func (t TupVal) Get(idx int) Expression {
	if idx < t.Len() {
		return t[idx]
	}
	return NewNone()
}
func (t TupVal) TypeFnc() TyFnc                     { return Tuple }
func (t TupVal) Call(args ...Expression) Expression { return NewVector(append(t, args...)...) }
func (t TupVal) Type() TyComp {
	var types = make([]d.Typed, 0, len(t))
	for _, tup := range t {
		types = append(types, tup.Type())
	}
	return Def(Tuple, Def(types...))
}

//// RECORD TYPE
///
//
func NewRecord(types ...KeyPair) RecDef {
	return func(args ...Expression) RecVal {
		var tup = make(RecVal, 0, len(args))
		if len(args) > 0 {
			for n, arg := range args {
				if len(types) > n {
					if arg.Type().Match(Key | Pair) {
						if kp, ok := arg.(KeyPair); ok {
							if strings.Compare( // equal keys
								string(kp.KeyStr()),
								string(types[n].KeyStr()),
							) == 0 && // equal values
								types[n].Value().Type().Match(
									kp.Value().Type()) {
								tup = append(tup, kp)
							}
						}
					}
				}
			}
			return tup
		}
		return types
	}
}

func (t RecDef) Call(args ...Expression) Expression { return t(args...) }
func (t RecDef) TypeFnc() TyFnc                     { return Record | Constructor }
func (t RecDef) Type() TyComp {
	var types = make([]d.Typed, 0, len(t()))
	for _, field := range t() {
		types = append(types,
			Def(
				DefSym(field.KeyStr()),
				Def(field.Value().Type()),
			))
	}
	return Def(Record, Def(types...))
}
func (t RecDef) String() string {
	var strs = make([]string, 0, len(t()))
	for _, f := range t() {
		strs = append(strs,
			"("+f.KeyStr()+" ∷ "+f.Value().String()+")",
		)
	}
	return "{" + strings.Join(strs, " ") + "}"
}

/// RECORD VALUE
// tuple value is a slice of expressions, constructed by a tuple type
// constructor validated according to its type pattern.
func (t RecVal) TypeFnc() TyFnc { return Record }
func (t RecVal) Call(args ...Expression) Expression {
	var exprs = make([]Expression, 0, len(t)+len(args))
	for _, elem := range t {
		exprs = append(exprs, elem)
	}
	for _, arg := range args {
		if arg.Type().Match(Pair | Key) {
			if kp, ok := arg.(KeyPair); ok {
				exprs = append(exprs, kp)
			}
		}
	}
	return NewVector(exprs...)
}
func (t RecVal) Type() TyComp {
	var types = make([]d.Typed, 0, len(t))
	for _, tup := range t {
		types = append(types, tup.Type())
	}
	return Def(Record, Def(types...))
}
func (t RecVal) Len() int { return len(t) }
func (t RecVal) String() string {
	var strs = make([]string, 0, t.Len())
	for _, field := range t {
		strs = append(strs,
			`"`+field.Key().String()+`"`+" ∷ "+field.Value().String())
	}
	return "{" + strings.Join(strs, " ") + "}"
}
func (t RecVal) SortByKeys() RecVal {
	var exprs = make([]Expression, 0, len(t))
	var pairs = make([]KeyPair, 0, len(t))
	for _, field := range t {
		exprs = append(exprs, field)
	}
	exprs = By(func(i, j int) bool {
		return strings.Compare(
			string(t[i].KeyStr()),
			string(t[j].KeyStr()),
		) == 0
	}).Sort(exprs)
	for _, expr := range exprs {
		pairs = append(pairs, expr.(KeyPair))
	}
	return pairs
}
func (t RecVal) SearchForKey(key string) KeyPair {
	var search = sort.Search(len(t),
		func(i int) bool { return strings.Compare(t[i].KeyStr(), key) >= 0 },
	)
}

// record sorter is a helper struct to sort record fields inline
type recordSorter struct {
	exprs []Expression
	by    By
}

func newRecordSorter(pairs []Expression, by By) *recordSorter {
	var exprs = make([]Expression, 0, len(pairs))
	for _, pair := range pairs {
		exprs = append(exprs, pair)
	}
	return &recordSorter{exprs, by}
}

func (t recordSorter) Less(i, j int) bool { return t.by(i, j) }
func (t recordSorter) Swap(i, j int)      { t.exprs[j], t.exprs[i] = t.exprs[i], t.exprs[j] }
func (t recordSorter) Len() int           { return len(t.exprs) }

// sort interface. the'By' type implements 'sort.Less() int' and is the
// function type of a parameterized sort & search function.
type By func(a, b int) bool

// sort is a method of the by function type
func (by By) Sort(exprs []Expression) []Expression {
	var sorter = newRecordSorter(exprs, by)
	sort.Sort(sorter)
	return sorter.exprs
}
