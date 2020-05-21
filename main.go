package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
)

var (
	ibase, obase           int
	batch, showInput, help bool
	file                   string
)

func init() {
	flag.IntVar(&ibase, "i", 0, "input base, must be 0 or an integer between 2 and 62. default is determined by input number.")
	flag.IntVar(&obase, "o", 0, "output base, must be 0 or an integer between 2 and 62. default is determined by input base.")
	flag.BoolVar(&batch, "b", false, "batch mode, one line for a number, or separate numbers by semicolon.")
	flag.BoolVar(&showInput, "s", false, "show input on the result.")
	flag.BoolVar(&help, "h", false, "show the usage.")
	flag.StringVar(&file, "f", "", "read inputs from file. this will ignores the -b option.")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	flag.Parse()
	if help {
		Usage()
		return
	}
	if !checkBase(ibase) {
		exitErr("invalid input base")
	}
	if !checkBase(obase) {
		exitErr("invalid output base")
	}

	var inputs []*Input
	var origin string
	var err error
	if file != "" {
		origin, inputs, err = InputsFromFile(file)
	} else {
		origin, inputs, err = InputsFromArgsOrPipeline(flag.Args(), batch)
	}
	if err != nil {
		exitErr(fmt.Sprintf("%q\terr=%v ", origin, err))
	}

	//determine ibase by first input if ibase not set
	determineIbase(inputs[0].Trimmed)
	determineObase()

	//handle each input
	exitCode := 0
	for _, item := range inputs {
		if ibase == 10 {
			f, ok := TrySpreadExpOrScienceInteger(item.Trimmed)
			if ok {
				item.Trimmed = f
			}
		} else {
			item.Trimmed = trimPrefix(item.Trimmed)
		}
		err := convert(item.Origin, item.Trimmed)
		if err != nil {
			fmt.Println(err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func exitErr(msg string) {
	fmt.Println(msg)
	fmt.Println()
	fmt.Println("Run 'nbc -h' for the help message.")
	os.Exit(1)
}

func convert(origin, s string) error {
	input, ok := new(big.Int).SetString(s, ibase)
	if !ok {
		return fmt.Errorf("invalid input number (and/or input base): %q , after trim: %s", origin, s)
	}

	out := input.Text(obase)
	if showInput {
		fmt.Printf("%s%s : %s%s\n", prefix(ibase), s, prefix(obase), out)
	} else {
		fmt.Printf("%s%s\n", prefix(obase), out)
	}
	return nil
}

func checkBase(base int) bool {
	if base == 0 || (base > 1 && base <= big.MaxBase) {
		return true
	}
	return false
}

//determineObase determine the output base when flag -o is not set.
//when ibase is 10, the default obase is 16;
//otherwise the default obase is 10
func determineObase() {
	if obase == 0 {
		switch ibase {
		case 10:
			obase = 16
		default:
			obase = 10
		}
	}
}

//determineIbase only determine whether the inputbase is 2 or 16
// if flag -i not set.
// remove prefix if determined, otherwise return the origin string.
func determineIbase(istr string) {
	if ibase == 0 {
		if len(istr) > 2 {
			if istr[0] == '0' {
				switch istr[1] {
				case 'b', 'B':
					ibase = 2
				case 'x', 'X':
					ibase = 16
				default:
					//prefix '0' as base 8
					ibase = 8
				}
				return
			}
		}
		//if not base 2,8,16,then take it as base 10
		ibase = 10
	}
}

func trimPrefix(istr string) string {
	//remove prefix if ibase is set to 2 or 16
	if ibase == 2 || ibase == 8 || ibase == 16 {
		if istr[0] == '0' {
			switch istr[1] {
			case 'b', 'B', 'x', 'X':
				return istr[2:]
			default:
				return istr[1:]
			}
		}
	}
	return istr
}

func prefix(base int) string {
	switch base {
	case 2:
		return "0b"
	case 8:
		return "0"
	case 16:
		return "0x"
	default:
		return ""
	}
}
