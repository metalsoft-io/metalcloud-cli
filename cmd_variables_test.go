package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestVariablesListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	variable := metalcloud.Variable{
		VariableName: "test",
	}

	list := map[string]metalcloud.Variable{
		"variable": variable,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		Variables("").
		Return(&list, nil).
		AnyTimes()

	//test json
	format := "json"
	emptyStr := ""
	cmd := Command{
		Arguments: map[string]interface{}{
			"usage":  &emptyStr,
			"format": &format,
		},
	}

	ret, err := variablesListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(0))
	Expect(r["NAME"].(string)).To(Equal(variable.VariableName))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"usage":  &emptyStr,
			"format": &format,
		},
	}

	ret, err = variablesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"usage":  &emptyStr,
			"format": &format,
		},
	}

	ret, err = variablesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 0)))
	Expect(csv[1][1]).To(Equal("test"))

}

func TestVariablesDeleteCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	variable := metalcloud.Variable{
		VariableID:   10,
		VariableName: "test",
	}

	client.EXPECT().
		VariableGet(10).
		Return(&variable, nil).
		AnyTimes()

	client.EXPECT().
		VariableDelete(10).
		Return(nil).
		AnyTimes()

	list := map[string]metalcloud.Variable{
		"variable": variable,
	}
	client.EXPECT().
		Variables("").
		Return(&list, nil).
		AnyTimes()

	//test json

	id := "test"
	bTrue := true
	cmd := Command{
		Arguments: map[string]interface{}{
			"variable_id_or_name": &id,
			"autoconfirm":         &bTrue,
		},
	}

	_, err := variableDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())

	//check with int id
	idint := 10
	cmd = Command{
		Arguments: map[string]interface{}{
			"variable_id_or_name": &idint,
			"autoconfirm":         &bTrue,
		},
	}

	_, err = variableDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())

}
