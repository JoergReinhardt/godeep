package functions

import (
	"strings"

	d "github.com/joergreinhardt/gatwd/data"
	"github.com/joergreinhardt/gatwd/lex"
)

type (
	TyNest    func() (string, []Typed)
	TyTup     func() (string, []TyNest)
	TyRec     func() (string, []TyNest)
	TySig     func() (string, Propertys, TyTup)
	TyFnc     d.BitFlag
	Arity     d.Int8Val
	Propertys d.Uint8Val
)

//// nested TYPE
///
// composition type to define higher order types. it returns a type name that
// has either been passed during creation, or derived from the type names of
// it's elements and a slice of instances implementing the typed interface to
// implement recursively nested higher order types of arbitrary depth and
// complexity.
func NewNestedType(name string, types ...Typed) TyNest {
	if name == "" {
		name = deriveTypeName(types...)
	}
	return func() (string, []Typed) {
		return name, types
	}
}

// helper function to derive types name if no explicit name was passed to
// nested type constructor
func deriveTypeName(types ...Typed) string {
	var str string
	var length = len(types)
	for n, t := range types {
		str = str + t.TypeName()
		if n < length-1 {
			str = str + " "
		}
	}
	return str
}

// higher order type has the highest possible flag type value.
func (t TyNest) FlagType() uint8 { return 255 }

// return all sub type flags of nested type
func (t TyNest) SubTypes() []Typed { _, flags := t(); return flags }

// return name of nested type
func (t TyNest) TypeName() string { name, _ := t(); return name }

// returns string representation of recursively nested type.
func (t TyNest) String() string {
	var str = t.TypeName() + " "
	var length = len(t.SubTypes())
	for n, sub := range t.SubTypes() {
		str = str + sub.TypeName()
		if n < length-1 {
			str = str + " "
		}
	}
	return str
}

// match method of nested higher order type either matches its sub-flag(s) to
// the (possibly OR concatenated) argument flag(s) or matches type names and
// sub flags to nested flags, that may be nested recursively.
func (t TyNest) Match(arg d.Typed) bool {
	// return all flags of this type
	var flags = t.SubTypes()
	// if flag argument is non-nested, name is ignored.
	if arg.FlagType() < uint8(255) {
		// if types flag is atomic (only applys to instances of native)
		if len(flags) == 1 {
			// if argument is atomic as well
			if arg.Flag().Count() == 1 {
				// check if flag type matches
				if flags[0].FlagType() == arg.FlagType() {
					// check if flags match
					if flags[0].Match(arg) {
						return true
					}
					return false
				}
			}
		}
		// argument type is OR concatenated, decompose argument type.
		// (applys to all non-nested callables)
		var args = arg.Flag().Decompose()
		// number of concatenated flags match
		if len(args) == len(flags) {
			// range over flags and compare each to corresponding
			// argument
			for n, flag := range flags {
				// return false for first non-matching flag
				if !flag.Match(args[n]) {
					return false
				}
				// return true if all flags match
				return true
			}
		}
	}
	// if argument type itself is nested, check if type names match first.
	if strings.Compare(t.TypeName(), arg.TypeName()) == 0 {
		// cast argument as nested type and return all its sub flags
		var args = arg.(TyNest).SubTypes()
		// range over sub flags to compare each.
		for n, flag := range flags {
			// return false on first non matching sub flag. calls
			// nested types match method recursively for nested
			// nested-type sub flags.
			if !flag.Match(args[n]) {
				return false
			}
			// return true if all sub flags match
			return true
		}
	}
	// type names of nested argument type did not match → return false
	return false
}

// filter all native type flags and return as slice
func (t TyNest) NatFlags() []d.TyNat {
	var flags = []d.TyNat{}
	for _, flag := range t.SubTypes() {
		if flag.FlagType() == 1 {
			flags = append(flags, flag.(d.TyNat))
		}
	}
	return flags
}

// filter all functional type flags and return as slice
func (t TyNest) FncFlags() []TyFnc {
	var flags = []TyFnc{}
	for _, flag := range t.SubTypes() {
		if flag.FlagType() == 2 {
			flags = append(flags, flag.(TyFnc))
		}
	}
	return flags
}

// OR concatenate all sub flags cast as bit-flag
func (t TyNest) Flag() d.BitFlag {
	var flags d.BitFlag
	for _, flag := range t.SubTypes() {
		flags = flags | flag.Flag()
	}
	return flags
}

// OR concatenate all function type flags
func (t TyNest) TypeFnc() TyFnc {
	var flags TyFnc
	for _, flag := range t.FncFlags() {
		flags = flags | flag
	}
	return flags
}

// OR concatenate all native type flags
func (t TyNest) TypeNat() d.TyNat {
	var flags d.TyNat
	for _, flag := range t.NatFlags() {
		flags = flags | flag
	}
	return flags
}

// eval method returns data slice of name and all sub flags
func (t TyNest) Eval() d.Native {
	// allocate a slice of natives containing type name as first element
	var natives = []d.Native{d.StrVal(t.TypeName())}
	// range over sub flags and append each one to slice
	for _, flag := range t.SubTypes() {
		natives = append(natives, flag)
	}
	return d.NewSlice(natives...)
}

// call method returns a vector instance containing name and all sub type flags
func (t TyNest) Call(...Callable) Callable {
	// allocate slice of callables with type name as first element
	var args = []Callable{New(t.TypeName())}
	// range over sub flags and concatenate each to slice
	for _, flag := range t.SubTypes() {
		args = append(args, New(flag))
	}
	return NewVector(args...)
}

///////////////////////////////////////////////////////////////////////////////
//// TUPLE TYPE
///

///////////////////////////////////////////////////////////////////////////////
//// RECORD TYPE
///

///////////////////////////////////////////////////////////////////////////////
//// SIGNATURE TYPE
///

///////////////////////////////////////////////////////////////////////////////
//go:generate stringer -type=TyFnc
const (
	/// KIND FLAGS ///
	Type TyFnc = 1 << iota
	Data
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
	Predicate
	Undecided
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

	Kinds = Type | Data | Functor

	Truth = Undecided | False | True

	Ordered = Equal | Lesser | Greater

	Maybe = Just | None

	Alternatives = Either | Or

	Branch = If | Else

	Continue = Do | While

	IO = Buffer | Reader | Writer

	Consumeables = Collections | Applicable | Monad | IO
)

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
func (t TyFnc) Eval() d.Native                 { return t.TypeNat() }
func (t TyFnc) Flag() d.BitFlag                { return d.BitFlag(t) }
func (t TyFnc) Match(arg d.Typed) bool         { return t.Flag().Match(arg) }
func (t TyFnc) Uint() uint                     { return d.BitFlag(t).Uint() }

//// CALL PROPERTYS
///
// propertys of well defined callables
//
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

func FlagToProp(flag d.BitFlag) Propertys { return Propertys(uint8(flag.Uint())) }
func (p Propertys) PostFix() bool         { return p.Flag().Match(PostFix.Flag()) }
func (p Propertys) InFix() bool           { return !p.Flag().Match(PostFix.Flag()) }
func (p Propertys) Atomic() bool          { return p.Flag().Match(Atomic.Flag()) }
func (p Propertys) Thunk() bool           { return !p.Flag().Match(Atomic.Flag()) }
func (p Propertys) Eager() bool           { return p.Flag().Match(Eager.Flag()) }
func (p Propertys) Lazy() bool            { return !p.Flag().Match(Eager.Flag()) }
func (p Propertys) RightBound() bool      { return p.Flag().Match(RightBound.Flag()) }
func (p Propertys) LeftBound() bool       { return !p.Flag().Match(RightBound.Flag()) }
func (p Propertys) Mutable() bool         { return p.Flag().Match(Mutable.Flag()) }
func (p Propertys) Imutable() bool        { return !p.Flag().Match(Mutable.Flag()) }
func (p Propertys) SideEffect() bool      { return p.Flag().Match(SideEffect.Flag()) }
func (p Propertys) Pure() bool            { return !p.Flag().Match(SideEffect.Flag()) }
func (p Propertys) Primitive() bool       { return p.Flag().Match(Primitive.Flag()) }
func (p Propertys) Parametric() bool      { return !p.Flag().Match(Primitive.Flag()) }
func (p Propertys) TypeNat() d.TyNat      { return d.Flag }
func (p Propertys) TypeFnc() TyFnc        { return HigherOrder }
func (p Propertys) Flag() d.BitFlag       { return d.BitFlag(uint64(p)) }
func (p Propertys) Eval() d.Native        { return d.Int8Val(p) }

func (p Propertys) Match(flag d.BitFlag) bool      { return p.Flag().Match(flag) }
func (p Propertys) Call(args ...Callable) Callable { return p }
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

//// CALL ARITY
///
// arity of well defined callables
//
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

func (a Arity) Eval() d.Native            { return d.Int8Val(a) }
func (a Arity) Call(...Callable) Callable { return NewNative(a.Eval()) }
func (a Arity) Int() int                  { return int(a) }
func (a Arity) Flag() d.BitFlag           { return d.BitFlag(a) }
func (a Arity) TypeNat() d.TyNat          { return d.Flag }
func (a Arity) TypeFnc() TyFnc            { return HigherOrder }
func (a Arity) Match(arg Arity) bool      { return a == arg }
