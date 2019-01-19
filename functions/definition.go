package functions

import (
	"strconv"

	d "github.com/JoergReinhardt/godeep/data"
)

type function praedicates

func (f function) get(acc Data) Paired { return praedicates(f).Get(acc) }
func newFuncDef(
	uid int, // unique id for this type
	name string,
	prec d.BitFlag,
	kind Kind,
	fixity Property, // PostFix|InFix|PreFix
	lazy Property, // Eager|Lazy
	bound Property, // Left_Bound|Right_Bound
	mutable Property, // Mutable|Imutable
	fnc Function, //  Call(...Data) Data
	retType Flag, // type of return value
	argSig ...Data, // either []Flag, or []Paired
) FncDef {
	var argsAcc = NamedArgs
	var flags = []Flag{}
	var accs = []Data{}
	// infer argument accessor property (access via named-, or positional arguments)
	if len(argSig) > 0 {
		for _, arg := range argSig { // ‥.eiher []Paired ⇒ named arguments
			if pair, ok := arg.(Paired); ok {
				flags = append(flags, pair.Right().(Flag))
				accs = append(accs, pair.Left())
			} else { // ‥. or []Flag ⇒ positional arguments
				if fl, ok := arg.(Flag); ok {
					argsAcc = Positional
					flags = append(flags, fl)
				}
			}
		}
	}
	var flagSet = newFlagSet(flags...)
	var accSet = newArguments(accs...)
	return function(newPraedicates(
		newPair(d.StrVal("uid"), d.UintVal(uid)),
		newPair(d.StrVal("name"), d.StrVal(name)),
		newPair(d.StrVal("prec"), prec),
		newPair(d.StrVal("kind"), kind),
		newPair(d.StrVal("fixity"), fixity),
		newPair(d.StrVal("lazy"), lazy),
		newPair(d.StrVal("bound"), bound),
		newPair(d.StrVal("mutable"), mutable),
		newPair(d.StrVal("rettype"), retType),
		newPair(d.StrVal("fnc"), fnc),
		// inferred propertys
		newPair(d.StrVal("arity"), Arity(len(argSig))),
		newPair(d.StrVal("argsacc"), argsAcc),
		newPair(d.StrVal("argtypes"), flagSet),
		newPair(d.StrVal("accs"), accSet),
	).(praedicates))
}
func (f function) String() string {
	var rows = [][]d.Data{}
	pairs, _ := f()
	props := pairs[0:12]
	for i, pair := range props {
		rows = append(rows, []d.Data{})
		rows[i] = append(rows[i],
			d.StrVal(strconv.Itoa(i)),
			d.StrVal(pair.Left().String()),
			d.StrVal(pair.Right().String()),
		)
	}
	rows = append(rows, []d.Data{d.IntVal(12), d.StrVal("_____________"), d.StrVal("_____________")})

	var row12a = [][]d.Data{}
	flags := f.ArgTypes()
	for i, flag := range flags {
		row12a = append(row12a, []d.Data{})
		row12a[i] = append(row12a[i],
			d.StrVal(strconv.Itoa(i)),
			d.StrVal(flag.String()),
		)
	}

	var row12b = [][]d.Data{}
	args := f.Accs()
	for i, arg := range args {
		row12b = append(row12b, []d.Data{})
		row12b[i] = append(row12b[i],
			d.StrVal(strconv.Itoa(i)),
			d.StrVal(arg.Flag().String()),
			d.StrVal(arg.(Argumented).Data().String()),
		)
	}

	for i, _ := range flags {
		if f.AccessType() == NamedArgs {
			rows = append(rows, []d.Data{d.IntVal(len(rows)), row12b[i][1], row12b[i][2]})
		} else {
			rows = append(rows, []d.Data{d.IntVal(len(rows)), d.IntVal(i), flags[i]})
		}
	}
	return d.StringChainTable(rows...)
}
func (f function) Type() Flag           { return newFlag(f.UID(), f.Kind(), f.Prec()) }
func (f function) Flag() d.BitFlag      { return f.Prec() }
func (f function) UID() int             { return f.get(d.StrVal("uid")).Right().(Integer).Int() }
func (f function) Name() string         { return f.get(d.StrVal("name")).Right().(Symbolic).String() }
func (f function) Prec() d.BitFlag      { return f.get(d.StrVal("prec")).Right().(d.BitFlag) }
func (f function) Kind() Kind           { return f.get(d.StrVal("kind")).Right().(Kind) }
func (f function) Arity() Arity         { return f.get(d.StrVal("arity")).Right().(Arity) }
func (f function) Fix() Property        { return f.get(d.StrVal("fixity")).Right().(Property) }
func (f function) Lazy() Property       { return f.get(d.StrVal("lazy")).Right().(Property) }
func (f function) Bound() Property      { return f.get(d.StrVal("bound")).Right().(Property) }
func (f function) Mutable() Property    { return f.get(d.StrVal("mutable")).Right().(Property) }
func (f function) AccessType() Property { return f.get(d.StrVal("argsacc")).Right().(Property) }
func (f function) ArgTypes() []Flag     { return f.get(d.StrVal("argtypes")).Right().(FlagSet)() }
func (f function) Accs() []Argumented   { return f.get(d.StrVal("accs")).Right().(arguments).Args() }
func (f function) RetType() Flag        { return f.get(d.StrVal("rettype")).Right().(Flag) }
func (f function) Fnc() Function        { return f.get(d.StrVal("fnc")).Right().(Function) }
