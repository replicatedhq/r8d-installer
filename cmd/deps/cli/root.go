package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "r8d-deps",
	Short: "Utility application to manage r8d-installer dependencies",
	Long: `This is a utility application to manage r8d-installer dependencies.

Intended usage:
1. Run "r8d-deps build <manifest.toml>" to generate the download, build, and 
   package assets from the manifest file for the installer build.
2. Run "r8d-deps update <manifest.toml>" to update the manifest file with the 
   latest dependency versions.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.r8d-deps.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".r8d-deps" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".r8d-deps")
	}

	viper.SetEnvPrefix("R8D")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
