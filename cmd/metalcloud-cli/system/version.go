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
const maxMinor = 3

var AllowDevelop bool

func ValidateVersion(ctx context.Context, cliVersion string) error {
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
		return fmt.Errorf("invalid version: %s (CLI version: %s)", version.Version, cliVersion)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("invalid version: %s (CLI version: %s)", version.Version, cliVersion)
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return fmt.Errorf("invalid version: %s (CLI version: %s)", version.Version, cliVersion)
	}

	if major < minMajor || major > maxMajor || (major == minMajor && minor < minMinor) || (major == maxMajor && minor > maxMinor) {
		return fmt.Errorf("incompatible version: server is %s, CLI version is %s (compatible range: %d.%d - %d.%d)", version.Version, cliVersion, minMajor, minMinor, maxMajor, maxMinor)
	}

	return nil
}

func GetMinMaxVersion() (string, string) {
	return fmt.Sprintf("%d.%d", minMajor, minMinor), fmt.Sprintf("%d.%d", maxMajor, maxMinor)
}
