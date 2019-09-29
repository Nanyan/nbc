package main

import "strconv"

func TrySpreadExpOrScienceInteger(input string) (full string, isOk bool) {
	//var isSci bool
	notationIndex := -1
	dotIndex := -1
	bs := []byte(input)
	//10 base exp
	if len(bs) > 3 {
		if bs[0] == '1' && bs[1] == '0' && bs[2] == '^' {
			return spreadExpBase10(bs)
		}
	}

	for i, ch := range bs {
		switch ch {
		case '.':
			if dotIndex == -1 && i == 1 {
				dotIndex = i
			} else {
				return
			}
		case 'e', 'E':
			if notationIndex == -1 {
				notationIndex = i
			} else {
				return
			}
		default:
			if ch < '0' || ch > '9' {
				return
			}
		}
	}
	if notationIndex < 1 {
		return
	}
	exp, err := strconv.Atoi(string(bs[notationIndex+1:]))
	if err != nil {
		return
	}
	var r, main []byte
	if dotIndex > 0 {
		mantissa := bs[dotIndex+1 : notationIndex]
		if len(mantissa) > exp {
			//it's probably not an integer, don't handle detail
			return
		}
		main = make([]byte, dotIndex+len(mantissa))
		copy(main, bs[:dotIndex])
		copy(main[dotIndex:], mantissa)
		exp -= len(mantissa)
	} else {
		//no dot, then notation should be in index 1
		if notationIndex != 1 {
			return
		}
		main = bs[:notationIndex]
	}
	lm := len(main)
	r = make([]byte, lm+exp)
	copy(r[:lm], main)
	for i := 0; i < exp; i++ {
		r[lm+i] = '0'
	}
	return string(r), true
}

func spreadExpBase10(bs []byte) (full string, isOk bool) {
	exp, err := strconv.Atoi(string(bs[3:]))
	if err != nil {
		return
	}
	if exp > 0 {
		r := make([]byte, 1+exp)
		r[0] = '1'
		for i := 1; i <= exp; i++ {
			r[i] = '0'
		}
		return string(r), true
	} else if exp < 0 {
		return "", false
	}
	return "1", true
}
