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

func TestServersListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	server := metalcloud.ServerSearchResult{
		ServerID:          100,
		ServerProductName: "test",
	}

	list := []metalcloud.ServerSearchResult{
		server,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		ServersSearch("").
		Return(&list, nil).
		AnyTimes()

	//test json
	format := "json"
	emptyStr := ""
	cmd := Command{
		Arguments: map[string]interface{}{
			"filter": &emptyStr,
			"format": &format,
		},
	}

	ret, err := serversListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(100))
	Expect(r["PRODUCT_NAME"].(string)).To(Equal(server.ServerProductName))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"filter": &emptyStr,
			"format": &format,
		},
	}

	ret, err = serversListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"filter": &emptyStr,
			"format": &format,
		},
	}

	ret, err = serversListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 100)))
	Expect(csv[1][5]).To(Equal("test"))

}

func TestServerGetCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	serverType := metalcloud.ServerType{
		ServerTypeID:   100,
		ServerTypeName: "testtype",
	}

	server := metalcloud.Server{
		ServerID:          10,
		ServerProductName: "test",
		ServerTypeID:      100,
	}

	client.EXPECT().
		ServerGet(10, false).
		Return(&server, nil).
		AnyTimes()

	client.EXPECT().
		ServerTypeGet(100).
		Return(&serverType, nil).
		AnyTimes()

	//test json
	id := 10
	format := "json"

	cmd := Command{
		Arguments: map[string]interface{}{
			"id":     &id,
			"format": &format,
		},
	}

	ret, err := serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(10))
	Expect(r["PRODUCT_NAME"].(string)).To(Equal(server.ServerProductName))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"id":     &id,
			"format": &format,
		},
	}

	ret, err = serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"id":     &id,
			"format": &format,
		},
	}

	ret, err = serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 10)))
	Expect(csv[1][5]).To(Equal("test"))

}
