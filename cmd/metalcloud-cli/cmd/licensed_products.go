package cmd

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/spf13/cobra"
)

var licensedProductsCmd = &cobra.Command{
	Use:          "licensed-products",
	Aliases:      []string{"license", "licenses"},
	Short:        "Show licensed product categories",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getLicensedProducts(cmd.Context())
	},
}

func getLicensedProducts(ctx context.Context) error {
	logger.Get().Info().Msg("Getting licensed products")

	client := api.GetApiClient(ctx)

	products, httpRes, err := client.SystemAPI.GetLicensedProducts(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Servers:  %s\n", licensedStr(products.ServersAreLicensed))
	fmt.Printf("Switches: %s\n", licensedStr(products.SwitchesAreLicensed))
	fmt.Printf("VMs:      %s\n", licensedStr(products.VmsAreLicensed))
	fmt.Printf("Storage:  %s\n", licensedStr(products.StoragesAreLicensed))
	return nil
}

func licensedStr(licensed bool) string {
	if licensed {
		return "licensed"
	}
	return "not licensed"
}

func init() {
	rootCmd.AddCommand(licensedProductsCmd)
}
