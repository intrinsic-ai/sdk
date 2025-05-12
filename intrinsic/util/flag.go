// Copyright 2023 Intrinsic Innovation LLC

// Package flag enhances core flag package with additional functionality.
package flag

import (
	"flag"
	"fmt"
	"strings"
)

type stringList []string

func (sl *stringList) String() string {
	return strings.Join(*sl, ",")
}

func (sl *stringList) Get() any {
	return []string(*sl)
}

func (sl *stringList) Set(s string) error {
	*sl = stringList(strings.Split(s, ","))
	return nil
}

// StringList registers a flag of type []string which splits values on commas.
func StringList(name string, value []string, usage string) *[]string {
	var v []string = value
	flag.Var((*stringList)(&v), name, usage)
	return &v
}

type multiString []string

func (ms *multiString) String() string {
	if len(*ms) == 0 {
		return ""
	}
	return fmt.Sprint(*ms)
}

func (ms *multiString) Get() any {
	return []string(*ms)
}

func (ms *multiString) Set(value string) error {
	*ms = append(*ms, value)
	return nil
}

// MultiString registers a flag of type []string with appending semantics.
func MultiString(name string, value []string, usage string) *[]string {
	var v []string = value
	flag.Var((*multiString)(&v), name, usage)
	return &v
}
