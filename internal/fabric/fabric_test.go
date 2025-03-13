package fabric

import (
	"testing"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestUnmarshaling(t *testing.T) {
	data := []byte(`{
    "fabricType": "ethernet",
    "defaultNetworkProfileId": 1,
    "gnmiMonitoringEnabled": false,
    "syslogMonitoringEnabled": true,
    "zeroTouchEnabled": false,
    "allocateDefaultVlan": true,
    "asnRanges": ["65000-65010"],
    "defaultVlan": 10,
    "extraInternalIPsPerSubnet": 2,
    "lagRanges": ["100-200", "300-400"],
    "leafSwitchesHaveMlagPairs": false,
    "mlagRanges": ["30-40", "50-60"],
    "numberOfSpinesNextToLeafSwitches": 5,
    "preventVlanCleanup": ["1000-1100"],
    "preventCleanupFromUplinks": true,
    "reservedVlans": ["2000-2100", "2200-2300"],
    "vlanRanges": ["3000-3100", "2000-2100"],
    "vniPrefix": 5000
  }`)

	dataWithExtraFields := []byte(`{
	    "fabricType": "ethernet",
	    "defaultNetworkProfileId": 1,
	    "gnmiMonitoringEnabled": false,
	    "syslogMonitoringEnabled": true,
	    "zeroTouchEnabled": false,
	    "allocateDefaultVlan": true,
	    "asnRanges": ["65000-65010"],
	    "defaultVlan": 10,
	    "extraInternalIPsPerSubnet": 2,
	    "lagRanges": ["100-200", "300-400"],
	    "leafSwitchesHaveMlagPairs": false,
	    "mlagRanges": ["30-40", "50-60"],
	    "numberOfSpinesNextToLeafSwitches": 5,
	    "preventVlanCleanup": ["1000-1100"],
	    "preventCleanupFromUplinks": true,
	    "reservedVlans": ["2000-2100", "2200-2300"],
	    "vlanRanges": ["3000-3100", "2000-2100"],
	    "vniPrefix": 5000,
	    "vrfVlanRanges": ["400-450", "460-470"],
	    "reservedVrfs": ["4000"],
	    "preventVrfCleanup": ["3000-3100", "2000-2100"]
	  }`)

	config := sdk.NetworkFabricFabricConfiguration{}

	err := config.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	err = config.UnmarshalJSON(dataWithExtraFields)
	if err == nil {
		t.Fatal("expected error")
	}
}
