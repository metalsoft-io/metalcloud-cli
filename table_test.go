package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"

	. "github.com/onsi/gomega"
)

func TestTableSortWithSchema(t *testing.T) {

	data := [][]interface{}{
		{4, "str", 20.1},
		{6, "st11r", 22.1},
		{5, "wt11r444", 2.3},
		{5, "wt11r444", 2.1},
		{5, "at11r43", 2.2},
		{4, "xxxx", 2.2},
	}

	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      6,
			FieldPrecision: 2,
		},
	}

	TableSorter(schema).OrderBy("LABEL", "ID", "INST.").Sort(data)

	expected := [][]interface{}{
		{5, "at11r43", 2.2},
		{6, "st11r", 22.1},
		{4, "str", 20.1},
		{5, "wt11r444", 2.1},
		{5, "wt11r444", 2.3},
		{4, "xxxx", 2.2},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("the sorted array was not correct \nwas:\n%+v\n expected\n %+v\n", data, expected)
	}
}

func TestTableSortWithSchemaWithDateTime(t *testing.T) {

	data := [][]interface{}{
		{4, "str", "2013-11-29T13:00:01Z"},
		{6, "st11r", "2013-11-29T13:00:02Z"},
		{6, "st11r", "2014-11-29T13:00:03Z"},
		{6, "st11r", "2012-11-29T13:00:03Z"},
	}

	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName:   "DATE",
			FieldType:   TypeDateTime,
			FieldSize:   6,
			FieldFormat: defaultTimeFormat,
		},
	}

	TableSorter(schema).OrderBy("DATE").Sort(data)

	expected := [][]interface{}{
		{6, "st11r", "2012-11-29T13:00:03Z"},
		{4, "str", "2013-11-29T13:00:01Z"},
		{6, "st11r", "2013-11-29T13:00:02Z"},
		{6, "st11r", "2014-11-29T13:00:03Z"},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("the sorted array was not correct \nwas:\n%+v\n expected\n %+v\n", data, expected)
	}
}

func TestDefaultTimeFormat(t *testing.T) {

	layout := defaultTimeFormat

	s := "2012-11-29T13:00:03Z"

	tm, err := time.Parse(layout, s)

	if err != nil {
		t.Errorf("error converting time string %s", err)
	}

	if tm.Year() != 2012 || tm.Second() != 3 {
		t.Error("Date was not parsed correctly")
	}

}

func TestGetTableAsJSONRegressionTest1(t *testing.T) {
	RegisterTestingT(t)
	fw1 := metalcloud.FirewallRule{
		FirewallRuleDescription:    "test desc",
		FirewallRuleProtocol:       "tcp",
		FirewallRulePortRangeStart: 22,
		FirewallRulePortRangeEnd:   23,
	}

	fw2 := metalcloud.FirewallRule{
		FirewallRuleProtocol:       "udp",
		FirewallRulePortRangeStart: 22,
		FirewallRulePortRangeEnd:   22,
	}

	fw3 := metalcloud.FirewallRule{
		FirewallRuleProtocol:                  "tcp",
		FirewallRulePortRangeStart:            22,
		FirewallRulePortRangeEnd:              22,
		FirewallRuleSourceIPAddressRangeStart: "192.168.0.1",
		FirewallRuleSourceIPAddressRangeEnd:   "192.168.0.1",
	}

	fw4 := metalcloud.FirewallRule{
		FirewallRuleProtocol:                  "tcp",
		FirewallRulePortRangeStart:            22,
		FirewallRulePortRangeEnd:              22,
		FirewallRuleSourceIPAddressRangeStart: "192.168.0.1",
		FirewallRuleSourceIPAddressRangeEnd:   "192.168.0.100",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia-edited",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
			fw2,
			fw3,
			fw4,
		},
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           100,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
			fw2,
			fw3,
			fw4,
		},
	}

	list := ia.InstanceArrayOperation.InstanceArrayFirewallRules
	data := [][]interface{}{}
	idx := 0

	for _, fw := range list {

		portRange := "any"

		if fw.FirewallRulePortRangeStart != 0 {
			portRange = fmt.Sprintf("%d", fw.FirewallRulePortRangeStart)
		}

		if fw.FirewallRulePortRangeStart != fw.FirewallRulePortRangeEnd {
			portRange += fmt.Sprintf("-%d", fw.FirewallRulePortRangeEnd)
		}

		sourceIPRange := "any"

		if fw.FirewallRuleSourceIPAddressRangeStart != "" {
			sourceIPRange = fw.FirewallRuleSourceIPAddressRangeStart
		}

		if fw.FirewallRuleSourceIPAddressRangeStart != fw.FirewallRuleSourceIPAddressRangeEnd {
			sourceIPRange += fmt.Sprintf("-%s", fw.FirewallRuleSourceIPAddressRangeEnd)
		}

		data = append(data, []interface{}{
			idx,
			fw.FirewallRuleProtocol,
			portRange,
			sourceIPRange,
			fw.FirewallRuleEnabled,
			fw.FirewallRuleDescription,
		})

		idx++

	}

	schema := []SchemaField{
		{
			FieldName: "INDEX",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "PROTOCOL",
			FieldType: TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "PORT",
			FieldType: TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "SOURCE",
			FieldType: TypeString,
			FieldSize: 20,
		},

		{
			FieldName: "ENABLED",
			FieldType: TypeBool,
			FieldSize: 10,
		},
		{
			FieldName: "DESC.",
			FieldType: TypeString,
			FieldSize: 50,
		},
	}

	Expect(data[0][0]).NotTo(Equal(data[0][1]))
	Expect(data[0][1]).NotTo(Equal(data[0][2]))
	Expect(data[0][1]).NotTo(Equal(data[0][2]))

	ret, err := GetTableAsJSONString(data, schema)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

	Expect(m[0].(map[string]interface{})["INDEX"]).ToNot(Equal(m[1].(map[string]interface{})["INDEX"]))
	Expect(m[0].(map[string]interface{})["INDEX"]).ToNot(Equal(m[2].(map[string]interface{})["INDEX"]))
	Expect(m[1].(map[string]interface{})["INDEX"]).ToNot(Equal(m[2].(map[string]interface{})["INDEX"]))
}
