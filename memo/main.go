/*
Package memo is the tool for storing labels, in a shitty way, on
files in filesystems.

I am a chronir user of ~/.local/bin, and sometimes I want to remember
why the hell I've installed ~/.local/bin/satan. This is a shitty way
to just add notes in a familiar and quick way, and maintain a key
value store my computer to recheck why I originally did such a thing.

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
package memo

import (
	"bytes"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/psyomn/psy/common"
)

func memoDirPath() string      { return path.Join(common.ConfigDir(), "memo") }
func memoDataFilePath() string { return path.Join(memoDirPath(), "data.gobbin") }

type memoStore struct {
	Data map[string]string
}

func memoStoreNew() *memoStore {
	var store memoStore
	store.Data = make(map[string]string)
	return &store
}

func init() {
	if _, err := os.Stat(memoDirPath()); os.IsNotExist(err) {
		os.MkdirAll(memoDirPath(), os.ModePerm)
	}

	if _, err := os.Stat(memoDataFilePath()); os.IsNotExist(err) {
		initStore := memoStoreNew()
		store(initStore)
	}
}

func (s *memoStore) encode() (bytes.Buffer, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(*s)
	if err != nil {
		log.Println("problem encoding memo store file: ", err)
	}

	return buffer, err
}

func (s *memoStore) Add(key, value string) {
	s.Data[key] = value
}

func (s *memoStore) Get(key string) (string, bool) {
	val, ok := s.Data[key]
	return val, ok
}

func decode(cmdInFile string) *memoStore {
	var buff bytes.Buffer

	dat, err := ioutil.ReadFile(cmdInFile)
	if err != nil {
		log.Fatal("problem opening file:", cmdInFile, ":", err)
		os.Exit(1)
	}

	dec := gob.NewDecoder(&buff)
	var store memoStore
	buff.Write(dat)

	err = dec.Decode(&store)
	if err != nil {
		log.Println("problem decoding store: ", cmdInFile, ", ", err)
		os.Exit(1)
	}

	return &store
}

func mkconfig() {
	memodir := memoDirPath()
	mkdirError := os.MkdirAll(memodir, os.ModePerm)
	if mkdirError != nil {
		log.Println("problem creating memo dir directories")
	}
}

func store(memos *memoStore) {
	bytes, err := memos.encode()

	file, err := os.Create(memoDataFilePath())
	if err != nil {
		log.Println("problem opening file for storing gob: ", err)
		return
	}
	defer file.Close()

	_, err = file.Write(bytes.Bytes())
	if err != nil {
		log.Println("problem writing memo file")
		return
	}
}

func usage(fs *flag.FlagSet) error {
	fs.Usage()
	return errors.New("wrong usage")
}

// Run the memo command
// TODO: this needs some cleaning up and a better argument parsing strategy
func Run(args common.RunParams) common.RunReturn {
	type memoFlags struct {
		fileName string
		memo     string
		list     bool
	}

	sess := memoFlags{}

	memoCmd := flag.NewFlagSet("memo", flag.ExitOnError)
	memoCmd.StringVar(&sess.fileName, "file", sess.fileName, "<message> - the filename to write a memo about")
	memoCmd.BoolVar(&sess.list, "list", sess.list, "list all current memos")
	memoCmd.Parse(args)

	if sess.list {
		theStore := decode(memoDataFilePath())
		for k, v := range theStore.Data {
			fmt.Printf("%v\t%v\n", k, v)
		}
		return nil
	}

	if _, err := os.Stat(sess.fileName); os.IsNotExist(err) {
		return errors.New("fool! you can't memo what does not exist")
	}

	if sess.fileName == "" {
		return nil
	}

	absPath, err := filepath.Abs(sess.fileName)
	if err != nil {
		log.Println("problem getting abs path:", err)
		return nil
	}

	if len(memoCmd.Args()) > 0 {
		message := strings.Join(memoCmd.Args(), " ")
		theStore := decode(memoDataFilePath())
		theStore.Add(absPath, message)
		store(theStore)
		return nil
	}

	// read operations
	store := decode(memoDataFilePath())
	value, ok := store.Get(absPath)
	if !ok {
		log.Println("could not find entry for:", sess.fileName)
		return nil
	}
	fmt.Println(value)

	return nil
}
