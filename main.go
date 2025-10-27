package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
)

var (
	ibase, obase           int
	ibaseMode, obaseMode   string
	batch, showInput, help bool
	file                   string
)

var (
	ibaseStr, obaseStr string // 支持 b58/b64
)

func init() {
	flag.StringVar(&ibaseStr, "i", "", "input base, must be 0, an integer between 2 and 62, or b58/b64.")
	flag.StringVar(&obaseStr, "o", "", "output base, must be 0, an integer between 2 and 62, or b58/b64.")
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
	// 解析 ibaseStr/obaseStr 支持 b58/b64
	ibase, ibaseMode = parseBase(ibaseStr)
	obase, obaseMode = parseBase(obaseStr)
	if ibaseMode == "invalid" {
		exitErr("invalid input base " + ibaseStr)
	}
	if obaseMode == "invalid" {
		exitErr("invalid output base " + obaseStr)
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
	if ibaseMode == "int" && ibase == 0 {
		determineIbase(inputs[0].Trimmed)
	}
	if obaseMode == "int" && obase == 0 {
		determineObase(ibase, ibaseMode)
	}

	exitCode := 0
	for _, item := range inputs {
		var err error
		switch {
		case ibaseMode == "b58":
			// base58 解码
			decoded, derr := decodeBase58(item.Trimmed)
			if derr != nil {
				fmt.Println(derr)
				exitCode = 1
				continue
			}
			err = convertBytes(item.Origin, decoded, obase, obaseMode)
		case ibaseMode == "b64":
			// base64 解码
			decoded, derr := decodeBase64(item.Trimmed)
			if derr != nil {
				fmt.Println(derr)
				exitCode = 1
				continue
			}
			err = convertBytes(item.Origin, decoded, obase, obaseMode)
		case obaseMode == "b58":
			// base58 编码
			var b []byte
			if ibase == 16 || ibase == 0 {
				// 默认输入为16进制
				b, err = hexToBytes(item.Trimmed)
			} else {
				// 其他进制转为10进制再转字节
				n, nerr := new(big.Int).SetString(item.Trimmed, ibase)
				if !nerr {
					err = fmt.Errorf("invalid input number: %q", item.Trimmed)
				} else {
					b = n.Bytes()
				}
			}
			if err == nil {
				out := encodeBase58(b)
				if showInput {
					fmt.Printf("%s : %s\n", item.Trimmed, out)
				} else {
					fmt.Printf("%s\n", out)
				}
			}
		case obaseMode == "b64":
			// base64 编码
			var b []byte
			if ibase == 16 || ibase == 0 {
				b, err = hexToBytes(item.Trimmed)
			} else {
				n, nerr := new(big.Int).SetString(item.Trimmed, ibase)
				if !nerr {
					err = fmt.Errorf("invalid input number: %q", item.Trimmed)
				} else {
					b = n.Bytes()
				}
			}
			if err == nil {
				out := encodeBase64(b)
				if showInput {
					fmt.Printf("%s : %s\n", item.Trimmed, out)
				} else {
					fmt.Printf("%s\n", out)
				}
			}
		default:
			// 普通进制转换
			if ibase == 10 {
				f, ok := TrySpreadExpOrScienceInteger(item.Trimmed)
				if ok {
					item.Trimmed = f
				}
			} else {
				item.Trimmed = trimPrefix(item.Trimmed)
			}
			err = convert(item.Origin, item.Trimmed)
		}
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

// determineObase determine the output base when flag -o is not set.
// when ibase is 10, the default obase is 16;
// otherwise the default obase is 10
func determineObase(ib int, ibaseMode string) {
	if obase == 0 {
		switch ibaseMode {
		case "int":
			switch ib {
			case 10:
				obase = 16
			default:
				obase = 10
			}
		case "b58", "b64":
			obase = 16
		default:
			obase = 10
		}
	}
}

// determineIbase only determine whether the inputbase is 2 or 16
// if flag -i not set.
// remove prefix if determined, otherwise return the origin string.
func determineIbase(istr string) {
	if ibase == 0 {
		if obaseMode == "b58" || obaseMode == "b64" {
			// when output is base58 or base64, input is assumed to be hex if not specified.
			ibase = 16
			return
		}
		//check prefix
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

// parseBase 支持 b58/b64/int
func parseBase(s string) (int, string) {
	switch s {
	case "b58", "base58":
		return 0, "b58"
	case "b64", "base64":
		return 0, "b64"
	case "", "0":
		return 0, "int"
	default:
		var n int
		_, err := fmt.Sscanf(s, "%d", &n)
		if err == nil && checkBase(n) {
			return n, "int"
		}
		return 0, "invalid"
	}
}

// base58 编码表
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func encodeBase58(b []byte) string {
	x := new(big.Int).SetBytes(b)
	var out []byte
	base := big.NewInt(58)
	zero := big.NewInt(0)
	for x.Cmp(zero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, base, mod)
		out = append([]byte{b58Alphabet[mod.Int64()]}, out...)
	}
	// 前导零
	for _, v := range b {
		if v == 0 {
			out = append([]byte{b58Alphabet[0]}, out...)
		} else {
			break
		}
	}
	return string(out)
}

func decodeBase58(s string) ([]byte, error) {
	x := big.NewInt(0)
	base := big.NewInt(58)
	for i := 0; i < len(s); i++ {
		idx := bytesIndex(b58Alphabet, s[i])
		if idx < 0 {
			return nil, fmt.Errorf("invalid base58 char: %c", s[i])
		}
		x.Mul(x, base)
		x.Add(x, big.NewInt(int64(idx)))
	}
	// 处理前导零
	n := 0
	for n < len(s) && s[n] == b58Alphabet[0] {
		n++
	}
	out := x.Bytes()
	if n > 0 {
		out = append(make([]byte, n), out...)
	}
	return out, nil
}

func bytesIndex(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func hexToBytes(s string) ([]byte, error) {
	if len(s) > 1 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// bytes 转换为指定进制输出
func convertBytes(origin string, b []byte, obase int, obaseMode string) error {
	switch obaseMode {
	case "int":
		n := new(big.Int).SetBytes(b)
		out := n.Text(obase)
		if obase == 16 {
			if len(out)%2 == 1 {
				out = "0" + out
			}
			out = "0x" + out
		}
		if showInput {
			fmt.Printf("%s : %s\n", origin, out)
		} else {
			fmt.Printf("%s\n", out)
		}
		return nil
	default:
		return fmt.Errorf("convertBytes: unsupported output base mode %s", obaseMode)
	}
}
