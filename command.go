package main

import (
	"flag"
)

//Command defines a command, arguments, description etc
type Command struct {
	Description  string
	Subject      string
	AltSubject   string
	Predicate    string
	AltPredicate string
	FlagSet      *flag.FlagSet
	Arguments    map[string]interface{}
	InitFunc     func(c *Command)
	ExecuteFunc  func(c *Command, client MetalCloudClient) (string, error)
}
