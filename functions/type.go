package functions

import (
	s "strings"

	d "github.com/joergreinhardt/gatwd/data"
)

type (
	TyDef     func() (string, []Expression)
	TyFlag    d.Uint8Val
	TyFnc     d.BitFlag
	Arity     d.Int8Val
	Propertys d.Int8Val
)

//go:generate stringer -type TyFlag
const (
	Flag_BitFlag TyFlag = 0 + iota
	Flag_KeyWord
	Flag_DataCons
	Flag_Function
	Flag_Native
	Flag_Token
	Flag_Arity
	Flag_Prop
	Flag_Lex

	Flag_Definition TyFlag = 255
)

func (t TyFlag) U() d.Uint8Val { return d.Uint8Val(t) }
func (t TyFlag) Match(match d.Uint8Val) bool {
	if match == t.U() {
		return true
	}
	return false
}

func typedToExpression(typ Typed) Expression {
	var expr Expression
	switch {
	case Flag_Native.Match(typ.FlagType()):
		expr = NewData(typ.(d.TyNat))
	case Flag_Function.Match(typ.FlagType()):
		expr = typ.(TyFnc)
	case Flag_DataCons.Match(typ.FlagType()):
		expr = typ.(TyFnc)
	case Flag_Arity.Match(typ.FlagType()):
		expr = typ.(Arity)
	case Flag_Prop.Match(typ.FlagType()):
		expr = typ.(Propertys)
	case Flag_KeyWord.Match(typ.FlagType()):
		expr = typ.(TyKeyWord)
	case Flag_Lex.Match(typ.FlagType()):
		expr = typ.(TyLex)
	case Flag_Definition.Match(typ.FlagType()):
		expr = typ.(TyDef)
	}
	return expr
}

//go:generate stringer -type=TyFnc
const (
	/// GENERIC TYPE
	Type TyFnc = 1 << iota
	/// FUNCTION TYPES
	Data
	Constant
	/// PARAMETER OPTIONS
	Property
	Argument
	Return
	Index
	Key
	/// TRUTH VALUE OTIONS
	True
	False
	Undecided
	/// ORDER VALUE OPTIONS
	EQ
	LT
	GT
	/// BOUND VALUE OPTIONS
	Min
	Max
	/// VALUE OPTIONS
	Switch
	Case
	Then
	Else
	Just
	None
	OR
	EI
	/// DATA TYPE CLASSES
	Numbers
	Letters
	Bytes
	Text
	/// SUM COLLECTION TYPES
	List
	Vector
	/// PRODUCT COLLECTION TYPES
	Pair
	Set
	Enum
	Tuple
	Record
	/// IMPURE
	State
	IO
	/// HIGHER ORDER TYPE
	HigherORder

	Kinds = Type | Data | Constant

	//// PARAMETERS
	Signature = Argument | Return
	Param     = Key | Index | Property

	Params = Signature | Param

	//// TRUTH & COMPARE
	Truth   = True | False
	Trinary = Truth | Undecided
	CMP     = LT | GT | EQ

	//// OPTIONALS
	If     = Then | Else
	Maybe  = Just | None
	Option = EI | OR

	Branches = Switch | Case | If | Maybe | Option

	//// COLLECTIONS
	Consumeables = List | Vector
	Collections  = Consumeables | Pair |
		Set | Record | Enum | Tuple

	AllTypes = Kinds | Params |
		Branches | Collections
)

//// TYPE DEFINITION
func Define(name string, retype Expression, paratypes ...Expression) TyDef {
	return func() (string, []Expression) {
		return name, append([]Expression{retype}, paratypes...)
	}
}

func (t TyDef) Type() Typed                        { return t }
func (t TyDef) String() string                     { return t.TypeName() }
func (t TyDef) Name() string                       { var name, _ = t(); return name }
func (t TyDef) Elems() []Expression                { var _, expr = t(); return expr }
func (t TyDef) Return() Expression                 { return t.Elems()[0] }
func (t TyDef) FlagType() d.Uint8Val               { return Flag_Definition.U() }
func (t TyDef) Flag() d.BitFlag                    { return t.Return().TypeFnc().Flag() }
func (t TyDef) TypeFnc() TyFnc                     { return t.Return().TypeFnc() }
func (t TyDef) Call(args ...Expression) Expression { return t }
func (t TyDef) Pattern() []Expression {
	var elems = t.Elems()
	if len(elems) > 1 {
		return elems[1:]
	}
	return []Expression{Type}
}
func (t TyDef) Arity() Arity {
	return Arity(len(t.Pattern()))
	return Arity(0)
}
func (t TyDef) ReturnName() string {
	var retname = t.Return().TypeName()
	if s.Contains(retname, " ") {
		retname = "(" + retname + ")"
	}
	return retname
}
func (t TyDef) PatternName() string {
	if t.Arity() > Arity(0) {
		var slice []string
		var sep = " → "
		var pattern = t.Pattern()
		if len(pattern) > 0 {
			for _, arg := range pattern {
				slice = append(slice,
					arg.TypeName())
			}
			return s.Join(slice, sep)
		}
	}
	return ""
}
func (t TyDef) TypeName() string {
	var sep = " → "
	var name = t.Name()
	if s.Contains(name, " ") {
		name = "(" + name + ")"
	}
	if name == "" {
		name = t.ReturnName()
	}
	if t.Arity() > Arity(0) {
		var slice []string
		slice = append(slice, t.PatternName(),
			name, t.ReturnName())
		return s.Join(slice, sep)
	}
	return name
}
func (t TyDef) Match(typ d.Typed) bool {
	switch typ.FlagType() {
	case Flag_BitFlag.U():
	case Flag_Native.U():
	case Flag_Function.U():
	}
	return false
}

// type TyFnc d.BitFlag
// encodes the kind of functional data as bitflag
func (t TyFnc) TypeFnc() TyFnc                     { return Type }
func (t TyFnc) TypeNat() d.TyNat                   { return d.Type }
func (t TyFnc) Flag() d.BitFlag                    { return d.BitFlag(t) }
func (t TyFnc) Uint() d.UintVal                    { return d.BitFlag(t).Uint() }
func (t TyFnc) FlagType() d.Uint8Val               { return Flag_Function.U() }
func (t TyFnc) Match(arg d.Typed) bool             { return t.Flag().Match(arg) }
func (t TyFnc) Call(args ...Expression) Expression { return t.TypeFnc() }
func (t TyFnc) Eval() d.Native                     { return t.TypeNat() }
func (t TyFnc) Type() Typed                        { return Define(t.TypeName(), t) }
func (t TyFnc) TypeName() string {
	var count = t.Flag().Count()
	// loop to print concatenated type classes correcty
	if count > 1 {
		var delim = "|"
		var str string
		for i, flag := range t.Flag().Decompose() {
			str = str + TyFnc(flag.Flag()).String()
			if i < count-1 {
				str = str + delim
			}
		}
		return str
	}
	return t.String()
}

//// CALL PROPERTYS
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

// CALL PROPERTY FLAG
func (p Propertys) MatchProperty(arg Propertys) bool {
	if p&arg != 0 {
		return true
	}
	return false
}

// PROPERTY CONVIENIENCE METHODS
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

func FlagToProp(flag d.BitFlag) Propertys { return Propertys(flag.Uint()) }

func (p Propertys) Flag() d.BitFlag                    { return d.BitFlag(uint64(p)) }
func (p Propertys) FlagType() d.Uint8Val               { return Flag_Prop.U() }
func (p Propertys) TypeNat() d.TyNat                   { return d.Type }
func (p Propertys) TypeFnc() TyFnc                     { return Type }
func (p Propertys) TypeName() string                   { return "Propertys" }
func (p Propertys) Match(flag d.Typed) bool            { return p.Flag().Match(flag) }
func (p Propertys) Eval() d.Native                     { return d.Int8Val(p) }
func (p Propertys) Call(args ...Expression) Expression { return p }
func (p Propertys) Type() Typed {
	return Define(p.TypeName(), Property)
}

//// CALL ARITY
///
// arity of well defined callables
//
//go:generate stringer -type Arity
const (
	Nary Arity = -1 + iota
	Nullary
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

func (a Arity) Eval() d.Native                { return a }
func (a Arity) FlagType() d.Uint8Val          { return Flag_Arity.U() }
func (a Arity) Int() int                      { return int(a) }
func (a Arity) Type() Typed                   { return Type }
func (a Arity) TypeFnc() TyFnc                { return Type }
func (a Arity) TypeNat() d.TyNat              { return d.Type }
func (a Arity) Match(arg d.Typed) bool        { return a == arg }
func (a Arity) TypeName() string              { return a.String() }
func (a Arity) Call(...Expression) Expression { return NewData(a) }
func (a Arity) Flag() d.BitFlag               { return d.BitFlag(a) }
