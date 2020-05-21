package main

import (
	"flag"
	"os"
	"text/tabwriter"
	"text/template"
)

var helpTemplate = `nbc -- A very very simple and easy to use number base converter.
Version: {{.Version}}

Usage: {{.Usage}}{{if .Flags}}

Options:  {{range .Flags}}
  {{"-"}}{{.Name}}{{" "}}{{ft .}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{if .Description}}

Description:
  {{.Description}}{{end}}{{if .Example}}

Example:
  {{.Example}}{{end}}
`

func Usage() {
	flags := make([]*flag.Flag, 0, 4)
	flag.VisitAll(func(i *flag.Flag) {
		flags = append(flags, i)
	})

	type helpData struct {
		Version     string
		Usage       string
		Description string
		Flags       []*flag.Flag
		Example     string
	}

	funcMap := template.FuncMap{
		"ft": getType,
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(helpTemplate))

	data := helpData{
		Version: "v1.0.0",
		Usage:   "nbc [options] numString",
		Description: `The base argument must be 0 or a value between 2 and 62. If the base is 0, the string prefix determines the actual conversion base. A prefix of "0x" or "0X" selects base 16; the "0" prefix selects base 8, and a "0b" or "0B" prefix selects base 2. Otherwise the selected base is 10.
If the output base is not set(or set to 0), then the actual output base will be determined by the input base, when -i is 10, then default output base is 16; otherwise output base is default 10.
For bases <= 36, lower and upper case letters are considered the same: The letters 'a' to 'z' and 'A' to 'Z' represent digit values 10 to 35. For bases > 36, the upper case letters 'A' to 'Z' represent the digit values 36 to 61.

The input number also support scientific notation base 10, such as 1e18 or 10^18.`,
		Flags: flags,
		Example: "nbc 0xff    	//convert hex number 0xff to decimal\n  nbc 10^18   	//convert a big decimal to hexadecimal\n  nbc 1e18    	//also convert a big decimal to hexadecimal\n  nbc -o 2 16 	//convert decimal 16 to binary string\n  nbc -b -s \"100;200\"	//convert batch decimal to hex, showing input in the result, as 100 : 0x64\n  echo 0x10|nbc 	//use pipeline input\n  nbc -f input.txt	//read inputs from file\n",
	}
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func getType(f *flag.Flag) string {
	t, _ := flag.UnquoteUsage(f)
	return t
}
