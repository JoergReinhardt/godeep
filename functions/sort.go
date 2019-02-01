/*
SORT & SEARCH

  implements golang sort and search slices of data and pairs of data. since
  'data' can be of collection type, this implements search and sort for pretty
  much every type thinkable of. generalizes over contained types by using godeeps
  capabilitys.
*/
package functions

import (
	"fmt"
	"math/big"
	"sort"
	"strings"

	d "github.com/JoergReinhardt/godeep/data"
)

// type class based comparison functions
func compareSymbolic(a, b Symbolic) int {
	return strings.Compare(a.String(), b.String())
}
func compareRational(a, b d.RatioVal) int {
	return ((*big.Rat)(&a)).Cmp((*big.Rat)(&b))
}
func compareUnsigned(a, b Unsigned) int {
	if a.Uint() < b.Uint() {
		return -1
	}
	if a.Uint() > b.Uint() {
		return 1
	}
	return 0
}
func compareInteger(a, b Integer) int {
	if a.Int() < b.Int() {
		return -1
	}
	if a.Int() > b.Int() {
		return 1
	}
	return 0
}
func compareIrrational(a, b Irrational) int {
	if a.Float() < b.Float() {
		return -1
	}
	if a.Float() > b.Float() {
		return 1
	}
	return 0
}
func compareFlag(a, b d.BitFlag) int {
	if a.Flag() < b.Flag() {
		return -1
	}
	if a.Flag() > b.Flag() {
		return 1
	}
	return 0
}
func compInt2BooIncl(i int) bool { return i >= 0 }
func compInt2BooExcl(i int) bool { return i > 0 }

// type to sort slices of data
type dataSorter []Functional

func (d dataSorter) Empty() bool {
	if len(d) > 0 {
		for _, dat := range d {
			if !ElemEmpty(dat) {
				return false

			}
		}
	}
	return true
}
func newDataSorter(dat ...Functional) dataSorter { return dataSorter(dat) }
func (d dataSorter) Len() int                    { return len(d) }
func (d dataSorter) Swap(i, j int)               { d[i], d[j] = d[j], d[i] }
func (ds dataSorter) Sort(argType d.Type) {
	sort.Slice(ds, newDataLess(argType, ds))
}
func (ds dataSorter) Search(praed Functional) int {
	var idx = sort.Search(len(ds), newDataFind(ds, praed))
	if idx < len(ds) {
		if strings.Compare(ds[idx].String(), praed.String()) == 0 {
			return idx
		}
	}
	return -1
}

func newDataLess(argType d.Type, ds dataSorter) func(i, j int) bool {
	var f = argType.Flag()
	switch {
	case f.Match(d.Symbolic.Flag()):
		return func(i, j int) bool {
			if strings.Compare(ds[j].(Symbolic).String(), ds[i].(Symbolic).String()) >= 0 {
				return true
			}
			return false
		}
	case f.Match(d.Flag.Flag()):
		return func(i, j int) bool {
			return ds[j].Eval().Flag() == ds[i].Eval().Flag()
		}
	case f.Match(d.Unsigned.Flag()):
		return func(i, j int) bool {
			return ds[j].Eval().(Unsigned).Uint() == ds[i].Eval().(Unsigned).Uint()
		}
	case f.Match(d.Integer.Flag()):
		return func(i, j int) bool {
			return ds[j].Eval().(Integer).Int() == ds[i].Eval().(Integer).Int()
		}
	case f.Match(d.Irrational.Flag()):
		return func(i, j int) bool {
			return ds[j].Eval().(Irrational).Float() == ds[i].Eval().(Irrational).Float()
		}
	}
	return nil
}
func newDataFind(ds dataSorter, praed Functional) func(int) bool {
	var f = praed.Eval().Flag()
	var fn func(int) bool
	switch {
	case f.Match(d.Symbolic.Flag()):
		fn = func(i int) bool {
			fmt.Printf("%s %s", ds[i], praed.String())
			return strings.Compare(ds[i].(d.Data).String(), praed.(d.Data).String()) >= 0
		}
	case f.Match(d.Flag.Flag()):
		fn = func(i int) bool {
			return ds[i].(d.Typed).Flag() >= praed.(d.Typed).Flag()
		}
	case f.Match(d.Unsigned.Flag()):
		fn = func(i int) bool {
			return ds[i].(Unsigned).Uint() >= praed.(Unsigned).Uint()
		}
	case f.Match(d.Integer.Flag()):
		fn = func(i int) bool {
			return ds[i].(Integer).Int() >= praed.(Integer).Int()
		}
	case f.Match(d.Irrational.Flag()):
		fn = func(i int) bool {
			return ds[i].(Irrational).Float() >= praed.(Irrational).Float()
		}
	}
	return fn
}

// pair sorter has the methods to search for a pair in-/, and sort slices of
// pairs. pairs will be sorted by the left parameter, since it references the
// accessor (key) in an accessor/value pair.
type pairSorter []Paired

func newPairSorter(p ...Paired) pairSorter                         { return append(pairSorter{}, p...) }
func (a pairSorter) AppendKeyValue(key Functional, val Functional) { a = append(a, NewPair(key, val)) }
func (a pairSorter) Empty() bool {
	if len(a) > 0 {
		for _, p := range a {
			if !ElemEmpty(p) {
				return false
			}
		}
	}
	return true
}
func (p pairSorter) Len() int      { return len(p) }
func (p pairSorter) Swap(i, j int) { p[j], p[i] = p[i], p[j] }
func (p pairSorter) Sort(f d.Type) {
	less := newPraedLess(p, f)
	sort.Slice(p, less)
}
func (p pairSorter) Search(praed Functional) int {
	var idx = sort.Search(len(p), newPraedFind(p, praed))
	// when praedicate is a precedence type encoding bit-flag
	if idx != -1 {
		if praed.Flag().Match(d.Flag.Flag()) {
			if p[idx].Left().Flag() == praed.Flag() {
				return idx
			}
		}
		// otherwise check if key is equal to praedicate
		if idx < len(p) {
			if p[idx].Left().Eval() == praed.Eval() {
				return idx
			}
		}
	}
	return -1
}
func (p pairSorter) Get(praed Functional) Paired {
	idx := p.Search(praed)
	if idx != -1 {
		return p[idx]
	}
	return NewPair(New(d.NilVal{}), New(d.NilVal{}))
}
func (p pairSorter) Range(praed Functional) []Paired {
	var ran = []Paired{}
	idx := p.Search(praed)
	if idx != -1 {
		for pair := p[idx]; pair != nil; {
			ran = append(ran, pair)
		}
	}
	return ran
}

func newPraedLess(accs pairSorter, t d.Type) func(i, j int) bool {
	f := t.Flag()
	switch {
	case f.Match(d.Symbolic.Flag()):
		return func(i, j int) bool {
			chain := accs
			if strings.Compare(
				chain[i].(Paired).Left().Eval().String(),
				chain[j].(Paired).Left().Eval().String(),
			) <= 0 {
				return true
			}
			return false
		}
	case f.Match(d.Flag.Flag()):
		return func(i, j int) bool { // sort by value-, NOT accessor type
			chain := accs
			if chain[i].(Paired).Right().Eval().Flag() <=
				chain[j].(Paired).Right().Eval().Flag() {
				return true
			}
			return false
		}
	case f.Match(d.Unsigned.Flag()):
		return func(i, j int) bool {
			chain := accs
			if chain[i].(Paired).Left().Eval().(Unsigned).Uint() <=
				chain[j].(Paired).Left().Eval().(Unsigned).Uint() {
				return true
			}
			return false
		}
	case f.Match(d.Integer.Flag()):
		return func(i, j int) bool {
			chain := accs
			if chain[i].(Paired).Left().Eval().(Integer).Int() <=
				chain[j].(Paired).Left().Eval().(Integer).Int() {
				return true
			}
			return false
		}
	case f.Match(d.Irrational.Flag()):
		return func(i, j int) bool {
			chain := accs
			if chain[i].(Paired).Left().Eval().(Irrational).Float() <=
				chain[j].(Paired).Left().Eval().(Irrational).Float() {
				return true
			}
			return false
		}
	}
	return nil
}
func newPraedFind(accs pairSorter, praed Functional) func(i int) bool {
	var f = praed.Eval().Flag()
	var fn func(i int) bool
	switch { // parameters are accessor/value pairs to be applyed.
	case f.Match(d.Unsigned.Flag()):
		fn = func(i int) bool {
			return uint(accs[i].(Paired).Left().Eval().(Unsigned).Uint()) >=
				uint(praed.Eval().(Unsigned).Uint())
		}
	case f.Match(d.Integer.Flag()):
		fn = func(i int) bool {
			return int(accs[i].(Paired).Left().Eval().(Integer).Int()) >=
				int(praed.Eval().(Integer).Int())
		}
	case f.Match(d.Irrational.Flag()):
		fn = func(i int) bool {
			return int(accs[i].(Paired).Left().Eval().(Irrational).Float()) >=
				int(praed.Eval().(Irrational).Float())
		}
	case f.Match(d.Symbolic.Flag()):
		fn = func(i int) bool {
			return strings.Compare(
				accs[i].(Paired).Left().Eval().String(),
				praed.Eval().String()) >= 0
		}
	case f.Match(d.Flag.Flag()):
		fn = func(i int) bool {
			return accs[i].(Paired).Right().Eval().(d.BitFlag) >=
				praed.Eval().(d.BitFlag)
		}
	}
	return fn
}
