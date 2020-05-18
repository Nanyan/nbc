package main

import "strings"

type Input struct {
	Origin, Trimmed string
}

func NewInput(s string, batch bool) []*Input {
	var inputs []string
	if batch {
		first := strings.Split(s, "\n")
		for _, str := range first {
			inputs = append(inputs, strings.Split(str, ";")...)
		}
	} else {
		inputs = []string{s}
	}
	r := make([]*Input, 0, len(inputs))
	for _, o := range inputs {
		t := trimInput(o)
		if t != "" {
			r = append(r, &Input{
				Origin:  o,
				Trimmed: t,
			})
		}
	}
	return r
}

func trimInput(s string) string {
	return strings.Trim(s, "'\"\r\n ")
}
