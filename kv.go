//MIT License
//
//Copyright (c) 2019 Mike McDermott
//See LICENSE file for details
//
// Simple utility for acting like a kv store for the command line

package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func hashx(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 16)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Argument needed")
		os.Exit(1)
	}

	// todo: add flag args

	key_val := strings.Split(os.Args[1], "=")

	if len(key_val) > 2 {
		log.Fatal("Too many equals")
		os.Exit(1)

	}
	key := strings.TrimSpace(key_val[0])

	usr, _ := user.Current()
	lookup_key := hashx(key)
	lookup_path := usr.HomeDir + "/.local/share/xac/"
	lookup_file := lookup_path + lookup_key

	if _, err := os.Stat(lookup_path); os.IsNotExist(err) {
		_ = os.Mkdir(lookup_path, os.ModePerm)
	}

	// Store key=val in file so we can query available keys
	if len(key_val) == 1 {
		dat, err := ioutil.ReadFile(lookup_file)
		if err != nil {
			log.Fatal("Key not found: " + key)
			os.Exit(1)
		}
		fmt.Println(strings.Split(string(dat), "=")[1])
		os.Exit(0)

	} else {
		val := strings.TrimSpace(key_val[1])
		data := []byte(key + "=" + val)
		err := ioutil.WriteFile(lookup_file, data, 0644)
		check(err)
	}

}
