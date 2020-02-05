package main

import (
	"flag"
	"fmt"
	"strconv"

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
	Endpoint     string
}

func sameCommand(a *Command, b *Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}

const _nilDefaultStr = "__NIL__"
const _nilDefaultInt = -14234

//confirms command
func confirmCommand(c *Command, f func() string) (bool, error) {

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		return true, nil
	}

	return requestConfirmation(f())
}

//getPtrValueIfExistsOk returns a string or an int from a map of pointers if the key exists
func getPtrValueIfExistsOk(m map[string]interface{}, key string) (interface{}, bool) {

	if v := m[key]; v != nil {
		switch v.(type) {
		case *int:
			if *v.(*int) != _nilDefaultInt {
				return *v.(*int), true
			}
		case *string:
			if *v.(*string) != _nilDefaultStr {
				return *v.(*string), true
			}
		}
	}
	return nil, false
}

//getIDFromStringOk returns the id and true if valid number
func getIDFromStringOk(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	return i, err == nil
}

//verifyParam returns error if param is not present
func getParam(c *Command, label string, name string) (interface{}, error) {
	v := c.Arguments[label]
	if v == nil {
		return nil, fmt.Errorf("-%s cannot be nil", name)
	}
	switch v.(type) {
	case *int:
		if *v.(*int) <= 0 {
			return nil, fmt.Errorf("-%s cannot be <=0", name)
		}
		if *v.(*int) == _nilDefaultInt {
			return nil, fmt.Errorf("-%s is required", name)
		}
	case *string:
		if *v.(*string) == "" {
			return nil, fmt.Errorf("-%s cannot be empty", name)
		}
		if *v.(*string) == _nilDefaultStr {
			return nil, fmt.Errorf("-%s is required", name)
		}
	}
	return v, nil
}

//idOrLabel returns an int or a string contained in the interface. The last param is true if int is returned.
func idOrLabel(v interface{}) (int, string, bool) {
	switch v.(type) {
	case *int:
		return *v.(*int), "", true
	case *string:
		if i, ok := getIDFromStringOk(*v.(*string)); ok {
			return i, "", true
		}
		return 0, *v.(*string), false
	}
	return -1, "", false
}

type getIDOrDoFunc func(i string) (int, error)

func getIDOrDo(v interface{}, f getIDOrDoFunc) (int, error) {
	id, label, isID := idOrLabel(v)
	if !isID {
		return f(label)
	}
	return id, nil
}
