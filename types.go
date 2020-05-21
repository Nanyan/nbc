package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Input struct {
	Origin, Trimmed string
}

func InputsFromFile(file string) (string, []*Input, error) {
	origin := fmt.Sprintf("from file %q", file)
	f, err := os.Open(file)
	if err != nil {
		return origin, nil, err
	}
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0, 8)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	r := NewInputs(lines)
	if len(r) == 0 {
		return origin, nil, errors.New("no valid inputs")
	}
	return origin, r, nil
}

func InputsFromArgsOrPipeline(args []string, batch bool) (string, []*Input, error) {
	var origin string
	switch len(args) {
	case 0:
		//check pipeline input
		var ok bool
		origin, ok = checkGetPipelineInput()
		if !ok {
			return "", nil, errors.New("no input from pipeline or args")
		}
	case 1:
		origin = args[0]
	default:
		return "", nil, errors.New("too many args")
	}
	var inputs []string
	if batch {
		first := strings.Split(origin, "\n")
		for _, str := range first {
			inputs = append(inputs, strings.Split(str, ";")...)
		}
	} else {
		inputs = []string{origin}
	}
	r := NewInputs(inputs)
	if len(r) == 0 {
		return origin, nil, errors.New("no valid inputs")
	}
	return origin, r, nil
}

func checkGetPipelineInput() (string, bool) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", false
	}
	if (info.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe ||
		((info.Mode()&os.ModeCharDevice) == os.ModeCharDevice && info.Size() > 0) {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", false
		}
		return string(bytes), true
	}
	return "", false
}

func NewInputs(inputs []string) []*Input {
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
