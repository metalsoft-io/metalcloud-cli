package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const allowDevelop = true
const minMajor = 7
const minMinor = 0

func ValidateVersion(ctx context.Context) error {
	client, err := GetApiClient(ctx)
	if err != nil {
		return err
	}

	version, _, err := client.SystemAPI.GetVersion(ctx).Execute()
	if err != nil {
		return fmt.Errorf("failed to get version: %v", err)
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
