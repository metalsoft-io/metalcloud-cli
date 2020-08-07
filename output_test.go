package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	. "github.com/onsi/gomega"
)

func TestGetTableHeader(t *testing.T) {

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
	expected := "| ID    | LABEL               | INST. |"

	actual := GetTableHeader(schema)

	if actual != expected {
		t.Errorf("Header is not correct, \nexpected:  %s\n     was: %s", expected, actual)
	}
}

func TestGetTableRow(t *testing.T) {
	RegisterTestingT(t)
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
		{
			FieldName: "INTF",
			FieldSize: 5,
			FieldType: TypeInterface,
		},
	}

	row := []interface{}{10, "test", 33.3, map[string]string{"test": "test1", "test2": "test3"}}

	actual := GetTableRow(row, schema)

	Expect(actual).To(ContainSubstring("test1"))
	Expect(actual).To(ContainSubstring("test3"))
}

func TestGetTableDelimiter(t *testing.T) {
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

	expected := "+-------+---------------------+-------+"

	actual := GetTableDelimiter(schema)

	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected: %s\n     was: %s", expected, actual)
	}
}

func TestGetTableAsString(t *testing.T) {
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

	expected :=
		`+-------+---------------------+-------+
| ID    | LABEL               | INST. |
+-------+---------------------+-------+
| 4     | str                 | 20.10 |
| 5     | st11r               | 22.10 |
| 6     | st11r444            | 2.10  |
+-------+---------------------+-------+
`
	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	actual := GetTableAsString(data, schema)

	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected:\n%s\nwas:\n%s", expected, actual)
	}
}

func TestGetTableAsJSONString(t *testing.T) {
	RegisterTestingT(t)
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

	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	ret, err := GetTableAsJSONString(data, schema)
	if err != nil {
		t.Errorf("%s", err)
	}

	var m []interface{}

	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

	Expect(int(m[0].(map[string]interface{})["ID"].(float64))).To(Equal(data[0][0]))
	Expect(m[0].(map[string]interface{})["LABEL"]).To(Equal(data[0][1]))
	Expect(m[2].(map[string]interface{})["LABEL"]).To(Equal(data[2][1]))
	Expect(float64(m[2].(map[string]interface{})["INST."].(float64))).To(Equal(data[2][2]))

}

func TestGetTableAsCSVString(t *testing.T) {

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

	expected :=
		`ID,LABEL,INST.
4,str,20.100000
5,st11r,22.100000
6,st11r444,2.100000
`

	data := [][]interface{}{
		{4, "str", 20.1},
		{5, "st11r", 22.1},
		{6, "st11r444", 2.1},
	}

	actual, err := GetTableAsCSVString(data, schema)
	if err != nil {
		t.Errorf("%s", err)
	}
	if actual != expected {
		t.Errorf("Delimiter is not correct, \nexpected:\n%s\nwas:\n%s", expected, actual)
	}
}

func TestAdjustFieldSizes(t *testing.T) {
	RegisterTestingT(t)
	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 3,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5, //this is smaller than the largest
		},
		{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      0,
			FieldPrecision: 4,
		},
		{
			FieldName:      "VERY LONG FIELD NAME",
			FieldType:      TypeString,
			FieldSize:      4,
			FieldPrecision: 4,
		},
	}

	data := [][]interface{}{
		{4, "12345", 20.1, "tes"},
		{5, "12", 22.1, "te"},
		{6, "123456789", 1.2345, "t"},
	}

	AdjustFieldSizes(data, &schema)

	Expect(schema[0].FieldSize).To(Equal(3))
	Expect(schema[1].FieldSize).To(Equal(10))
	Expect(schema[2].FieldSize).To(Equal(8))
	//test if expands with LABEl
	Expect(schema[3].FieldSize).To(Equal(21))

}

func TestRenderTable(t *testing.T) {
	RegisterTestingT(t)
	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 3,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      0,
			FieldPrecision: 4,
		},
		{
			FieldName:      "VERY LONG FIELD NAME",
			FieldType:      TypeString,
			FieldSize:      4,
			FieldPrecision: 4,
		},
	}

	data := [][]interface{}{
		{4, "12345", 20.1, "tes"},
		{5, "12", 22.1, "te"},
		{6, "123456789", 1.2345, "t"},
	}

	s, err := renderTable("test", "", "", data, schema)

	Expect(err).To(BeNil())
	Expect(s).To(ContainSubstring("test"))
	Expect(s).To(ContainSubstring("VERY LONG"))

	s, err = renderTable("test", "", "json", data, schema)
	Expect(err).To(BeNil())
	var m []interface{}
	err = json.Unmarshal([]byte(s), &m)
	Expect(err).To(BeNil())

	s, err = renderTable("test", "", "csv", data, schema)
	Expect(err).To(BeNil())

	s, err = renderTable("test", "", "yaml", data, schema)
	err = yaml.Unmarshal([]byte(s), &m)
	Expect(err).To(BeNil())
}

func TestYAMLMArshalOfMetalcloudObjects(t *testing.T) {
	RegisterTestingT(t)

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	Expect(err).To(BeNil())

	b, err := yaml.Marshal(sw)
	Expect(err).To(BeNil())

	t.Log(string(b))

	var sw2 metalcloud.SwitchDevice

	err = yaml.Unmarshal(b, &sw2)
	Expect(err).To(BeNil())
	Expect(sw2.NetworkEquipmentPrimaryWANIPv4SubnetPool).To(Equal(sw.NetworkEquipmentPrimaryWANIPv4SubnetPool))
	//for some reason this doesn't work. don't know why yet
	//t.Logf("sw1=%+v", sw)
	//t.Logf("sw2=%+v", sw2)
	//Expect(reflect.DeepEqual(sw, sw2)).To(BeTrue())
}

func TestYAMLMArshalCaseSensitivity(t *testing.T) {
	RegisterTestingT(t)

	type dummy struct {
		WithCamelCase1 string
		WithCamelCase2 int `yaml:"withCamelCase2"`
		A              int
	}

	var d dummy

	s := `
withcamelcase1: test
withCamelCase2: 10
a: 12
`

	err := yaml.Unmarshal([]byte(s), &d)
	Expect(err).To(BeNil())
	Expect(d.A).To(Equal(12))
	Expect(d.WithCamelCase2).To(Equal(10))
	Expect(d.WithCamelCase1).To(Equal("test"))

}

func JSONUnmarshal(jsonString string) ([]interface{}, error) {
	var m []interface{}
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

//JSONFirstRowEquals checks if values of the table returned in the json match the values provided. Type is not checked (we check string equality)
func JSONFirstRowEquals(jsonString string, testVals map[string]interface{}) error {
	m, err := JSONUnmarshal(jsonString)
	if err != nil {
		return err
	}

	firstRow := m[0].(map[string]interface{})

	for k, v := range testVals {
		if fmt.Sprintf("%+v", firstRow[k]) != fmt.Sprintf("%+v", v) {
			return fmt.Errorf("values for key %s do not match:  expected '%+v' provided '%+v'", k, v, firstRow[k])
		}
	}

	return nil
}

func CSVUnmarshal(csvString string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(csvString))

	return reader.ReadAll()
}

//CSVFirstRowEquals checks if values of the table returned in the json match the values provided. Type is not checked (we check string equality)
func CSVFirstRowEquals(csvString string, testVals map[string]interface{}) error {
	m, err := CSVUnmarshal(csvString)
	if err != nil {
		return err
	}

	header := m[0]
	firstRow := map[string]string{}
	//turn first row into a map
	for k, v := range m[1] {
		firstRow[header[k]] = v
	}

	for k, v := range testVals {
		if fmt.Sprintf("%+v", firstRow[k]) != fmt.Sprintf("%+v", v) {
			return fmt.Errorf("values for key %s do not match:  expected '%+v' provided '%+v'", k, v, firstRow[k])
		}
	}

	return nil
}

func TestTransposeTable(t *testing.T) {
	RegisterTestingT(t)
	data := [][]interface{}{
		{11, 12, 13},
		{21, 22, 23},
		{31, 32, 33},
	}

	dataT := transposeTable(data)

	expectedDataT := [][]interface{}{
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
	}

	Expect(dataT).Should(Equal(expectedDataT))
}

func TestConvertToStringTable(t *testing.T) {
	RegisterTestingT(t)
	data := [][]interface{}{
		{11, "12", 13.4},
		{21, "22", 23.3},
		{31, "32", 33.4},
	}

	dataT := convertToStringTable(data)

	expectedDataT := [][]interface{}{
		{"11", "12", "13.4"},
		{"21", "22", "23.3"},
		{"31", "32", "33.4"},
	}

	Expect(dataT).Should(Equal(expectedDataT))
}

func TestRenderTransposedTable(t *testing.T) {
	RegisterTestingT(t)
	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 3,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName:      "INST.",
			FieldType:      TypeFloat,
			FieldSize:      0,
			FieldPrecision: 4,
		},
		{
			FieldName:      "VERY LONG FIELD NAME",
			FieldType:      TypeString,
			FieldSize:      4,
			FieldPrecision: 4,
		},
	}

	data := [][]interface{}{
		{4, "12345", 20.1, "tes"},
	}

	s, err := renderTransposedTable("test", "", "", data, schema)

	Expect(err).To(BeNil())
	Expect(s).To(ContainSubstring("KEY"))
	Expect(s).To(ContainSubstring("VALUE"))
	Expect(s).To(ContainSubstring("12345"))
	Expect(s).To(ContainSubstring("20.1"))

	s, err = renderTransposedTable("test", "", "json", data, schema)
	Expect(err).To(BeNil())
	var m []interface{}
	err = json.Unmarshal([]byte(s), &m)
	Expect(err).To(BeNil())

	s, err = renderTransposedTable("test", "", "csv", data, schema)
	Expect(err).To(BeNil())
}

func TestObjectToTable(t *testing.T) {

	RegisterTestingT(t)

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	Expect(err).To(BeNil())

	d, s, err := objectToTable(sw)
	Expect(err).To(BeNil())
	Expect(len(d)).To(Equal(40))
	Expect(d[1]).To(Equal("UK_RDG_EVR01_00_0001_00A9_01"))
	Expect(s[1].FieldName).To(Equal("network equipment identifier string"))
	Expect(s[39].FieldName).To(Equal("volume template id"))
	Expect(d[39]).To(Equal(0))
}

func TestObjToTableWithFormatter(t *testing.T) {
	RegisterTestingT(t)

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	Expect(err).To(BeNil())

	d, s, err := objectToTableWithFormatter(sw, NewStripPrefixFormatter("NetworkEquipment"))
	Expect(err).To(BeNil())
	Expect(len(d)).To(Equal(40))
	Expect(d[1]).To(Equal("UK_RDG_EVR01_00_0001_00A9_01"))
	Expect(s[1].FieldName).To(Equal("Identifier String"))
	Expect(s[39].FieldName).To(Equal("Volume Template Id"))
	Expect(d[39]).To(Equal(0))

}

func TestRenderTransposedTableHumanReadable(t *testing.T) {
	RegisterTestingT(t)

	schema := []SchemaField{
		{
			FieldName: "Field1",
			FieldType: TypeInt,
		},
		{
			FieldName: "Field2",
			FieldType: TypeString,
		},
	}

	data := [][]interface{}{
		{
			10,
			"test",
		},
	}

	s, err := renderTransposedTableHumanReadable("test", "test", data, schema)

	Expect(err).To(BeNil())
	Expect(s).To(Equal(`Field1: 10
Field2: test
`))

}

func TestRenderRawObject(t *testing.T) {
	RegisterTestingT(t)

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	Expect(err).To(BeNil())

	ret, err := renderRawObject(sw, "json", "")

	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeNil())
	Expect(ret).To(ContainSubstring("2A02:0CB8:0000:0000:0000:0000:0000:0000/53"))

	var sw2 metalcloud.SwitchDevice
	err = json.Unmarshal([]byte(ret), &sw2)
	Expect(err).To(BeNil())

	ret, err = renderRawObject(sw, "yaml", "")
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeNil())
	Expect(ret).To(ContainSubstring("2A02:0CB8:0000:0000:0000:0000:0000:0000/53"))

	ret, err = renderRawObject(sw, "csv", "")
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeNil())
	Expect(ret).To(ContainSubstring("2A02:0CB8:0000:0000:0000:0000:0000:0000/53"))

	ret, err = renderRawObject(sw, "", "NetworkEquipment")
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeNil())
	Expect(ret).To(ContainSubstring("2A02:0CB8:0000:0000:0000:0000:0000:0000/53"))

	t.Log(ret)

}
