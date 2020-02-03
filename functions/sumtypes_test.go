package functions

import (
	"fmt"
	"testing"

	d "github.com/joergreinhardt/gatwd/data"
)

var (
	addInts = Define(Dat(func(args ...d.Native) d.Native {
		var a, b = args[0].(d.IntVal), args[1].(d.IntVal)
		return a + b
	}),
		DecSym("+"),
		Declare(Dat(0).Type()),
		Declare(Dat(0).Type(), Dat(0).Type()),
	)
	addUints = Define(Dat(func(args ...d.Native) d.Native {
		var a, b = args[0].(d.UintVal), args[1].(d.UintVal)
		return a + b
	}),
		DecSym("+"),
		Declare(Dat(uint(0)).Type()),
		Declare(Dat(uint(0)).Type(), Dat(uint(0)).Type()),
	)
	addFloats = Define(Dat(func(args ...d.Native) d.Native {
		var a, b = args[0].(d.FltVal), args[1].(d.FltVal)
		return a + b
	}),
		DecSym("+"),
		Declare(Dat(0.0).Type()),
		Declare(Dat(0.0).Type(), Dat(0.0).Type()),
	)
)

func TestExpression(t *testing.T) {
	fmt.Printf(`
 defines expression to perform addition on integers of type data/Native.

 - should not take any arguments except d.IntVal
   - first argument is a symbol definition expressing the name
   - second argument is the return type (derived from instance)
   - third argument is the arguments types (derived from instances) wrapped
     by a composed type
 - should return a partial, when only one argument is passed
 - should return atomic integer result, when two args are passed
 - should return a vector of resulting integers‥.
 - ‥.where the last element might be a partialy applyed addition, if an odd
   number of arguments has been passed` + "\n\n")
	fmt.Printf("add ints expression definition type: %s\n"+
		"type-ident: %s\ntype-args: %s\nreturn type: %s\n",
		addInts.Type(), addInts.TypeId(),
		addInts.TypeArgs(), addInts.TypeRet())

	fmt.Printf("addInts: %s argtype : %s identype: %s, retype: %s fnctype: %s\n",
		addInts,
		addInts.Type().TypeArgs(),
		addInts.Type().TypeId(),
		addInts.Type().TypeRet(),
		addInts.Type().TypeFnc(),
	)

	var wrong = addInts.Call(Dat("string one"), Dat(true))
	fmt.Printf("called with argument of wrong type: %s\n", wrong)
	if !wrong.Type().Match(None) {
		t.Fail()
	}

	var partial = addInts.Call(Dat(23))
	fmt.Printf("partial: %s argtype : %s identype: %s, retype: %s\n",
		partial,
		partial.Type().TypeArgs(),
		partial.Type().TypeId(),
		partial.Type().TypeRet())

	fmt.Printf("manno ey: %s\n", partial.Type())
	if !partial.Type().Match(Partial) {
		t.Fail()
	}

	var wrongpart = partial.Call(Dat("string"))
	fmt.Printf("partial called with argument of wrong type: %s\n", wrongpart)
	if !wrongpart.Type().Match(None) {
		t.Fail()
	}

	var complete = partial.Call(Dat(42))
	fmt.Printf("complete: %s\n", complete)
	if data, ok := complete.(Evaluable); ok {
		if num, ok := data.Eval().(d.IntVal); ok {
			if num.Int() != 65 {
				t.Fail()
			}
		}
	}

	var result2 = addInts.Call(Dat(23), Dat(42))
	fmt.Printf("result2: %s argtype : %s identype: %s, retype: %s\n",
		result2, result2.Type().TypeArgs(),
		result2.Type().TypeId(),
		result2.Type().TypeRet())
	fmt.Printf("result2: %s\n", result2)
	if vec, ok := result2.(VecVal); ok {
		if vec.Len() != 2 {
			t.Fail()
		}
	}

	var result3 = addInts.Call(Dat(23), Dat(42), Dat(23))
	fmt.Printf("result3 type: %s type-ident: %s type-return: %s type-arguments: %s\n",
		result3.Type().String(),
		result3.Type().TypeId(),
		result3.Type().TypeRet(),
		result3.Type().TypeArgs(),
	)

	fmt.Printf("result3: %s\n", result3)
	fmt.Printf("result3 element 0 type: %s\n", result3.(VecVal)()[0].Type())
	if vec, ok := result3.(VecVal); ok {
		if !vec()[0].Type().Match(Data) {
			t.Fail()
		}
	}

	complete = result3.(VecVal)()[1].Call(Dat(42))
	fmt.Printf("completed result3[1] partial: %s\n", complete)
	if complete.(Atom).Eval().(d.Numeral).Int() != 65 {
		t.Fail()
	}

	var result4 = addInts.Call(
		Dat(23), Dat(42), Dat(23), Dat(42),
	)
	fmt.Printf("result4: %s\n", result4)
	if vec, ok := result4.(VecVal); ok {
		if !vec()[0].Type().Match(Dat(0).Type()) {
			t.Fail()
		}
	}
}
func TestTuple(t *testing.T) {
	var tup = DefTuple(Dat(0).Type(), Dat(uint(0)).Type(), Dat(0.0).Type())
	fmt.Printf("tuple constructor: %s\n", tup)

	var app = tup.Call(Dat(1), Dat(uint(1)), Dat(0.1))
	fmt.Printf("applied tuple: %s\n", app)
	var vec = app.(TupVal).Vector()
	fmt.Printf("vector: %s\n", vec)
	if vec()[0].(Atom)().(d.IntVal) != 1 &&
		vec()[0].(Atom)().(d.UintVal) != 1 &&
		vec()[0].(Atom)().(d.FltVal) != 0.1 {
		t.Fail()
	}

	var celltype = tup.GetCellType(2)
	fmt.Printf("type of cell₂ (expects float): %s\n", celltype)
	if !celltype.Match(Dat(0.0).Type()) {
		t.Fail()
	}

	var elem = app.(TupVal).Get(2)
	fmt.Printf("tuple idx₂ == 0.1? %s\n", elem)
	if elem.(Atom)().(d.FltVal) != 0.1 {
		t.Fail()
	}

	var elems = app.Call(Dat(0), Dat(2))
	fmt.Printf("tuple idx₀₂ == [1, 0.1]? %s\n", elems)
	if elems.(TupVal).Vector()()[0].(Atom)().(d.IntVal) != 1 && elems.(TupVal).Vector()()[1].(Atom)().(d.FltVal) != 0.1 {
		t.Fail()
	}

	app = tup.Call(Dat(1))
	fmt.Printf("partialy applied to tuple: %s\n", app)
	if !app.TypeFnc().Match(Partial) {
		t.Fail()
	}

	app = app.Call(Dat(uint(1)))
	fmt.Printf("partialy applied to tuple: %s\n", app)
	if !app.TypeFnc().Match(Partial) {
		t.Fail()
	}

	app = app.Call(Dat(1.1))
	fmt.Printf("completed tuple application: %s\n", app)
	fmt.Printf("return type: %s\n", app.Type())
	if !app.TypeFnc().Match(Tuple) {
		t.Fail()
	}

	app = tup.Call(Dat(1.1), Dat(true), Dat("string"))
	fmt.Printf("wrong args applied to tuple: %s\n", app)
	if !app.TypeFnc().Match(None) {
		t.Fail()
	}

}

func TestCase(t *testing.T) {
}
func TestSwitch(t *testing.T) {
}
func TestRecord(t *testing.T) {
}
