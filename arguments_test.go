package main

import (
	"flag"
	"testing"
)

func TestSimpleArgument(t *testing.T) {

	var executed = false

	cmd := Command{
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "c",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_instance_count": c.FlagSet.Int("instance_count", 0, "Instance count of this instance array"),
				"instance_array_instance_label": c.FlagSet.String("label", "", "Instance array's label"),
			}
		},
		ExecuteFunc: func(c *Command, client MetalCloudClient) (string, error) {
			executed = true
			return "retstr", nil
		},
	}

	cmd.InitFunc(&cmd)

	argv := []string{
		"-instance_count=3",
		"-label=test",
	}

	err := cmd.FlagSet.Parse(argv)
	if err != nil {
		t.Errorf("%s", err)
	}

	iaCount := cmd.Arguments["instance_array_instance_count"].(*int)
	if iaCount == nil || *iaCount != 3 {
		t.Errorf("instance_array_instance_count expected to be %d\n\twas %d", 3, *iaCount)
	}

	iaLabel := cmd.Arguments["instance_array_instance_label"].(*string)

	if iaLabel == nil || *iaLabel != "test" {
		t.Errorf("instance_array_label expected to be %s\n\twas %s", "test", *iaLabel)
	}

	argv = []string{
		"instance_countasdad=3",
		"la33bel=\"test\"",
	}

	err = cmd.FlagSet.Parse(argv)
	if err != nil {
		t.Errorf("%s", err)
	}

	ret, err := cmd.ExecuteFunc(&cmd, nil)
	if err != nil {
		t.Errorf("%s", err)
	}

	if !executed || ret != "retstr" {
		t.Errorf("ExecuteFunction not called properly")
	}
}
