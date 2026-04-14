# getopt

This is a set of functions for parsing command line options.
Both traditional one-character options, and GNU-style --long-options are supported.

You set up a special structure describing the names and types of the options you want your program to support.
In the structure you also give addresses of variables to fill with value for the various options.
By calling `getopt.Parse`, all options in the incoming `argv` are parsed and _removed from the array_.
What is left are the non-optional arguments to your program.

# Usage

To parse your command line, you need to create and initialize an array of `getopt.Option`.
Each element in the array describes a long and short version of an option and specifies what type of option it is and how to handle it.

The structure fields:

* `ShortName` is the short option name without the leading '-'.
* `LongName` is the long option name without the leading "--".
* `Type` specifies what type of option this is.
* `Target` is a pointer to a location in which to store the value.

The different argument types:

* `Flag` indicates an option that takes no arguments. Target location is expected to be a `bool` variable, which will be set to `true` only if the flag is present,
* `String` expects a string argument.
* `Integer` expects a signed int argument.
* `UInteger` expects an unsigned int argument.

## Notes

* A dash (`-') by itself is not taken as any kind of an option, as several programs use the dash to indicate the special files stdin and stdout.
  It is thus left as a normal argument to the program.

* Two dashes (`--') as an argument is taken to mean that the rest of the arguments should not be scanned for options.
  This simplifies giving names of e.g. files that start with a dash.

* Short (one-character) options accept parameters in two ways, either directly following the option in the same argv-entry, or in the next
  argv-entry:
  ```text
  -sPARAMETER
  -s PARAMETER
  ```

* Long options accept parameters in two ways:
  ```text
  --long-option=PARAMETER
  --long-option PARAMETER
  ```

Ported from C code I wrote back in 1996.
