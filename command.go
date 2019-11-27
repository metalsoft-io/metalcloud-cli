package main

import (
	"flag"
	"fmt"
	"strconv"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
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
	ExecuteFunc  func(c *Command, client interfaces.MetalCloudClient) (string, error)
}

func sameCommand(a *Command, b *Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}

const _nilDefaultStr = "__NIL__"
const _nilDefaultInt = -14234

func getIDFromCommand(c *Command, label string) (metalcloud.ID, error) {
	v := c.Arguments[label]

	if v == nil {
		return nil, fmt.Errorf("id is required")
	}

	switch v.(type) {
	case *int:
		return *c.Arguments[label].(*int), nil
	case *string:
		if id, err := strconv.Atoi(*v.(*string)); err == nil {
			return id, nil
		}
		return *c.Arguments[label].(*string), nil
	}

	return nil, fmt.Errorf("could not determinte the type of the passed ID")
}
