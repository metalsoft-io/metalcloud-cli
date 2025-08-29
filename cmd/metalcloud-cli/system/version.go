package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

const minMajor = 7
const minMinor = 0
const maxMajor = 7
const maxMinor = 0

var AllowDevelop bool

func ValidateVersion(ctx context.Context) error {
	if AllowDevelop {
		return nil
	}

	client := api.GetApiClient(ctx)

	version, httpRes, err := client.SystemAPI.GetVersion(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	versionParts := strings.Split(strings.Trim(version.Version, "v "), ".")
	if len(versionParts) < 2 {
		return fmt.Errorf("invalid version: %s", version.Version)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("invalid version: %s", version.Version)
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return fmt.Errorf("invalid version: %s", version.Version)
	}

	if major < minMajor || major > maxMajor || (major == minMajor && minor < minMinor) || (major == maxMajor && minor > maxMinor) {
		return fmt.Errorf("incompatible version: %s", version.Version)
	}

	return nil
}

func GetMinMaxVersion() (string, string) {
	return fmt.Sprintf("%d.%d", minMajor, minMinor), fmt.Sprintf("%d.%d", maxMajor, maxMinor)
}
