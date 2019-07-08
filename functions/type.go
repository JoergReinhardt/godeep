package functions

import (
	"fmt"
	s "strings"
	u "unicode"

	d "github.com/joergreinhardt/gatwd/data"
)

type (
	TySig     func() (TyDef, []TyDef)
	TyDef     func() (string, []Typed)
	TyFlag    d.Uint8Val
	TyFnc     d.BitFlag
	Arity     d.Int8Val
	Propertys d.Int8Val
)

//go:generate stringer -type TyFlag
const (
	Flag_BitFlag TyFlag = 0 + iota
	Flag_Native
	Flag_Functional
	Flag_DataCons
	Flag_Arity
	Flag_Prop

	Flag_Def TyFlag = 255
)

func (t TyFlag) U() d.Uint8Val { return d.Uint8Val(t) }
func (t TyFlag) Match(match d.Uint8Val) bool {
	if match == t.U() {
		return true
	}
	return false
}

//go:generate stringer -type=TyFnc
const (
	/// GENERIC TYPE
	Type TyFnc = 1 << iota
	/// FUNCTION TYPES
	Data
	Constant
	Function
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
	Lesser
	Greater
	Equal
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
	Either
	Or
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
	HigherOrder

	Kinds = Type | Data | Constant | Function

	//// PARAMETERS
	Signature = Argument | Return
	Parameter = Key | Index | Property

	Parameters = Signature | Parameter

	//// TRUTH & COMPARE
	Truth   = True | False
	Trinary = Truth | Undecided
	Compare = Lesser | Greater | Equal

	Tests = Truth | Trinary | Compare

	//// OPTIONALS
	If     = Then | Else
	Maybe  = Just | None
	Option = Either | Or

	Branches = Switch | Case | If | Maybe | Option

	//// COLLECTIONS
	Consumeables = List | Vector
	Collections  = Consumeables | Pair |
		Set | Record | Enum | Tuple

	AllTypes = Kinds | Parameters | Tests |
		Branches | Collections
)

var fncTypes = map[string]TyFnc{}
var natTypes = map[string]d.TyNat{}

func init() {
	for _, nat := range d.FetchTypes() {
		natTypes[nat.TypeName()] = nat
	}
	for _, fnc := range fetchTypes() {
		fncTypes[fnc.TypeName()] = fnc
	}
}
func searchNatType(name string) (d.TyNat, bool) {
	if val, ok := natTypes[name]; ok {
		return val, ok
	}
	return d.Nil, false
}
func searchFncType(name string) (TyFnc, bool) {
	if val, ok := fncTypes[name]; ok {
		return val, ok
	}
	return None, false
}
func fetchTypes() []TyFnc {
	var tt = []TyFnc{}
	var i uint
	var t TyFnc = 0
	for t < Type {
		t = 1 << i
		i = i + 1
		tt = append(tt, TyFnc(t))
	}
	return tt
}

type TvKind uint8

//go:generate stringer -type=TvKind
const (
	FunctionName TvKind = 0 + iota
	FunctionType
	NativeType
	ParamType
	ClassType
)

type tval struct {
	Kind    TvKind
	NatType d.TyNat
	FncType TyFnc
}

func newTval(kind TvKind, nat d.TyNat, fnc TyFnc) tval {
	return tval{kind, nat, fnc}
}

type tvalm map[string]tval

func splitSignature(signature string) VecCol {
	var slice = []Expression{}
	var str = s.Split(signature, " ")
	for _, str := range str {
		slice = append(slice, NewData(d.StrVal(str)))
	}
	return NewVector(slice...)
}

func lDelim(tm tvalm, sig, elems VecCol) (tvalm, VecCol, VecCol) {
	if sig.Len() > 0 {
		var tok = sig.Head().String()
		if s.HasPrefix(tok, "(") || tok == "(" {
			_, sig = sig.ConsumeVec()
			// shave of the delimiter from the token
			tok = s.TrimLeft(tok, "(")
			// generate a new token to replace the popped one
			sig = sig.Append(NewData(d.StrVal(tok)))
			// create sub element to take the delimiter expression
			var subelems = NewVector()
			// pass on and reassing the subelement,signature and type map
			tm, sig, subelems = parseSig(tm, sig, subelems)
			// append the sub element as element
			elems = elems.Append(subelems)
		}
	}
	return tm, sig, elems
}

func rDelim(tm tvalm, sig, elems VecCol) (tvalm, VecCol, VecCol) {
	if sig.Len() > 0 {
		var tok = sig.Head().String()
		if s.HasSuffix(tok, ")") || tok == ")" {
			_, sig = sig.ConsumeVec()
			// shave delimiter and replace popped toke
			tok = s.TrimRight(tok, ")")
			sig = sig.Append(NewData(d.StrVal(tok)))
		}
	}
	return tm, sig, elems
}

func parseSigElem(tm tvalm, sig, elems VecCol) (tvalm, VecCol, VecCol) {
	if sig.Len() > 0 {
		var val tval
		var tok = sig.Head().String()
		if u.IsUpper([]rune(tok)[0]) {
			if nat, ok := searchNatType(tok); ok {
				_, sig = sig.ConsumeVec()
				elems = elems.Append(NewData(d.StrVal(tok)))
				val = newTval(NativeType, nat, None)
				tm[tok] = val
				return tm, sig, elems
			}
			if fnc, ok := searchFncType(tok); ok {
				_, sig = sig.ConsumeVec()
				elems = elems.Append(NewData(d.StrVal(tok)))
				val = newTval(FunctionType, d.Nil, fnc)
				tm[tok] = val
				return tm, sig, elems
			}
			_, sig = sig.ConsumeVec()
			elems = elems.Append(NewData(d.StrVal(tok)))
			val = newTval(ParamType, d.Nil, None)
			tm[tok] = val
			return tm, sig, elems
		}
		if elems.Len() > 0 {
			if elems.Last().TypeFnc().Match(Vector) {
				fmt.Printf("last element matched vector")
				return tm, sig, elems
			}
		}
		_, sig = sig.ConsumeVec()
		elems = elems.Append(NewData(d.StrVal(tok)))
		val = newTval(FunctionName, d.Nil, None)
		tm[tok] = val
		return tm, sig, elems
	}
	return tm, sig, elems
}

func stripArrows(tm tvalm, sig, elems VecCol) (tvalm, VecCol, VecCol) {
	if sig.Len() > 0 {
		var tok = sig.Head().String()
		if s.ContainsAny(tok, "∷:->→=>⇒") {
			// pop the arrow token
			_, sig = sig.ConsumeVec()
		}
	}
	return tm, sig, elems
}

func parseSig(tm tvalm, sig, elems VecCol) (tvalm, VecCol, VecCol) {
	for sig.Len() > 0 {
		if elems.Len() > 0 {
			if elems.Last().TypeFnc().Match(Vector) {
				fmt.Printf("last element matched vector")
				continue
			}
		}
		tm, sig, elems = stripArrows(tm, sig, elems)
		tm, sig, elems = lDelim(tm, sig, elems)
		tm, sig, elems = rDelim(tm, sig, elems)
		tm, sig, elems = parseSigElem(tm, sig, elems)
	}
	return tm, sig, elems
}

// remove arrows and continue examining the next token
// if left delimiter → parse subelement
//		if s.Contains(tok, "|") {
//			elem, tval := parseTypeValSet()
//			elems = append(elems, elem)
//			tvm[tok] = tval
//		}
// if right delim → pop after token has been consumed

//func parseTypeValSet(tok string) (Expression, tval) {
//	var fnc TyFnc
//	var nat TyNat
//	return NewKeyPair(tok, Type), newTval(TypeValue, Type.Flag())
//}

//// TYPE DEFINITION
func Define(name string, retype Typed, paratypes ...Typed) TyDef {
	return func() (string, []Typed) {
		return name, append([]Typed{retype}, paratypes...)
	}
}

func (t TyDef) Type() TyDef                        { return t }
func (t TyDef) String() string                     { return t.TypeName() }
func (t TyDef) Name() string                       { var name, _ = t(); return name }
func (t TyDef) Elems() []Typed                     { var _, expr = t(); return expr }
func (t TyDef) Return() Typed                      { return t.Elems()[0] }
func (t TyDef) FlagType() d.Uint8Val               { return Flag_Def.U() }
func (t TyDef) Flag() d.BitFlag                    { return t.Return().TypeFnc().Flag() }
func (t TyDef) TypeFnc() TyFnc                     { return t.Return().TypeFnc() }
func (t TyDef) Call(args ...Expression) Expression { return t }
func (t TyDef) Pattern() []Typed {
	var elems = t.Elems()
	if len(elems) > 1 {
		return elems[1:]
	}
	return []Typed{Type}
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
	case Flag_Functional.U():
	}
	return false
}

// type TyFnc d.BitFlag
// encodes the kind of functional data as bitflag
func (t TyFnc) TypeFnc() TyFnc                     { return Type }
func (t TyFnc) TypeNat() d.TyNat                   { return d.Type }
func (t TyFnc) Flag() d.BitFlag                    { return d.BitFlag(t) }
func (t TyFnc) Uint() d.UintVal                    { return d.BitFlag(t).Uint() }
func (t TyFnc) FlagType() d.Uint8Val               { return Flag_Functional.U() }
func (t TyFnc) Match(arg d.Typed) bool             { return t.Flag().Match(arg) }
func (t TyFnc) Call(args ...Expression) Expression { return t.TypeFnc() }
func (t TyFnc) Eval() d.Native                     { return t.TypeNat() }
func (t TyFnc) Type() TyDef                        { return Define(t.TypeName(), t) }
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
func (p Propertys) Type() TyDef {
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
func (a Arity) TypeFnc() TyFnc                { return Type }
func (a Arity) TypeNat() d.TyNat              { return d.Type }
func (a Arity) Match(arg d.Typed) bool        { return a == arg }
func (a Arity) TypeName() string              { return a.String() }
func (a Arity) Call(...Expression) Expression { return NewData(a) }
func (a Arity) Flag() d.BitFlag               { return d.BitFlag(a) }
