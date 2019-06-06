package functions

import (
	d "github.com/joergreinhardt/gatwd/data"
	"github.com/joergreinhardt/gatwd/lex"
)

type (
	TyComp    func() (string, []Typed)
	TyFnc     d.BitFlag
	Arity     d.Int8Val
	Propertys d.Uint8Val
)

//go:generate stringer -type=TyFnc
const (
	/// KIND FLAGS ///
	Type TyFnc = 1 << iota
	Native
	Atom
	Key
	Index
	/// EXPRESSION CALL PROPERTYS
	CallArity
	CallPropertys
	/// TYPE CLASSES
	Numbers
	Strings
	Bytes
	/// COLLECTION TYPES
	Element
	List
	Vector
	Tuple
	Record
	Enum
	Set
	Pair
	/// FUNCTORS AND MONADS
	Constructor
	Functor
	Applicable
	Monad
	/// MONADIC SUB TYPES
	Undecided
	Predicate
	False
	True
	Equal
	Lesser
	Greater
	Just
	None
	Case
	Switch
	Either
	Or
	If
	Else
	Do
	While
	/// IO
	Buffer
	Reader
	Writer
	/// HIGHER ORDER TYPE
	HigherOrder

	Collections = List | Vector | Tuple | Record | Enum |
		Set | Pair

	Options = Undecided | False | True | Equal | Lesser |
		Greater | Just | None | Case | Switch | Either |
		Or | If | Else | Do | While

	Parameters = CallPropertys | CallArity

	Kinds = Type | Native | Atom | Functor

	Truth = Undecided | False | True

	Ordered = Equal | Lesser | Greater

	Maybe = Just | None

	Alternatives = Either | Or

	Branch = If | Else

	Continue = Do | While

	IO = Buffer | Reader | Writer

	Consumeables = Collections | Applicable | Monad | IO
)

///////////////////////////////////////////////////////////////////////////////
//go:generate stringer -type Propertys
const (
	Default Propertys = 0
	PostFix Propertys = 1
	InFix   Propertys = 1 + iota
	// ⌐: PreFix
	Atomic
	// ⌐: Thunk
	Eager
	// ⌐: Lazy
	RightBound
	// ⌐: Left_Bound
	Mutable
	// ⌐: Imutable
	SideEffect
	// ⌐: Pure
	Primitive
	// ⌐: Parametric
)

func FlagToProp(flag d.BitFlag) Propertys          { return Propertys(uint8(flag.Uint())) }
func (p Propertys) PostFix() bool                  { return p.Flag().Match(PostFix.Flag()) }
func (p Propertys) InFix() bool                    { return !p.Flag().Match(PostFix.Flag()) }
func (p Propertys) Atomic() bool                   { return p.Flag().Match(Atomic.Flag()) }
func (p Propertys) Thunk() bool                    { return !p.Flag().Match(Atomic.Flag()) }
func (p Propertys) Eager() bool                    { return p.Flag().Match(Eager.Flag()) }
func (p Propertys) Lazy() bool                     { return !p.Flag().Match(Eager.Flag()) }
func (p Propertys) RightBound() bool               { return p.Flag().Match(RightBound.Flag()) }
func (p Propertys) LeftBound() bool                { return !p.Flag().Match(RightBound.Flag()) }
func (p Propertys) Mutable() bool                  { return p.Flag().Match(Mutable.Flag()) }
func (p Propertys) Imutable() bool                 { return !p.Flag().Match(Mutable.Flag()) }
func (p Propertys) SideEffect() bool               { return p.Flag().Match(SideEffect.Flag()) }
func (p Propertys) Pure() bool                     { return !p.Flag().Match(SideEffect.Flag()) }
func (p Propertys) Primitive() bool                { return p.Flag().Match(Primitive.Flag()) }
func (p Propertys) Parametric() bool               { return !p.Flag().Match(Primitive.Flag()) }
func (p Propertys) TypeNat() d.TyNat               { return d.Flag }
func (p Propertys) TypeFnc() TyFnc                 { return HigherOrder }
func (p Propertys) Flag() d.BitFlag                { return d.BitFlag(uint64(p)) }
func (p Propertys) Eval(a ...d.Native) d.Native    { return p }
func (p Propertys) Call(args ...Callable) Callable { return p }
func (p Propertys) Match(flag d.BitFlag) bool      { return p.Flag().Match(flag) }
func (p Propertys) MatchProperty(arg Propertys) bool {
	if p&arg != 0 {
		return true
	}
	return false
}
func (p Propertys) Print() string {

	var flags = p.Flag().Decompose()
	var str string
	var l = len(flags)

	if l > 1 {
		for i, typed := range flags {

			if typed.FlagType() == 1 {

				str = str + typed.(d.TyNat).String()
			}

			if typed.FlagType() == 2 {

				str = str + typed.(TyFnc).String()
			}

			if typed.FlagType() == 3 {

				str = str + typed.(lex.TySyntax).String()
			}

			if i < l-1 {
				str = str + " "
			}
		}
	}

	return p.String()
}

//go:generate stringer -type Arity
const (
	Nullary Arity = 0 + iota
	Unary
	Binary
	Ternary
	Quaternary
	Quinary
	Senary
	Septenary
	Octonary
	Novenary
	Denary
)

func (a Arity) Eval(...d.Native) d.Native { return d.Int8Val(a) }
func (a Arity) Call(...Callable) Callable { return NewAtom(a.Eval()) }
func (a Arity) Int() int                  { return int(a) }
func (a Arity) Flag() d.BitFlag           { return d.BitFlag(a) }
func (a Arity) TypeNat() d.TyNat          { return d.Flag }
func (a Arity) TypeFnc() TyFnc            { return HigherOrder }
func (a Arity) Match(arg Arity) bool      { return a == arg }

// type TyFnc d.BitFlag
// encodes the kind of functional data as bitflag
func (t TyFnc) FlagType() uint8 { return 2 }
func (t TyFnc) TypeName() string {
	var count = t.Flag().Count()
	if count > 1 {
		var str string
		for i, flag := range t.Flag().Decompose() {
			str = str + TyFnc(flag.Flag()).String()
			if i < count-1 {
				str = str + "·"
			}
		}
		return str
	}
	return t.String()
}
func (t TyFnc) TypeFnc() TyFnc                 { return Type }
func (t TyFnc) TypeNat() d.TyNat               { return d.Flag }
func (t TyFnc) Call(args ...Callable) Callable { return t.TypeFnc() }
func (t TyFnc) Eval(args ...d.Native) d.Native { return t.TypeNat() }
func (t TyFnc) Flag() d.BitFlag                { return d.BitFlag(t) }
func (t TyFnc) Match(arg d.Typed) bool         { return t.Flag().Match(arg) }
func (t TyFnc) Uint() uint                     { return d.BitFlag(t).Uint() }

//// COMPOSED TYPE
///
// composition type to mark higher order types. it returns a type name and a
// slice of typed interface implementing instances.  so recursively nested
// composed types of arbitrary depth can be defined
func DefineComposedType(name string, types ...Typed) TyComp {
	return func() (string, []Typed) {
		return name, types
	}
}

// higher order type has the highest flag type assigned
func (t TyComp) FlagType() uint8  { return 254 }
func (t TyComp) TypeName() string { name, _ := t(); return name }
func (t TyComp) Types() []Typed   { _, flags := t(); return flags }
func (t TyComp) NatFlags() []d.TyNat {
	var flags = []d.TyNat{}
	for _, flag := range t.Types() {
		if flag.FlagType() == 1 {
			flags = append(flags, flag.(d.TyNat))
		}
	}
	return flags
}
func (t TyComp) FncFlags() []TyFnc {
	var flags = []TyFnc{}
	for _, flag := range t.Types() {
		if flag.FlagType() == 2 {
			flags = append(flags, flag.(TyFnc))
		}
	}
	return flags
}
func (t TyComp) Flag() d.BitFlag {
	var flags d.BitFlag
	for _, flag := range t.Types() {
		flags = flags | flag.Flag()
	}
	return flags
}
func (t TyComp) TypeFnc() TyFnc {
	var flags TyFnc
	for _, flag := range t.FncFlags() {
		flags = flags | flag
	}
	return flags
}
func (t TyComp) TypeNat() d.TyNat {
	var flags d.TyNat
	for _, flag := range t.NatFlags() {
		flags = flags | flag
	}
	return flags
}
func (t TyComp) Call(...Callable) Callable {
	var args = []Callable{}
	for _, arg := range t.FncFlags() {
		args = append(args, arg)
	}
	return NewVector(args...)
}
func (t TyComp) Eval(...d.Native) d.Native {
	var args = []d.Native{}
	for _, arg := range t.NatFlags() {
		args = append(args, arg)
	}
	return d.NewSlice(args...)
}
func (t TyComp) Match(arg d.Typed) bool { return t.Flag().Match(arg) }
