package functions

import (
	"fmt"
	"testing"

	d "github.com/joergreinhardt/gatwd/data"
)

func TestData(t *testing.T) {
	var nat = DeclareData(d.UintVal(12))
	fmt.Printf("uint converted to native: %s, typeFnc: %s typeNat: %s \n",
		nat, nat.TypeFnc(), nat.Type())

	var pair = DeclareData(d.NewPair(Declare("key"), Declare(42)))
	fmt.Printf("key pair converted to native: %s, typeFnc: %s typeNat: %s\n",
		pair, pair.TypeFnc(), pair.Type())

	var nest = DeclareData(d.NewPair(d.StrVal("key"), d.NewPair(d.StrVal("inner key"), d.StrVal("Value"))))
	fmt.Printf("nested pair converted to native: %s, typeFnc: %s typeNat: %s\n",
		nest, nest.TypeFnc(), nest.Type())
}
