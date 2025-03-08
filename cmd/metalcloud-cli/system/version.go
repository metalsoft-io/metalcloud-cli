package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

const allowDevelop = true
const minMajor = 7
const minMinor = 0

func ValidateVersion(ctx context.Context) error {
	client := api.GetApiClient(ctx)

	version, httpRes, err := client.SystemAPI.GetVersion(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if allowDevelop && version.Version == "develop" {
		return nil
	}

	versionParts := strings.Split(version.Version, ".")
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

	if major < minMajor || (major == minMajor && minor < minMinor) {
		return fmt.Errorf("incompatible version: %s", version.Version)
	}

	return nil
}
