package functions

import (
	"fmt"
	"testing"

	d "github.com/joergreinhardt/gatwd/data"
)

func TestConstant(t *testing.T) {
	var cons = NewConstant(func() Expression {
		return NewNative("this is a constant")
	})
	fmt.Printf("constant string: %s\n", cons)
	fmt.Printf("constant call: %s\n", cons.Call())
	fmt.Printf("constant function type: %s\n", cons.TypeFnc().TypeName())
	fmt.Printf("constant type: %s\n", cons.Type())
}
func TestFunction(t *testing.T) {
	var fnc = NewFunction(func(args ...Expression) Expression {
		if len(args) > 0 {
			return NewNative(args[0].String())
		}
		return NewNative(None.String())
	}, Def(Data, d.String))
	fmt.Printf("fnc: %s\n", fnc)
	fmt.Printf("fnc('test'): %s\n", fnc(NewNative("test")))
	fmt.Printf("fnc.Call('test'): %s\n", fnc.Call(NewNative("test")))
	fmt.Printf("fnc.Type(): %s\n", fnc.Type())
}

func TestArgType(t *testing.T) {
	var at = DeclareArguments(Def(Data, d.Int), Def(Data, d.Int), Def(Data, d.Int))
	fmt.Printf("declared arguments: %s\n", at)
	if !at.Type().Match(Def(Def(Data, d.Int), Def(Data, d.Int), Def(Data, d.Int))) {
		t.Fail()
	}

	var result = at.Call(NewNative(1))
	fmt.Printf("match pass int: %s result type: %s\n", result, result.Type())
	if !result.Type().Match(Def(Data, d.Int)) {
		t.Fail()
	}

	result = at.Call(NewNative(1), NewNative(1), NewNative(1))
	fmt.Printf("match pass three ints: %s result type: %s\n", result, result.Type())
	if !result.Type().Match(Def(Vector, Def(Data, d.Int))) {
		t.Fail()
	}

	result = at.Call(NewNative(1.0))
	fmt.Printf("match pass float: %s result type: %s\n", result, result.Type())
	if !result.Type().Match(d.Float) {
		t.Fail()
	}

	at = DeclareArguments(Def(Data, d.Int), Def(Data, d.Float))
	fmt.Printf("declared arguments: %s\n", at)
	if !at.MatchArgs(NewNative(1), NewNative(1.0)) {
		t.Fail()
	}
}

func TestDeclaredExpression(t *testing.T) {

	var addInt = NewFunction(func(args ...Expression) Expression {
		var a, b = args[0].(DataConst).Eval().(d.IntVal), args[1].(DataConst).Eval().(d.IntVal)
		return NewData(a + b)
	}, Def(Data, d.Int))

	var result = addInt(NewNative(23), NewNative(42))
	fmt.Printf("result from applying ints to addInt: %s\n", result)

	var expr = DeclareExpression(addInt, Def(Data, d.Int), Def(Data, d.Int))
	fmt.Printf("declared expression: %s\n", expr)

	fmt.Printf("result from applying ints to addInt: %s\n",
		expr(NewNative(23), NewNative(42)))

	fmt.Printf("result from applying two floats to addInt: %s\n",
		expr(NewNative(23.0), NewNative(42.0)))

	fmt.Printf("result from applying four ints to addInt: %s\n",
		expr(NewNative(23), NewNative(42), NewNative(23), NewNative(42)))
	fmt.Printf("result type: %s\n", expr.Type())

	fmt.Printf("result from applying five ints to addInt: %s\n",
		expr(NewNative(23), NewNative(42), NewNative(23),
			NewNative(42), NewNative(42)))

	fmt.Printf("result from applying six ints to addInt: %s\n",
		expr(NewNative(23), NewNative(42), NewNative(23),
			NewNative(42), NewNative(42), NewNative(42)))

	result = expr(NewNative(23), NewNative(42), NewNative(23), NewNative(42),
		NewNative(42), NewNative(42), NewNative(42), NewNative(42))
	fmt.Printf("result from applying eight ints to addInt: %s\n", result)
	fmt.Printf("result from applying two more ints oversatisfyed expr: %s\n",
		result.Call(NewNative(42), NewNative(42)))

	var partial = expr.Call(NewNative(23))
	fmt.Printf("result from applying one int to addInt: %s, expression: %s arg type: %s len: %d\n",
		partial, partial.(DeclaredExpr).Unbox(), partial.(DeclaredExpr).ArgType(),
		partial.(DeclaredExpr).ArgType().Len())

	fmt.Printf("result from applying second int to partial addInt: %s\n",
		partial.Call(NewNative(42)))
}