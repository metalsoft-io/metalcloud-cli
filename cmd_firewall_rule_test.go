package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestFirewallRuleListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

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
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
			fw2,
			fw3,
			fw4,
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	//test json
	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	ret, err := firewallRulesListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["INDEX"].(float64))).To(Equal(0))
	Expect(r["PROTOCOL"].(string)).To(Equal(fw1.FirewallRuleProtocol))
	Expect(r["PORT"].(string)).To(Equal("22-23"))

	r = m[1].(map[string]interface{})
	Expect(int(r["INDEX"].(float64))).To(Equal(1))
	Expect(r["PORT"].(string)).To(Equal("22"))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	ret, err = firewallRulesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	ret, err = firewallRulesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 0)))
	Expect(csv[1][2]).To(Equal("22-23"))
	Expect(csv[2][1]).To(Equal(fw2.FirewallRuleProtocol))

}
