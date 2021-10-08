package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/spf13/cobra"
)

var (
	cfg       *Config
	awsConfig *Outback.AwsConfig
)

// Flags
var (
	flagCluster    string
	flagService    string
	flagConfigName string
	flagTimeout    int
)

// RootCmd represents the base command when called
// without any subcommands
var rootCmd = &cobra.Command{
	Use:     "outback",
	Short:   "Outback is an ecs deployment tool",
	Long:    ``,
	Version: "v2.0.0",
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
	rootCmd.PersistentFlags().StringVar(&flagConfigName, "config", "config", "outback config name")
	rootCmd.PersistentFlags().IntVarP(&flagTimeout, "timeout", "t", 5, "Deployment Timeout Time")
}

func loadConfig() {
	cwd, err := os.Getwd()

	if flagConfigName != "" {
		viper.AddConfigPath(cwd + "/.outback")
		viper.AddConfigPath(cwd + "/.ufo")
		viper.SetConfigName(flagConfigName)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("outback config not found (searched for both .config and .ufo)")
		}

		handleError(err)

		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Printf("Unable to unmarshal config, %v", err)
			os.Exit(1)
		}
	}

	awsConfig = &Outback.AwsConfig{
		Profile: cfg.Profile,
		Region:  cfg.Region,
	}
}
