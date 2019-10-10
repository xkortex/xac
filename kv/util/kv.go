//MIT License
//
//Copyright (c) 2019 Mike McDermott
//See LICENSE file for details
//
// Simple utility for acting like a kv store for the command line

package util

import (
	"bufio"
	"fmt"
	"github.com/Wessie/appdirs"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
	FlagHelp
	FlagList
	FlagDelete
	FlagNamespace
	FlagPop
	FlagVerbose
)

var flagsMap = map[string]FlagsMode{
	"--":        FlagEndOfOpts,
	"-h":        FlagHelp,
	"--help":    FlagHelp,
	"-?":        FlagHelp,
	"-l":        FlagList,
	"--list":    FlagList,
	"-d":        FlagDelete,
	"--delete":  FlagDelete,
	"-n":        FlagNamespace,
	"--name":    FlagNamespace,
	"-v":        FlagVerbose,
	"--verbose": FlagVerbose,
}

type StdInContainer struct {
	Stdin     string
	Has_stdin bool
}

type ParsedArgs struct {
	flags    []string // anything starting with '-' including just '-' and '--'
	args     []string // things which are arguments without dashes
	key_val  string   // an argument which may be 'key' or 'key=val'
	key      string
	val      string
	name     string      // namespace
	mode     ProgramMode // mode: getting/setting etc
	err_code int
}

func Panic_if(e error) {
	if e != nil {
		panic(e)
	}
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

func Hashx(s string) string {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	Panic_if(err)
	return strconv.FormatUint(h.Sum64(), 16)
}

func Read_value(lookup_file string, key string) (string, error) {
	dat, err := ioutil.ReadFile(lookup_file)
	if err != nil {
		return "", &keyError{key: key}
	}
	value := string(dat)
	return value, nil
}

func Store_value(lookup_file string, key string, value string) {
	data := []byte(value)
	lookup_basepath := filepath.Dir(lookup_file)
	if _, err := os.Stat(lookup_basepath); os.IsNotExist(err) {
		_ = os.MkdirAll(lookup_basepath, os.ModePerm)
	}
	err := ioutil.WriteFile(lookup_file, data, 0644)
	Panic_if(err)
}

func Pop_value(lookup_file string, key string) (string, error) {
	val, err := Read_value(lookup_file, key)
	if err != nil {
		return "", err
	}
	err = os.Remove(lookup_file)
	return val, err

}

func Delete_value(lookup_file string, key string) (error) {
	Vprint("Deleting [%s] (%s)", key, lookup_file)
	_, err := Pop_value(lookup_file, key)
	return err
}

func Get_stdin() (StdInContainer, error) {
	info, err := os.Stdin.Stat()
	Panic_if(err)
	out_struct := StdInContainer{Has_stdin: false}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		//fmt.Println("Stdin is from a terminal")
		return out_struct, nil
	}

	// data is being piped to Stdin
	//fmt.Println("data is being piped to Stdin")

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	out_struct.Stdin = string(output)
	out_struct.Has_stdin = true
	return out_struct, nil
}

func GetLookupPath(namespace string, key string) string {
	lookup_path := appdirs.UserDataDir("kv", "", "", false)
	if len(namespace) > 0 {
		namespace = filepath.Join(".ns", namespace)
	}
	if len(key) == 0 {
		return filepath.Join(lookup_path, namespace)
	}

	lookup_key := url.PathEscape(key)
	return filepath.Join(lookup_path, namespace, lookup_key)
}

// This should always return a key, or fail/go to help
func parseCLI() (ParsedArgs) {
	parsed_args := ParsedArgs{mode: ModeHelp, err_code: 0}

	// separate flags from args - this is rapidly getting out of hand
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
			Vprint("Verbose is on")
		}

	}

	// deal with Stdin, if present
	stdin_struct, err := Get_stdin()
	Panic_if(err)
	if (len(os.Args) == 1) && stdin_struct.Has_stdin {
		fmt.Println("<!> Error:  Piping from Stdin, but no key provided")
		parsed_args.err_code = 1
		return parsed_args
	}

	if (len(os.Args) == 1) && !stdin_struct.Has_stdin {
		fmt.Println("<!> Error:  must provide one or more arguments")
		parsed_args.err_code = 1
		return parsed_args
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

		}

		if stdin_struct.Has_stdin {
			if (len(key_val) == 2) {
				log.Fatal("Cannot use `key=val` with pipe in")
			} else {
				parsed_args.key = strings.TrimSpace(key_val[0])
				if parsed_args.mode == 0 {
					parsed_args.mode = ModeSimpleSet
					parsed_args.val = strings.TrimSpace(stdin_struct.Stdin)
				} else {
					log.Fatal("Flags incompatible with pipe in:", parsed_args.flags)
				}
			}
		} else {
			if (len(key_val) == 2) {
				parsed_args.key = strings.TrimSpace(key_val[0])
				parsed_args.val = strings.TrimSpace(key_val[1])
				parsed_args.mode = ModeSimpleSet
			} else {
				parsed_args.key = strings.TrimSpace(key_val[0])
				// if any existing flags are set, such as delete, use that mode
				if parsed_args.mode == 0 {
					parsed_args.mode = ModeSimpleGet
				}
			}
		}

	}

	// default - show help
	return parsed_args
}
func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func List_all(lookup_path string) {
	var files []string

	err := filepath.Walk(lookup_path, Visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fi, err := os.Stat(file)
		Panic_if(err)
		if fi.Mode().IsRegular() {
			dat, err := ioutil.ReadFile(file)
			Panic_if(err)
			// todo: deal with namespaces and formatting and a bunch of issues
			key, _ := url.PathUnescape(filepath.Base(file))
			rel, err := filepath.Rel(lookup_path, file)
			Panic_if(err)
			rel = strings.Replace(filepath.Dir(rel), ".ns", ".", 1)
			if len(rel) > 0 {
				rel = rel + "/"
			}
			tmp := []string{key, string(dat)}
			fmt.Print( rel )
			fmt.Println(strings.Join(tmp, ": "))
		}
	}
	// todo: json output
}

func show_help() {
	fmt.Println(` Usage: kv [OPTIONS] KEY[=VALUE]
    kv is a simple utility for getting and setting key-value pairs.
	Examples:
		$ kv foo=bar			# Set foo to bar
		$ echo spam | kv foo 	# set foo to spam
		$ kv foo				# Get value of foo
		spam

	Options:
		-h, -?, --help			Show help

		-l, --list 				List all kv pairs

		-v, --verbose			Verbose mode on

		-d, --delete KEY		Delete KEY

    `)
}

func kv_main() {

	parsed_args := parseCLI()
	// todo: add flag args

	if parsed_args.mode == ModeHelp {
		show_help()
		os.Exit(parsed_args.err_code)
	}

	lookup_key := Hashx(parsed_args.key)
	lookup_path := appdirs.UserDataDir("xac", "", "", false)
	lookup_file := lookup_path + "/" + lookup_key
	Vprint("path: ", lookup_path)
	Vprint("lfile: ", lookup_file)
	Vprint("flags: ", parsed_args.flags)
	Vprint("args: ", parsed_args.args)
	Vprint("mode: ", parsed_args.mode)

	if _, err := os.Stat(lookup_path); os.IsNotExist(err) {
		_ = os.MkdirAll(lookup_path, os.ModePerm)
	}

	switch parsed_args.mode {
	case ModeSimpleSet:
		// Store key=val in file so we can query available keys
		Store_value(lookup_file, parsed_args.key, parsed_args.val)
	case ModeSimpleGet:
		val, err := Read_value(lookup_file, parsed_args.key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
	case ModeList:
		List_all(lookup_path)
	case ModeDelete:
		err := Delete_value(lookup_file, parsed_args.key)
		if err != nil {
			log.Fatal(err)
		}
	default:
		show_help()
	}

}
