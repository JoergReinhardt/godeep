package data

import (
	"fmt"
	"math/big"
	"time"
)

type CompoundError func() []error

func NewCompoundError(errs ...error) CompoundError {
	return func() []error { return errs }
}
func (e CompoundError) Error() string {
	var str string
	for n, e := range e() {
		str = str + string(n) + ": " + e.Error() + "\n"
	}
	return str
}

// returns a nil value instance
func NewNil() NilVal { return NilVal{} }
func NewError(errs ...error) ErrorVal {
	if len(errs) == 0 {
		return ErrorVal{}
	}
	if len(errs) == 1 {
		return ErrorVal{errs[0]}
	}
	return ErrorVal{NewCompoundError(errs...)}
}
func NewErrorFromString(str string) ErrorVal {
	return ErrorVal{fmt.Errorf(str)}
}

// returns a null value according to the native type passed in. if the flag
// turns out to be composed from multiple types, null value will be an instance
// of slice, pair, or map. otherwise an atomic native instance will be
// returned.
func NewUboxNull(nat TyNat) Native {
	return newUnboxed(nat)
}
func NewNull(nat TyNat) Native {
	// should the type flag turn out to be composed
	if nat.Flag().Count() > 0 {
		var nats []Typed
		switch {
		case nat.Match(Slice):
			// mask slice flag
			nats = nat.Flag().Mask(Slice).Flag().Decompose()
			// a single native type should be left
			if len(nats) == 1 {
				// return a slice of the particular type
				return NewUnboxed(nats[0].(NativeTyped).Type())
			}
		case nat.Match(Pair):
			// mask pair flag and decompose remaining flags
			nats = nat.Flag().Mask(Pair).Flag().Decompose()
			// two native types should be left
			if len(nats) == 2 {
				// return an empty pair composed of both type
				// elements
				return NewPair(
					newNull(nats[0].(TyNat)),
					newNull(nats[1].(TyNat)))
			}
		case nat.Match(Map):
			// mask the map flag and reassign native type
			nat = TyNat(nat.Flag().Mask(Map).Flag())
			// should leave the native type of the maps key
			if len(nats) > 0 {
				switch {
				case nat.Match(String):
					return NewStringMap()
				case nat.Match(Uint):
					return NewUintMap()
				case nat.Match(Type):
					return NewFLagMap()
				default:
					// if remaining type doesn't match any
					// of the existing set types, return a
					// generic set
					return NewValMap()
				}
			}
		default:
			// return the nil value, if the composed type flag
			// turns out to not be parseable
			return newNull(nat)
		}
	}
	// for non composed types, return an atomic null instance (returns a
	// nil type, if not parseable)
	return newNull(nat)
}

func New(vals ...interface{}) Native { dat, _ := newWithTypeInfo(vals...); return dat }

func NewData(args ...Native) Native {
	if len(args) > 0 {
		if len(args) > 1 {
			// try to return unboxed natives if possible, falls
			// back to return slice of native instances if not.
			return SliceToNatives(NewSlice(args...))
		}
		if args[0].Type() == Slice {
			return SliceToNatives(args[0].(DataSlice))
		}
		// a single native argument has been passed, return unchanged
		return args[0]
	}
	// no argument has been passed, return nil value
	return NilVal{}
}

func newUnboxedVector(f BitFlag, vals ...Native) Native { return conNativeVector(f, vals...) }

// converts untyped arguments to instances of native type, followed by a bit
// flag to indicate the derived type
func newWithTypeInfo(args ...interface{}) (rval Native, flag BitFlag) {

	// no arguments passed, return nil instance
	if len(args) == 0 {
		return nil, Nil.Type().Flag()
	}

	// multiple arguments have been passed
	if len(args) > 1 {

		// allocate slice of natives
		var nats = make([]Native, 0, len(args))

		// range over arguments
		for _, arg := range args {

			// allocate native instance to temporary assign
			// converted argument to, when created
			var nat Native

			// recursively create native instances and corresponding
			// type flags
			nat, flag = newWithTypeInfo(arg)

			// append native instance to preallocated slice of
			// natives
			nats = append(nats, nat)

			// OR concatenate flag type flags created by previously
			// converted arguments
			flag = flag | nat.Type().Flag()
		}

		// if flag length is one, all arguments yielded identical type.
		// return unboxed vector and type pure flag, to indicate all
		// members type
		if FlagLength(flag) == 1 {
			// return unboxed vector of natives
			return conNativeVector(flag, nats...), flag
		}

		// argument types are mixed, return slice of native instances
		// and multi typed flag
		return NewSlice(nats...), flag
	}

	// a single argument has been passed, assign to temporary value
	var temp = args[0]

	// switch on temporary values type, convert and assign corresponding
	// instance of typed native to return value.
	switch temp.(type) {
	case bool:
		rval = BoolVal(temp.(bool))
	case int, int64:
		rval = IntVal(temp.(int))
	case int8:
		rval = Int8Val(temp.(int8))
	case int16:
		rval = Int16Val(temp.(int16))
	case int32:
		rval = Int32Val(temp.(int32))
	case uint, uint64:
		rval = UintVal(temp.(uint))
	case uint16:
		rval = Uint16Val(temp.(uint16))
	case uint32:
		rval = Int32Val(temp.(int32))
	case float32:
		rval = Flt32Val(temp.(float32))
	case float64:
		rval = FltVal(temp.(float64))
	case complex64:
		rval = ImagVal(temp.(complex64))
	case complex128:
		rval = ImagVal(temp.(complex128))
	case byte:
		rval = ByteVal(temp.(byte))
	case []byte:
		rval = BytesVal(temp.([]byte))
	case string:
		rval = StrVal(temp.(string))
	case error:
		rval = ErrorVal{temp.(error)}
	case time.Time:
		rval = TimeVal(temp.(time.Time))
	case time.Duration:
		rval = DuraVal(temp.(time.Duration))
	case *big.Int:
		v := BigIntVal(*temp.(*big.Int))
		rval = &v
	case *big.Float:
		v := BigFltVal(*temp.(*big.Float))
		rval = &v
	case *big.Rat:
		v := RatioVal(*temp.(*big.Rat))
		rval = &v
	case func(...Native) Native:
		rval = Expression(temp.(func(...Native) Native))
	case TyNat:
		rval = BitFlag(temp.(TyNat))
	case BitFlag:
		rval = TyNat(temp.(BitFlag))
	case Native:
		rval = temp.(Native)
	case []Native:
		rval = DataSlice(temp.([]Native))
	}
	// return typed native instance and corresponding type flag
	return rval, rval.Type().Flag()
}

// returns unboxed vector from arguments of the type that has been passed as
// flag. argument types need to be prechecked!
func conNativeVector(flag BitFlag, args ...Native) (nat Sliceable) {

	var slice = []Native{}

	switch {
	case FlagMatch(flag, Nil.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(NilVal))
		}
	case FlagMatch(flag, Bool.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(BoolVal))
		}
	case FlagMatch(flag, Int.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(IntVal))
		}
	case FlagMatch(flag, Int8.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Int8Val))
		}
	case FlagMatch(flag, Int16.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Int16Val))
		}
	case FlagMatch(flag, Int32.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Int32Val))
		}
	case FlagMatch(flag, Uint.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(UintVal))
		}
	case FlagMatch(flag, Uint8.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Uint8Val))
		}
	case FlagMatch(flag, Uint16.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Uint16Val))
		}
	case FlagMatch(flag, Uint32.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Uint32Val))
		}
	case FlagMatch(flag, Float.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(FltVal))
		}
	case FlagMatch(flag, Flt32.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Flt32Val))
		}
	case FlagMatch(flag, Imag.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Imag64Val))
		}
	case FlagMatch(flag, Imag64.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(Imag64Val))
		}
	case FlagMatch(flag, Byte.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(ByteVal))
		}
	case FlagMatch(flag, Rune.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(RuneVal))
		}
	case FlagMatch(flag, Bytes.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(BytesVal))
		}
	case FlagMatch(flag, String.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(StrVal))
		}
	case FlagMatch(flag, BigInt.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(BigIntVal))
		}
	case FlagMatch(flag, BigFlt.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(BigFltVal))
		}
	case FlagMatch(flag, Ratio.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(RatioVal))
		}
	case FlagMatch(flag, Time.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(TimeVal))
		}
	case FlagMatch(flag, Duration.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(DuraVal))
		}
	case FlagMatch(flag, Error.Type().Flag()):
		for _, v := range args {
			slice = append(slice, v.(ErrorVal))
		}
	}
	return NewUnboxed(flag.Type(), slice...)
}

func NewEmpty(flag BitFlag) (val Native) {

	var value Native

	switch {

	case FlagMatch(flag, Nil.Type().Flag()):
		value = NilVal{}.Null()

	case FlagMatch(flag, Bool.Type().Flag()):
		value = BoolVal(false).Null()

	case FlagMatch(flag, Int.Type().Flag()):
		value = IntVal(0).Null()

	case FlagMatch(flag, Int8.Type().Flag()):
		value = Int8Val(0).Null()

	case FlagMatch(flag, Int16.Type().Flag()):
		value = Int16Val(0).Null()

	case FlagMatch(flag, Int32.Type().Flag()):
		value = Int32Val(0).Null()

	case FlagMatch(flag, Uint.Type().Flag()):
		value = UintVal(0).Null()

	case FlagMatch(flag, Uint8.Type().Flag()):
		value = Uint8Val(0).Null()

	case FlagMatch(flag, Uint16.Type().Flag()):
		value = Uint16Val(0).Null()

	case FlagMatch(flag, Uint32.Type().Flag()):
		value = Uint32Val(0).Null()

	case FlagMatch(flag, Float.Type().Flag()):
		value = FltVal(0).Null()

	case FlagMatch(flag, Flt32.Type().Flag()):
		value = Flt32Val(0).Null()

	case FlagMatch(flag, Imag.Type().Flag()):
		value = ImagVal(0).Null()

	case FlagMatch(flag, Imag64.Type().Flag()):
		value = Imag64Val(0).Null()

	case FlagMatch(flag, Byte.Type().Flag()):
		value = ByteVal(0).Null()

	case FlagMatch(flag, Rune.Type().Flag()):
		value = RuneVal(0).Null()

	case FlagMatch(flag, Bytes.Type().Flag()):
		value = BytesVal{}.Null()

	case FlagMatch(flag, String.Type().Flag()):
		value = StrVal(0).Null()

	case FlagMatch(flag, BigInt.Type().Flag()):
		value = BigIntVal{}.Null()

	case FlagMatch(flag, BigFlt.Type().Flag()):
		value = BigFltVal{}.Null()

	case FlagMatch(flag, Ratio.Type().Flag()):
		value = RatioVal{}.Null()

	case FlagMatch(flag, Time.Type().Flag()):
		value = TimeVal{}.Null()

	case FlagMatch(flag, Duration.Type().Flag()):
		value = DuraVal(0).Null()

	case FlagMatch(flag, Error.Type().Flag()):
		value = ErrorVal{nil}.Null()

	}

	return value
}
