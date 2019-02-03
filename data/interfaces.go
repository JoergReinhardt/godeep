package data

// VALUES AND TYPES
///////////////////
// propertys intendet for internal use
type Reproduceable interface{ Copy() Data }
type Destructable interface{ Clear() }
type Stringer interface{ String() string }

//// SER DEFINED DATA & FUNCTION TYPES ///////
type Typed interface {
	Flag() BitFlag
	String() string
} //<- lowest common denominator
type DataTyped interface {
	Flag() Type
	String() string
} //<- lowest common denominator
type DataType interface {
	Typed
}
type Data interface {
	Typed
}
type Ident interface {
	Data
	Ident() Data
}
type Evaluable interface {
	Eval() Data
}

type Nullable interface {
	Data
	Null() Data
}

type Accessor interface {
	Acc() Data
	Arg() Data
}
type Accessable interface {
	Get(acc Data) (Data, bool)
	Set(acc Data, dat Data)
}
type Paired interface {
	Left() Data
	Right() Data
	Both() (Data, Data)
}
type Mapped interface {
	Data
	Keys() []Data
	Data() []Data
	Accs() []Paired
	Get(acc Data) (Data, bool)
	Set(Data, Data) Mapped
}
type NativeVal interface {
	Data
	Null() func() Data
	DataFnc() func(Data) Data
}
type UnsignedVal interface{ Uint() uint }
type IntegerVal interface{ Int() int }
type Sliceable interface {
	Data
	Empty() bool
	Len() int
	Slice() []Data
}
type Vectorized interface {
	Data
	Len() int
	Empty() bool
	Slice() []Data
}
type NativeVec interface {
	Data
	Len() int
	Empty() bool
	Slice() []Data
}
type Collected interface {
	Data
	Empty() bool //<-- no more nil pointers & 'out of index'!
}
type Consumeable interface {
	Collected
	Head() Data
	Tail() Consumeable
	Shift() Consumeable
}
