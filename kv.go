//MIT License
//
//Copyright (c) 2019 Mike McDermott
//See LICENSE file for details
//
// Simple utility for acting like a kv store for the command line

package main

import (
	"bufio"
	"fmt"
	"github.com/Wessie/appdirs"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var G_VERBOSE = new(bool)

// We basically only care about two levels of logging, at least right now
func vprint(a ...interface{}) () {
	if *G_VERBOSE {
		fmt.Println(a...)
	}
}

type keyError struct {
	key string
}

func (e *keyError) Error() string {
	return fmt.Sprintf("KeyError: %s", e.key)
}

type ProgramMode int

const (
	ModeHelp ProgramMode = iota
	ModeSimpleGet
	ModeSimpleSet
	ModePipeSet
	ModeList
	ModeDelete
)

type FlagsMode int

const (
	FlagNone FlagsMode = iota
	FlagEndOfOpts
	FlagList
	FlagDelete
	FlagPop
	FlagVerbose
)

var flagsMap = map[string]FlagsMode{
	"--":        FlagEndOfOpts,
	"-l":        FlagList,
	"--list":    FlagList,
	"-d":        FlagDelete,
	"--delete":  FlagDelete,
	"-v":        FlagVerbose,
	"--verbose": FlagVerbose,
}

type StdInContainer struct {
	stdin     string
	has_stdin bool
}

type ParsedArgs struct {
	flags    []string // anything starting with '-' including just '-' and '--'
	args     []string // things which are arguments without dashes
	key_val  string   // an argument which may be 'key' or 'key=val'
	key      string
	val      string
	mode     ProgramMode // mode: getting/setting etc
	err_code int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func anyStringInSlice(options []string, list []string) bool {
	for _, b := range list {
		for _, a := range options {
			if b == a {
				return true
			}
		}

	}
	return false
}

// parse flags to a list of enums
func checkFlags(list []string) (map[FlagsMode]bool, error) {
	// collection of bool flags that are on/off
	var bOptions = map[FlagsMode]bool{
		FlagList:    false,
		FlagDelete:  false,
		FlagVerbose: false,
	}
	for _, element := range list {
		if flagEnum, ok := flagsMap[element]; ok {
			if _, ok2 := bOptions[flagEnum]; ok2 {
				bOptions[flagEnum] = true
			}
		} else {
			return nil, fmt.Errorf("Unrecognized flag: %s", element)
		}

	}
	return bOptions, nil
}

func hashx(s string) string {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	check(err)
	return strconv.FormatUint(h.Sum64(), 16)
}

func read_value(lookup_file string, key string) (string, error) {
	dat, err := ioutil.ReadFile(lookup_file)
	if err != nil {
		return "", &keyError{key: key}
	}
	value := strings.Split(string(dat), "=")[1]
	return value, nil
}

func store_value(lookup_file string, key string, value string) {
	data := []byte(key + "=" + value)
	err := ioutil.WriteFile(lookup_file, data, 0644)
	check(err)
}

func pop_value(lookup_file string, key string) (string, error) {
	val, err := read_value(lookup_file, key)
	if err != nil {
		return "", err
	}
	err = os.Remove(lookup_file)
	return val, err

}

func delete_value(lookup_file string, key string) (error) {
	vprint("Deleting [%s] (%s)", key, lookup_file)
	_, err := pop_value(lookup_file, key)
	return err
}

func get_stdin() (StdInContainer, error) {
	info, err := os.Stdin.Stat()
	check(err)
	out_struct := StdInContainer{has_stdin: false}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		//fmt.Println("stdin is from a terminal")
		return out_struct, nil
	}

	// data is being piped to stdin
	//fmt.Println("data is being piped to stdin")

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	out_struct.stdin = string(output)
	out_struct.has_stdin = true
	return out_struct, nil
}

// This should always return a key, or fail/go to help
func parseCLI() (ParsedArgs) {
	parsed_args := ParsedArgs{mode: ModeHelp, err_code: 0}

	// separate flags from args
	endOfOpts := false
	for _, arg := range os.Args[1:] {
		if arg[0:2] == "--" {
			endOfOpts = true
		} else if arg[0] == '-' {
			if endOfOpts {
				log.Fatal("Cannot have flags after end-of-options (--)")
			}
			parsed_args.flags = append(parsed_args.flags, arg)
		} else {
			parsed_args.args = append(parsed_args.args, arg)
		}
	}

	// deal with flags
	if len(parsed_args.flags) > 0 {
		bOptions, err := checkFlags(parsed_args.flags)
		if err != nil {
			log.Fatal(err)
		}

		if val, ok := bOptions[FlagList]; ok && val {
			parsed_args.mode = ModeList
		}
		if val, ok := bOptions[FlagDelete]; ok && val {
			parsed_args.mode = ModeDelete
		}
		if val, ok := bOptions[FlagVerbose]; ok && val {
			*G_VERBOSE = true
			vprint("Verbose is on")
		}

	}

	// deal with stdin, if present
	stdin_struct, err := get_stdin()
	check(err)
	if (len(os.Args) == 1) && stdin_struct.has_stdin {
		fmt.Println("<!> Error:  Piping from stdin, but no key provided")
		parsed_args.err_code = 1
		return parsed_args
	}

	if (len(os.Args) == 1) && !stdin_struct.has_stdin {
		fmt.Println("<!> Error:  must provide one or more arguments")
		parsed_args.err_code = 1
		return parsed_args
	}

	if stdin_struct.has_stdin {
		panic("/\\--/\\ Cannot handle pipe yet (under construction)")

	}

	// deal with args now
	if len(parsed_args.args) > 1 {
		panic("/\\--/\\ Too many arguments (under construction)")
	}

	if len(parsed_args.args) == 1 {
		key_val := strings.Split(parsed_args.args[0], "=")

		if len(key_val) > 2 {
			log.Fatal("Too many equals")
			os.Exit(1)

		} else if (len(key_val) == 2) {
			parsed_args.key = strings.TrimSpace(key_val[0])
			parsed_args.val = strings.TrimSpace(key_val[1])
			parsed_args.mode = ModeSimpleSet
		} else {
			parsed_args.key = strings.TrimSpace(key_val[0])
			if parsed_args.mode == 0 {
				parsed_args.mode = ModeSimpleGet
			}
		}
	}

	// default - show help
	return parsed_args
}
func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func list_all(lookup_path string) {
	var files []string

	err := filepath.Walk(lookup_path, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fi, err := os.Stat(file)
		check(err)
		if fi.Mode().IsRegular() {
			dat, err := ioutil.ReadFile(file)
			check(err)
			fmt.Println(string(dat))
		}
	}
}

func show_help() {
	fmt.Println(` Usage: kv [OPTIONS] KEY[=VALUE]
    kv is a simple utility for getting and setting key-value pairs. 
    Getting: kv KEY
    Setting: kv KEY=VALUE`)
}

func main() {
	*G_VERBOSE = false

	parsed_args := parseCLI()
	// todo: add flag args

	if parsed_args.mode == ModeHelp {
		show_help()
		os.Exit(parsed_args.err_code)
	}

	lookup_key := hashx(parsed_args.key)
	lookup_path := appdirs.UserDataDir("xac", "", "", false)
	lookup_file := lookup_path + "/" + lookup_key
	if *G_VERBOSE {
		fmt.Println("path: ", lookup_path)
		fmt.Println("lfile: ", lookup_file)
		fmt.Println("flags: ", parsed_args.flags)
		fmt.Println("args: ", parsed_args.args)
		fmt.Println("mode: ", parsed_args.mode)

	}

	if _, err := os.Stat(lookup_path); os.IsNotExist(err) {
		_ = os.MkdirAll(lookup_path, os.ModePerm)
	}

	switch parsed_args.mode {
	case ModeSimpleSet:
		// Store key=val in file so we can query available keys
		store_value(lookup_file, parsed_args.key, parsed_args.val)
	case ModeSimpleGet:
		val, err := read_value(lookup_file, parsed_args.key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
	case ModeList:
		list_all(lookup_path)
	case ModeDelete:
		err := delete_value(lookup_file, parsed_args.key)
		if err != nil {
			log.Fatal(err)
		}
	default:
		show_help()
	}

}
