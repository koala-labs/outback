package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

var (
	cfg    *Config
	ufoCfg UFO.Config
)

// Flags
var (
	flagCluster    string
	flagService    string
	flagConfigName string
)

// RootCmd represents the base command when called
// without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ufo",
	Short:   "Ufo is an ecs deployment tool",
	Long:    ``,
	Version: "0.18.5.30",
}

// Execute adds all child commands so the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)
	// Here you will define your flags and configuration settings
	// Cobra supports Persistent Flags which if defined here will be global for your application

	rootCmd.PersistentFlags().StringVarP(&flagCluster, "cluster", "c", "", "AWS ECS Cluster")
	rootCmd.PersistentFlags().StringVarP(&flagService, "service", "s", "", "Service in an ECS cluster")
	rootCmd.Flags().StringVar(&flagConfigName, "config", "config", "ufo config name")
}

func loadConfig() {
	cwd, err := os.Getwd()

	if flagConfigName != "" {
		viper.AddConfigPath(cwd + "/.ufo")
		viper.SetConfigName(flagConfigName)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("ufo config not found")
		}

		handleError(err)

		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Printf("Unable to decode config, %v", err)
			os.Exit(1)
		}
	}

	ufoCfg = UFO.Config{
		Region:  &cfg.Region,
		Profile: &cfg.Profile,
	}
}
