package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var firewallRuleCmds = []Command{

	Command{
		Description:  "Lists instance array firewall rules",
		Subject:      "firewall_rule",
		AltSubject:   "fw",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list firewall rules", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id": c.FlagSet.Int("ia", _nilDefaultInt, "(Required) The instance array id"),
				"format":            c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: firewallRuleListCmd,
	},
	Command{
		Description:  "Add instance array firewall rule",
		Subject:      "firewall_rule",
		AltSubject:   "fw",
		Predicate:    "add",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("add firewall rules", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id":                   c.FlagSet.Int("ia", _nilDefaultInt, "(Required) The instance array id"),
				"firewall_rule_protocol":              c.FlagSet.String("protocol", _nilDefaultStr, "The protocol of the firewall rule. Possible values: all, icmp, tcp, udp."),
				"firewall_rule_ip_address_type":       c.FlagSet.String("ip_address_type", "ipv4", "The IP address type of the firewall rule. Possible values: ipv4, ipv6."),
				"firewall_rule_port":                  c.FlagSet.String("port", _nilDefaultStr, "The port to filter on. It can also be a range with the start and end values separated by a dash."),
				"firewall_rule_source_ip_address":     c.FlagSet.String("source", _nilDefaultStr, "The source address to filter on. It can also be a range with the start and end values separated by a dash."),
				"firewall_rule_desination_ip_address": c.FlagSet.String("destination", _nilDefaultStr, "The destination address to filter on. It can also be a range with the start and end values separated by a dash."),
				"firewall_rule_description":           c.FlagSet.String("description", _nilDefaultStr, "The firewall rule's description."),
			}
		},
		ExecuteFunc: firewallRuleAddCmd,
	},
	Command{
		Description:  "Remove instance array firewall rule",
		Subject:      "firewall_rule",
		AltSubject:   "fw",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete firewall rules", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id":                   c.FlagSet.Int("ia", _nilDefaultInt, "(Required) The instance array id"),
				"firewall_rule_ip_address_type":       c.FlagSet.String("ip_address_type", "ipv4", "The IP address type of the firewall rule. Possible values: ipv4, ipv6."),
				"firewall_rule_protocol":              c.FlagSet.String("protocol", _nilDefaultStr, "The protocol of the firewall rule. Possible values: all, icmp, tcp, udp."),
				"firewall_rule_port":                  c.FlagSet.String("port", _nilDefaultStr, "The port to filter on. It can also be a range with the start and end values separated by a dash."),
				"firewall_rule_source_ip_address":     c.FlagSet.String("source", _nilDefaultStr, "The source address to filter on. It can also be a range with the start and end values separated by a dash."),
				"firewall_rule_desination_ip_address": c.FlagSet.String("destination", _nilDefaultStr, "The destination address to filter on. It can also be a range with the start and end values separated by a dash."),
			}
		},
		ExecuteFunc: firewallRuleDeleteCmd,
	},
}

func firewallRuleListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	instanceArrayID := c.Arguments["instance_array_id"]

	if instanceArrayID == nil || *instanceArrayID.(*int) == 0 {
		return "", fmt.Errorf("-ia <instance_array_id> is required")
	}

	retIA, err := client.InstanceArrayGet(*instanceArrayID.(*int))
	if err != nil {
		return "", err
	}

	if !retIA.InstanceArrayOperation.InstanceArrayFirewallManaged {
		return "", fmt.Errorf("the instance array %s [#%d] has firewall management disabled", retIA.InstanceArrayLabel, retIA.InstanceArrayID)
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "INDEX",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "PROTOCOL",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "PORT",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "SOURCE",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "DEST",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "TYPE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "ENABLED",
			FieldType: TypeBool,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "DESC.",
			FieldType: TypeString,
			FieldSize: 50,
		},
	}

	status := retIA.InstanceArrayServiceStatus
	if retIA.InstanceArrayServiceStatus != "ordered" && retIA.InstanceArrayOperation.InstanceArrayDeployType == "edit" && retIA.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
		status = "edited"
	}

	list := retIA.InstanceArrayOperation.InstanceArrayFirewallRules
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

		destinationIPRange := "any"

		if fw.FirewallRuleDestinationIPAddressRangeStart != "" {
			sourceIPRange = fw.FirewallRuleSourceIPAddressRangeStart
		}

		if fw.FirewallRuleDestinationIPAddressRangeStart != fw.FirewallRuleDestinationIPAddressRangeEnd {
			sourceIPRange += fmt.Sprintf("-%s", fw.FirewallRuleDestinationIPAddressRangeEnd)
		}

		data = append(data, []interface{}{
			idx,
			fw.FirewallRuleProtocol,
			portRange,
			sourceIPRange,
			destinationIPRange,
			fw.FirewallRuleIPAddressType,
			fw.FirewallRuleEnabled,
			fw.FirewallRuleDescription,
		})

		idx++

	}

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		sb.WriteString(fmt.Sprintf("Instance Array %s (%d) [%s] has the following firewall rules:\n", retIA.InstanceArrayLabel, retIA.InstanceArrayID, status))

		AdjustFieldSizes(data, &schema)
		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d firewall rules\n\n", len(list)))
	}

	return sb.String(), nil
}

func firewallRuleAddCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	instanceArrayID := c.Arguments["instance_array_id"]

	if instanceArrayID == nil || *instanceArrayID.(*int) == 0 {
		return "", fmt.Errorf("-ia <instance_array_id> is required")
	}

	retIA, err := client.InstanceArrayGet(*instanceArrayID.(*int))
	if err != nil {
		return "", err
	}

	fw := metalcloud.FirewallRule{}

	if v := c.Arguments["firewall_rule_protocol"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleProtocol = *v.(*string)
	}

	if v := c.Arguments["firewall_rule_ip_address_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleIPAddressType = *v.(*string)
	}

	if v := c.Arguments["firewall_rule_port"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRulePortRangeStart, fw.FirewallRulePortRangeEnd, err = portStringToRange(*v.(*string))
		if err != nil {
			return "", err
		}
	}

	if v := c.Arguments["firewall_rule_source_ip_address"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleSourceIPAddressRangeStart, fw.FirewallRuleSourceIPAddressRangeEnd, err = addressStringToRange(*v.(*string))
		if err != nil {
			return "", err
		}
	}

	if v := c.Arguments["firewall_rule_desination_ip_address"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleDestinationIPAddressRangeStart, fw.FirewallRuleDestinationIPAddressRangeEnd, err = addressStringToRange(*v.(*string))
		if err != nil {
			return "", err
		}
	}

	if v := c.Arguments["firewall_rule_description"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleDescription = *v.(*string)
	}

	retIA.InstanceArrayOperation.InstanceArrayFirewallRules = append(
		retIA.InstanceArrayOperation.InstanceArrayFirewallRules,
		fw)

	bFalse := false
	_, err = client.InstanceArrayEdit(retIA.InstanceArrayID, *retIA.InstanceArrayOperation, &bFalse, nil, nil, nil)

	return "", err
}

func firewallRuleDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	instanceArrayID := c.Arguments["instance_array_id"]

	if instanceArrayID == nil || *instanceArrayID.(*int) == 0 {
		return "", fmt.Errorf("-ia <instance_array_id> is required")
	}

	retIA, err := client.InstanceArrayGet(*instanceArrayID.(*int))
	if err != nil {
		return "", err
	}

	fw := metalcloud.FirewallRule{}

	if v := c.Arguments["firewall_rule_protocol"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleProtocol = *v.(*string)
	}

	if v := c.Arguments["firewall_rule_ip_address_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleIPAddressType = *v.(*string)
	}

	if v := c.Arguments["firewall_rule_port"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRulePortRangeStart, fw.FirewallRulePortRangeEnd, err = portStringToRange(*v.(*string))
		if err != nil {
			return "", err
		}
	}

	if v := c.Arguments["firewall_rule_source_ip_address"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleSourceIPAddressRangeStart, fw.FirewallRuleSourceIPAddressRangeEnd, err = addressStringToRange(*v.(*string))
		if err != nil {
			return "", err
		}
	}

	if v := c.Arguments["firewall_rule_desination_ip_address"]; v != nil && *v.(*string) != _nilDefaultStr {
		fw.FirewallRuleDestinationIPAddressRangeStart, fw.FirewallRuleDestinationIPAddressRangeEnd, err = addressStringToRange(*v.(*string))
		if err != nil {
			return "", err
		}
	}

	newFW := []metalcloud.FirewallRule{}
	found := false
	for _, f := range retIA.InstanceArrayOperation.InstanceArrayFirewallRules {
		if !fwRulesEqual(f, fw) {
			newFW = append(newFW, f)
		} else {
			found = true
		}
	}

	if !found {
		return "", fmt.Errorf("No matching firewall rule was found %v", fw)
	}

	retIA.InstanceArrayOperation.InstanceArrayFirewallRules = newFW
	bFalse := false
	_, err = client.InstanceArrayEdit(retIA.InstanceArrayID, *retIA.InstanceArrayOperation, &bFalse, nil, nil, nil)

	return "", err
}

func fwRulesEqual(a, b metalcloud.FirewallRule) bool {
	return a.FirewallRuleProtocol == b.FirewallRuleProtocol &&
		a.FirewallRulePortRangeStart == b.FirewallRulePortRangeStart &&
		a.FirewallRulePortRangeEnd == b.FirewallRulePortRangeEnd &&
		a.FirewallRuleSourceIPAddressRangeStart == b.FirewallRuleSourceIPAddressRangeStart &&
		a.FirewallRuleSourceIPAddressRangeEnd == b.FirewallRuleSourceIPAddressRangeEnd &&
		a.FirewallRuleDestinationIPAddressRangeStart == b.FirewallRuleDestinationIPAddressRangeStart &&
		a.FirewallRuleDestinationIPAddressRangeEnd == b.FirewallRuleDestinationIPAddressRangeEnd
}

func portStringToRange(s string) (int, int, error) {
	port, err := strconv.Atoi(s)

	if err == nil && port > 0 {
		return port, port, nil
	}

	re := regexp.MustCompile(`^(\d+)\-(\d+)$`)
	matches := re.FindStringSubmatch(s)

	if matches == nil {
		return 0, 0, fmt.Errorf("Could not parse port definition %s", s)
	}

	startPort, err := strconv.Atoi(matches[1])

	if err != nil && startPort > 0 {
		return 0, 0, fmt.Errorf("Could not parse port definition %s", s)
	}

	endPort, err := strconv.Atoi(matches[2])
	if err != nil && endPort > 0 {
		return 0, 0, fmt.Errorf("Could not parse port definition %s", s)
	}

	return startPort, endPort, nil

}

func addressStringToRange(s string) (string, string, error) {

	if s == "" {
		return "", "", fmt.Errorf("address cannot be empty")
	}

	components := strings.Split(s, "-")

	if len(components) == 1 && components[0] != "" {
		return s, s, nil //single address, we return it
	}

	if len(components) != 2 || components[0] == "" || components[1] == "" {
		return "", "", fmt.Errorf("cannot parse address %s", s)
	}

	return components[0], components[1], nil

}
