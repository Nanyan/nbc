package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
)

var (
	ibase, obase int
	help         bool
)

func init() {
	flag.IntVar(&ibase, "i", 0, "input base, must be 0 or a value between 2 and 62. default is determined by input number.")
	flag.IntVar(&obase, "o", 0, "output base, must be 0 or a value between 2 and 62. default is determined by input base.")
	flag.BoolVar(&help, "h", false, "show the usage.")
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

	args := flag.Args()
	if len(args) != 1 {
		exitErr("this app take one and only one argument as the input number")
	}

	inputStr := args[0]

	f, ok := TrySpreadExpOrScienceInteger(inputStr)
	if ok {
		ibase = 10
		inputStr = f
	} else {
		inputStr = determineIbase(inputStr)
	}

	input, ok := new(big.Int).SetString(inputStr, ibase)
	if !ok {
		exitErr("invalid input number and/or input base")
	}

	determineObase()
	out := input.Text(obase)
	fmt.Printf("%s%s\n", oPrefix(), out)

}

func exitErr(msg string) {
	fmt.Println(msg)
	Usage()
	os.Exit(1)
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
func determineIbase(istr string) string {
	if ibase == 0 {
		if len(istr) > 2 {
			if istr[0] == '0' {
				switch istr[1] {
				case 'b', 'B':
					ibase = 2
					return istr[2:]
				case 'x', 'X':
					ibase = 16
					return istr[2:]
				default:
					//prefix '0' as base 8
					ibase = 8
					return istr[1:]
				}
			}
		}
		//if not base 2,8,16,then take it as base 10
		ibase = 10
	} else {
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
	}
	return istr
}

func oPrefix() string {
	switch obase {
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
