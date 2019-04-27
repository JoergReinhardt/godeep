package functions

import (
	d "github.com/joergreinhardt/gatwd/data"
	"github.com/joergreinhardt/gatwd/lex"
)

//go:generate stringer -type=TyFnc
const (
	/// TYPE RELATED FLAGS ///
	Type TyFnc = 1 << iota
	Constructor
	Expression
	CallPropertys
	CallArity
	Data
	/// FUNCTORS AND MONADS ///
	Applicable
	Operator
	Functor
	Monad
	/// MONADIC SUB TYPES ///
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
	/// TYPE CLASSES ///
	Number
	Index
	Symbol
	Error
	/// COLLECTION TYPES ///
	Constant
	Pair
	Tuple
	Enum
	Set
	List
	Record
	Vector
	/// HIGHER ORDER TYPE IS THE HIGHEST FLAG ///
	HigherOrder

	Truth = False | True

	Maybe = Just | None

	Ordered = Greater | Lesser

	Kind = Data | Expression

	TrinaryTruth = False | Truth | Undecided

	Equality = Greater | Lesser | Equal

	Morphisms = Constructor | Operator | Functor |
		Applicable | Monad

	Options = False | True | Just | None | Case |
		Switch | Either | Or | If | Else |
		While | Do | Truth | Maybe | Ordered |
		Kind | Equality | TrinaryTruth

	Collections = Pair | Tuple | Enum | Set |
		List | Vector | Record

	Classes = Truth | TrinaryTruth | Equality |
		Ordered | Number
)

//// FIRST ORDER TYPE FUNCTION
///
// higher order types form a type tree by association, featuring leave
// nodes independent off all other types and types composed of other
// composed-/, or leave node types.  composed types are either defined
// as fixed sum of their sub types (enum-/, tuple-/, record types, or
// sets like the function argument set of a static function, as defined
// at sum-types definition; or as paramrtric product types, able to
// generate derived types, according to the arguments types passed at
// time of call evaluation.
//
// parent-/, subtype relationships are created, whenever a native types
// flag has more than a single bit set.  all functional types are
// parametric composed types per default by providing the 'TypeNat'
// method to return it's own native type.
//
// all container types are type agnostic regarding their contained
// elements, beginning with pairs, native and functional alike, and
// need consideration as product types.
type (
	TyFnc     d.BitFlag
	Arity     d.Uint8Val
	Propertys d.Uint8Val
)

//////////////////////////////////////////////////////////////////////
//// SEMANTIC CALL PROPERTYS
///

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

func (p Propertys) PostFix() bool    { return p.Flag().Match(PostFix.Flag()) }
func (p Propertys) InFix() bool      { return !p.Flag().Match(PostFix.Flag()) }
func (p Propertys) Atomic() bool     { return p.Flag().Match(Atomic.Flag()) }
func (p Propertys) Thunk() bool      { return !p.Flag().Match(Atomic.Flag()) }
func (p Propertys) Eager() bool      { return p.Flag().Match(Eager.Flag()) }
func (p Propertys) Lazy() bool       { return !p.Flag().Match(Eager.Flag()) }
func (p Propertys) RightBound() bool { return p.Flag().Match(RightBound.Flag()) }
func (p Propertys) LeftBound() bool  { return !p.Flag().Match(RightBound.Flag()) }
func (p Propertys) Mutable() bool    { return p.Flag().Match(Mutable.Flag()) }
func (p Propertys) Imutable() bool   { return !p.Flag().Match(Mutable.Flag()) }
func (p Propertys) SideEffect() bool { return p.Flag().Match(SideEffect.Flag()) }
func (p Propertys) Pure() bool       { return !p.Flag().Match(SideEffect.Flag()) }
func (p Propertys) Primitive() bool  { return p.Flag().Match(Primitive.Flag()) }
func (p Propertys) Parametric() bool { return !p.Flag().Match(Primitive.Flag()) }

func (p Propertys) TypeNat() d.TyNative { return d.Flag }
func (p Propertys) TypeFnc() TyFnc      { return HigherOrder }

func (p Propertys) Flag() d.BitFlag { return d.BitFlag(uint64(p)) }

func FlagToProp(flag d.BitFlag) Propertys { return Propertys(uint8(flag.Uint())) }

func (p Propertys) Eval(a ...d.Native) d.Native { return p }

func (p Propertys) Call(args ...Callable) Callable { return p }

func (p Propertys) MatchProperty(arg Propertys) bool {
	if p&arg != 0 {
		return true
	}
	return false
}

func (p Propertys) Match(flag d.BitFlag) bool { return p.Flag().Match(flag) }
func (p Propertys) Print() string {

	var flags = p.Flag().Decompose()
	var str string
	var l = len(flags)

	if l > 1 {
		for i, typed := range flags {

			if typed.FlagType() == 0 {

				str = str + typed.(d.TyNative).String()
			}

			if typed.FlagType() == 1 {

				str = str + typed.(TyFnc).String()
			}

			if typed.FlagType() == 2 {

				str = str + typed.(lex.SyntaxItemFlag).String()
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

func (a Arity) Call(...Callable) Callable { return NewFromData(a.Eval()) }

func (a Arity) Int() int            { return int(a) }
func (a Arity) Flag() d.BitFlag     { return d.BitFlag(a) }
func (a Arity) TypeNat() d.TyNative { return d.Flag }
func (a Arity) TypeFnc() TyFnc      { return HigherOrder }

func (a Arity) Match(arg Arity) bool { return a == arg }

// type TyFnc d.BitFlag
// encodes the kind of functional data as bitflag
func (t TyFnc) FlagType() int8                 { return 2 }
func (t TyFnc) TypeName() string               { return t.String() }
func (t TyFnc) TypeFnc() TyFnc                 { return Type }
func (t TyFnc) TypeNat() d.TyNative            { return d.Flag }
func (t TyFnc) Call(args ...Callable) Callable { return t.TypeFnc() }
func (t TyFnc) Eval(args ...d.Native) d.Native { return t.TypeNat() }
func (t TyFnc) Flag() d.BitFlag                { return d.BitFlag(t) }
func (t TyFnc) Uint() uint                     { return d.BitFlag(t).Uint() }
