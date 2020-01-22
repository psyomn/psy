/*
more experimental than anything.

Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/psyomn/psy/barf"
	"github.com/psyomn/psy/common"
	"github.com/psyomn/psy/gh"
	"github.com/psyomn/psy/memo"
	"github.com/psyomn/psy/mock"
	"github.com/psyomn/psy/uploader"
)

type command struct {
	name string
	fn   func(common.RunParams) common.RunReturn
	desc string
}

var commands []command

func init() {
	commands = []command{
		{"barf", barf.Run, "run code barfer"},
		{"memo", memo.Run, "description on files in the system"},
		{"upld", uploader.Run, "run the uploader tool"},
		{"mock", mock.Run, "run tcp/udp mocker"},
		{"gh", gh.Run, "personal github utils"},
		{"help", help, "print help"},
	}
}

func help(_ common.RunParams) common.RunReturn {
	fmt.Println("usage:")
	for _, c := range commands {
		fmt.Println("\t", c.name, "\t", c.desc)
	}
	return nil
}

func handleRet(ret common.RunReturn) {
	if ret != nil {
		log.Println("error: ", ret)
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		help(nil)
		os.Exit(1)
	}

	cmd := args[1]
	rest := args[2:]
	var callfn func(common.RunParams) common.RunReturn

	{
		found := false
		for _, c := range commands {
			if cmd == c.name {
				callfn = c.fn
				found = true
				break
			}
		}

		if !found {
			log.Println("no such command: ", cmd)
			os.Exit(1)
		}
	}

	err := callfn(rest)
	if err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}
