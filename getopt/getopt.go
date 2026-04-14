package getopt

// Test example for flag: https://go.dev/src/flag/flag_test.go

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type OptType int

const (
	Flag     OptType = iota // Option with no argument
	String                  // Option with a string argument
	Integer                 // Option with an integer argument
	UInteger                // Option with an unsigned integer argument
)

type Option struct {
	ShortName byte        // Short option character (e.g., 'v')
	LongName  string      // Long option name (e.g., "verbose")
	Type      OptType     // Type of the option
	Target    interface{} // Pointer to variable or function
}

func fatal(format string, a ...interface{}) {
	fmt.Print("")
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(99)
}

func match(opt []Option, s string, lng bool) int {
	var nopt, q, matchlen int
	var p int

	nopt = len(opt)
	if lng {
		p = strings.IndexByte(s, '=')
		if p != -1 {
			matchlen = p
		} else {
			matchlen = len(s)
		}
	}
	for q = 0; q < nopt; q++ {
		if lng {
			if opt[q].LongName == "" {
				continue
			}
			if len(opt[q].LongName) >= matchlen && s[:matchlen] == opt[q].LongName[:matchlen] {
				return q
			}
		} else {
			if opt[q].ShortName == 0 {
				continue
			}
			if s[0] == opt[q].ShortName {
				return q
			}
		}
	}
	return -1
}

func toOptionString(opt *Option) string {
	if len(opt.LongName) > 0 {
		return "--" + opt.LongName
	}
	return "-" + string(opt.ShortName)
}

func needsArgument(opt *Option) bool {
	return opt.Type != Flag
}

func argvRemove(argv *[]string, i int) {
	if i < len(*argv) {
		*argv = append((*argv)[:i], (*argv)[i+1:]...)
	}
}

func execute(opt *Option, arg string, lng bool) {
	switch opt.Type {
	case Flag:
		ptr := opt.Target.(*bool)
		*ptr = true

	case String:
		ptr := opt.Target.(*string)
		*ptr = arg

	case Integer:
		var tmp int64
		var err error
		tmp, err = strconv.ParseInt(arg, 10, 64)
		if err != nil {
			if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
				fatal("number `%s' is out of range\n", arg)
			} else {
				fatal("invalid number `%s'\n", arg)
			}
		}
		assignInteger(opt.Target, tmp, arg)

	case UInteger:
		var tmp uint64
		var err error
		tmp, err = strconv.ParseUint(arg, 10, 64)
		if err != nil {
			if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
				fatal("number `%s' is out of range\n", arg)
			} else {
				fatal("invalid number `%s'\n", arg)
			}
		}
		assignUInteger(opt.Target, tmp, arg)

	default:
		fatal("unhandled option `%s'\n", arg)
	}
}

func assignInteger(target interface{}, n int64, arg string) {
	switch target.(type) {
	case *int:
		if n > math.MaxInt || n < math.MinInt {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*int) = int(n)
	case *int8:
		if n > math.MaxInt8 || n < math.MinInt8 {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*int8) = int8(n)
	case *int16:
		if n > math.MaxInt16 || n < math.MinInt16 {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*int16) = int16(n)
	case *int32:
		if n > math.MaxInt32 || n < math.MinInt32 {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*int32) = int32(n)
	case *int64:
		*target.(*int64) = n
	default:
		fatal("invalid type `%T`, expected signed integer type\n", target)
	}
}

func assignUInteger(target interface{}, n uint64, arg string) {
	switch target.(type) {
	case *uint:
		if n > math.MaxUint {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*uint) = uint(n)
	case *uint8:
		if n > math.MaxUint8 {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*uint8) = uint8(n)
	case *uint16:
		if n > math.MaxUint16 {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*uint16) = uint16(n)
	case *uint32:
		if n > math.MaxUint32 {
			fatal("number `%s' is out of range\n", arg)
		}
		*target.(*uint32) = uint32(n)
	case *uint64:
		*target.(*uint64) = n
	default:
		fatal("invalid type `%T`, expected unsigned integer type\n", target)
	}
}

func Parse(argv *[]string, opt []Option, allowNegNum bool) {
	var ai int     /* argv index. */
	var optarg int /* argv index of option argument, or -1 if none. */
	var mi int     /* Match index in opt. */
	var done bool
	var arg string /* Pointer to argument to an option. */
	var o string   /* pointer to an option string (internal) */

	for ai < len(*argv) {
		/* "--" indicates that the rest of the argv-array does not contain options. */
		if (*argv)[ai] == "--" {
			argvRemove(argv, ai)
			break
		}

		if allowNegNum && (*argv)[ai][0] == '-' && len((*argv)[ai]) > 1 && unicode.IsDigit(rune((*argv)[ai][1])) {
			ai++
			continue
		} else if strings.HasPrefix((*argv)[ai], "--") {
			/* long option */
			/* find matching option */
			mi = match(opt, (*argv)[ai][2:], true)
			if mi < 0 {
				fatal("unrecognized option `%s'\n", (*argv)[ai])
			}

			/* possibly locate the argument to this option. */
			arg = ""
			p := strings.IndexByte((*argv)[ai], '=')
			if p != -1 {
				arg = (*argv)[ai][p+1:]
			}

			/* does this option take an argument? */
			optarg = -1
			if needsArgument(&opt[mi]) {
				/* option needs an argument. find it. */
				if p == -1 {
					optarg = ai + 1
					if optarg == len(*argv) {
						fatal("option `%s' requires an argument\n", toOptionString(&opt[mi]))
					}
					arg = (*argv)[optarg]
				}
			} else {
				if p != -1 {
					fatal("option `%s' doesn't allow an argument\n", toOptionString(&opt[mi]))
				}
			}
			/* perform the action of this option. */
			execute(&opt[mi], arg, true)
			/* remove option and any argument from the argv-array. */
			if optarg >= 0 {
				argvRemove(argv, ai)
			}
			argvRemove(argv, ai)
		} else if (*argv)[ai][0] == '-' {
			/* A dash by itself is not considered an option. */
			if len((*argv)[ai]) == 1 {
				ai++
				continue
			}
			/* Short option(s) following */
			o = (*argv)[ai][1:]
			done = false
			optarg = -1
			for len(o) > 0 && !done {
				/* find matching option */
				mi = match(opt, o, false)
				if mi < 0 {
					fatal("unrecognized option `-%c'\n", o[0])
				}

				/* does this option take an argument? */
				optarg = -1
				arg = ""
				if needsArgument(&opt[mi]) {
					/* option needs an argument. find it. */
					arg = o[1:]
					if arg == "" {
						optarg = ai + 1
						if optarg == len(*argv) {
							fatal("option `%s' requires an argument\n",
								toOptionString(&opt[mi]))
						}
						arg = (*argv)[optarg]
					}
					done = true
				}
				/* perform the action of this option. */
				execute(&opt[mi], arg, false)
				o = o[1:]
			}
			/* remove option and any argument from the argv-array. */
			if optarg >= 0 {
				argvRemove(argv, ai)
			}
			argvRemove(argv, ai)
		} else {
			/* a non-option argument */
			ai++
		}
	}
}
