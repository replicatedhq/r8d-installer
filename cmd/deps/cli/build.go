package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/replicatedhq/r8d-installer/pkg/deps"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Downloads, builds, and packages assets for installer into pkg/component/<name>",
	Long: `Downloads, builds, and packages assets for installer into pkg/component/<name> 
based on the provided configuration (manifest file, env variables, cli flags).

Types of assets:
1. YAML Manifests
2. Executables/Binaries
3. Container Images

Usage:
  r8d-deps build --config=<manifest.toml>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manifest := deps.Manifest{
			RKE2: viper.GetString("rke2"),
		}
		return deps.Build(manifest)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
