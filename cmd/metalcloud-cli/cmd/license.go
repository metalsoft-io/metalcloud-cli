package cmd

import (
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/license"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	licenseFlags = struct {
		source string
	}{}

	licenseCmd = &cobra.Command{
		Use:     "license [command]",
		Aliases: []string{"lic", "licenses"},
		Short:   "Manage the system license",
		Long: `Manage the license installed on the MetalSoft system.

A license controls which product categories are enabled (servers, switches, VMs,
storage) and the resource allowances granted to the installation.

Available Commands:
  status        Show the validity status of the current license
  get           Show the installed license document (Base64)
  allowance     Show the resource allowance granted by the license
  products      Show which product categories are licensed
  request       Show the license request document to send to MetalSoft
  add           Install a license on the system

Typical workflow:
  # 1. Generate the request document and send it to MetalSoft
  metalcloud-cli license request > license-request.txt

  # 2. Install the signed license returned by MetalSoft
  metalcloud-cli license add --source ./license.txt

  # 3. Verify it took effect
  metalcloud-cli license status`,
	}

	licenseStatusCmd = &cobra.Command{
		Use:          "status",
		Short:        "Show the validity status of the current license",
		Long:         `Show the validity status of the current license, including whether it is valid, its expiration date, and the hostname it is bound to.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LICENSES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return license.LicenseStatus(cmd.Context())
		},
	}

	licenseGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Show the installed license document",
		Long:         `Show the license document currently installed on the system as a Base64-encoded signed blob. In text output the raw document is printed so it can be piped or saved to a file.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LICENSES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return license.LicenseGet(cmd.Context())
		},
	}

	licenseAllowanceCmd = &cobra.Command{
		Use:          "allowance",
		Short:        "Show the resource allowance granted by the license",
		Long:         `Show the resource allowance granted by the current license: the number of servers, switches, and devices, and the amount of storage permitted.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LICENSES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return license.LicenseAllowance(cmd.Context())
		},
	}

	licenseProductsCmd = &cobra.Command{
		Use:          "products",
		Aliases:      []string{"licensed-products"},
		Short:        "Show which product categories are licensed",
		Long:         `Show which product categories (servers, switches, VMs, storage) are enabled by the current license.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LICENSES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return license.LicensedProducts(cmd.Context())
		},
	}

	licenseRequestCmd = &cobra.Command{
		Use:          "request",
		Short:        "Show the license request document",
		Long:         `Show the license request document for this system as a Base64-encoded blob. Send this document to MetalSoft to obtain a signed license. In text output the raw document is printed so it can be piped or saved to a file.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LICENSES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return license.LicenseRequest(cmd.Context())
		},
	}

	licenseAddCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"install", "create"},
		Short:   "Install a license on the system",
		Long: `Install a license on the system.

The license must be provided as a Base64-encoded signed document, exactly as
returned by MetalSoft. It is forwarded verbatim to preserve the original signed
bytes.

Required Flags:
  --source    Source of the license document. Can be 'pipe' to read from stdin
              or a path to a file containing the Base64-encoded license.

Examples:
  # Install a license from a file
  metalcloud-cli license add --source ./license.txt

  # Install a license from stdin
  cat license.txt | metalcloud-cli license add --source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LICENSES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := utils.ReadConfigFromPipeOrFile(licenseFlags.source)
			if err != nil {
				return err
			}

			licenseDoc := strings.TrimSpace(string(content))
			if licenseDoc == "" {
				return fmt.Errorf("license document is empty")
			}

			return license.LicenseAdd(cmd.Context(), licenseDoc)
		},
	}
)

func init() {
	rootCmd.AddCommand(licenseCmd)

	licenseCmd.AddCommand(licenseStatusCmd)
	licenseCmd.AddCommand(licenseGetCmd)
	licenseCmd.AddCommand(licenseAllowanceCmd)
	licenseCmd.AddCommand(licenseProductsCmd)
	licenseCmd.AddCommand(licenseRequestCmd)

	licenseCmd.AddCommand(licenseAddCmd)
	licenseAddCmd.Flags().StringVar(&licenseFlags.source, "source", "", "Source of the license document. Can be 'pipe' or a path to a file containing the Base64-encoded license.")
	licenseAddCmd.MarkFlagsOneRequired("source")
}
