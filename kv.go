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

	key_val := strings.Split(os.Args[1], "=")

	if len(key_val) > 2 {
		log.Fatal("Too many equals")
		os.Exit(1)

	}
	key := strings.TrimSpace(key_val[0])

	usr, _ := user.Current()
	lookup_key := hashx(key)
	lookup_file := usr.HomeDir + "/.local/share/xac/" + lookup_key


	if len(key_val) == 1 {
		dat, err := ioutil.ReadFile(lookup_file)
		if err != nil {
			log.Fatal("Key not found: " + key)
			os.Exit(1)
		}
		fmt.Println(string(dat))
		os.Exit(0)

	} else {
		data := []byte(strings.TrimSpace(key_val[1]))
		err := ioutil.WriteFile(lookup_file, data, 0644)
		check(err)
	}

}
