package firewall

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
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
	/*
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
	*/

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:              11,
		InstanceArrayLabel:           "testia-edited",
		InstanceArrayDeployType:      "edit",
		InstanceArrayDeployStatus:    "not_started",
		InstanceArrayFirewallManaged: true,
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
			//	fw2,
			//	fw3,
			//	fw4,
		},
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:              11,
		InstanceArrayLabel:           "testia",
		InfrastructureID:             infra.InfrastructureID,
		InstanceArrayOperation:       &iao,
		InstanceArrayServiceStatus:   "active",
		InstanceArrayFirewallManaged: true,
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
			//	fw2,
			//	fw3,
			//	fw4,
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	//test json
	format := "json"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	ret, err := firewallRuleListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["INDEX"].(float64))).To(Equal(0))
	Expect(r["PROTOCOL"].(string)).To(Equal(fw1.FirewallRuleProtocol))
	Expect(r["PORT"].(string)).To(Equal("22-23"))

	//test plaintext
	format = ""
	cmd = command.Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	ret, err = firewallRuleListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = command.Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	ret, err = firewallRuleListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 0)))
	Expect(csv[1][2]).To(Equal("22-23"))

}

func TestFirewallRuleListWithFWDisabledCmd(t *testing.T) {
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

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:              11,
		InstanceArrayLabel:           "testia-edited",
		InstanceArrayDeployType:      "edit",
		InstanceArrayDeployStatus:    "not_started",
		InstanceArrayFirewallManaged: false,
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
		},
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:              11,
		InstanceArrayLabel:           "testia",
		InfrastructureID:             infra.InfrastructureID,
		InstanceArrayOperation:       &iao,
		InstanceArrayServiceStatus:   "active",
		InstanceArrayFirewallManaged: false,
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{
			fw1,
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	//test json
	format := "json"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"format":            &format,
			"instance_array_id": &ia.InstanceArrayID,
		},
	}

	_, err := firewallRuleListCmd(&cmd, client)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(SatisfyAll(ContainSubstring("firewall"), ContainSubstring("disabled")))

}

func TestFirewallRuleAddCmd(t *testing.T) {
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
	protocol := "tcp"
	port := "22-34"
	source := "192.172.10.10-192.172.10.15"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"firewall_rule_protocol":          &protocol,
			"firewall_rule_port":              &port,
			"firewall_rule_source_ip_address": &source,
			"instance_array_id":               &ia.InstanceArrayID,
		},
	}

	expectedIAO := iao

	fw := metalcloud.FirewallRule{
		FirewallRuleProtocol:                  "tcp",
		FirewallRulePortRangeStart:            22,
		FirewallRulePortRangeEnd:              34,
		FirewallRuleSourceIPAddressRangeStart: "192.172.10.10",
		FirewallRuleSourceIPAddressRangeEnd:   "192.172.10.15",
	}

	expectedIAO.InstanceArrayFirewallRules = append(expectedIAO.InstanceArrayFirewallRules, fw)

	client.EXPECT().
		InstanceArrayEdit(ia.InstanceArrayID, expectedIAO, gomock.Any(), nil, nil, nil).
		Return(&ia, nil).
		AnyTimes()

	_, err := firewallRuleAddCmd(&cmd, client)
	Expect(err).To(BeNil())

}

func TestFirewallRuleRemoveCmd(t *testing.T) {
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
	protocol := "tcp"
	port := "22"
	source := "192.168.0.1"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"firewall_rule_protocol":          &protocol,
			"firewall_rule_port":              &port,
			"firewall_rule_source_ip_address": &source,
			"instance_array_id":               &ia.InstanceArrayID,
		},
	}

	expectedIAO := iao

	expectedIAO.InstanceArrayFirewallRules = []metalcloud.FirewallRule{
		fw1,
		fw2, //we skipped fw3 which should be deleted
		fw4,
	}

	client.EXPECT().
		InstanceArrayEdit(ia.InstanceArrayID, expectedIAO, gomock.Any(), nil, nil, nil).
		Return(&ia, nil).
		AnyTimes()

	_, err := firewallRuleDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())

}

func TestPortStringToRange(t *testing.T) {
	RegisterTestingT(t)

	s, e, err := portStringToRange("12")
	Expect(err).To(BeNil())
	Expect(s).To(Equal(12))
	Expect(e).To(Equal(12))

	s, e, err = portStringToRange("12-35")
	Expect(err).To(BeNil())
	Expect(s).To(Equal(12))
	Expect(e).To(Equal(35))

	s, e, err = portStringToRange("-35")
	Expect(err).NotTo(BeNil())

	s, e, err = portStringToRange("44-")
	Expect(err).NotTo(BeNil())

	s, e, err = portStringToRange("44&33")
	Expect(err).NotTo(BeNil())
}

func TestAddressStringToRange(t *testing.T) {
	RegisterTestingT(t)

	s, e, err := addressStringToRange("192.168.0.1")
	Expect(err).To(BeNil())
	Expect(s).To(Equal("192.168.0.1"))
	Expect(e).To(Equal("192.168.0.1"))

	s, e, err = addressStringToRange("192.168.0.1-192.168.0.100")
	Expect(err).To(BeNil())
	Expect(s).To(Equal("192.168.0.1"))
	Expect(e).To(Equal("192.168.0.100"))

	s, e, err = addressStringToRange("192.168.0.1-")
	Expect(err).NotTo(BeNil())

	s, e, err = addressStringToRange("-192.168.0.1")
	Expect(err).NotTo(BeNil())

	s, e, err = addressStringToRange("-192.168.0.1--192.168.0.1")
	Expect(err).NotTo(BeNil())

}
