# nbc

A very very simple and easy to use number base converter written in go.

## Install

```bash
go install github.com/nanyan/nbc@latest
```

## Usage

It just converts an input num string to another base num string as the output. For more details, see the help message.

```bash
~ » nbc -h
nbc -- A very very simple and easy to use number base converter.
Version: v1.0.0

Usage: nbc [options] numString

Options:
  -b         batch mode, one line for a number, or separate numbers by semicolon.
  -f string  read inputs from file. this will ignores the -b option.
  -h         show the usage.
  -i int     input base, must be 0 or an integer between 2 and 62. default is determined by input number.
  -o int     output base, must be 0 or an integer between 2 and 62. default is determined by input base.
  -s         show input on the result.

Description:
  The base argument must be 0 or a value between 2 and 62. If the base is 0, the string prefix determines the actual conversion base. A prefix of "0x" or "0X" selects base 16; the "0" prefix selects base 8, and a "0b" or "0B" prefix selects base 2. Otherwise the selected base is 10.
If the output base is not set(or set to 0), then the actual output base will be determined by the input base, when -i is 10, then default output base is 16; otherwise output base is default 10.
For bases <= 36, lower and upper case letters are considered the same: The letters 'a' to 'z' and 'A' to 'Z' represent digit values 10 to 35. For bases > 36, the upper case letters 'A' to 'Z' represent the digit values 36 to 61.

The input number also support scientific notation base 10, such as 1e18 or 10^18.

Example:
  nbc 0xff             //convert hex number 0xff to decimal
  nbc 10^18            //convert a big decimal to hexadecimal
  nbc 1e18             //also convert a big decimal to hexadecimal
  nbc -o 2 16          //convert decimal 16 to binary string
  nbc -b -s "100;200"  //convert batch decimal to hex, showing input in the result, as 100 : 0x64
  echo 0x10|nbc        //use pipeline input
  nbc -f input.txt     //read inputs from file
```
