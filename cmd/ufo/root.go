package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

// Config
const (
	configPath = "/.ufo/config.json"
	configDir  = "/.ufo/"
	configFile = "/config.json"
)

var (
	cfg    *Config
	ufoCfg UFO.Config
)

// Flags
var (
	flagCluster string
	flagService string
)

// RootCmd represents the base command when called
// without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ufo",
	Short:   "Ufo is an ecs deployment tool",
	Long:    ``,
	Version: "18.4.1",
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
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings
	// Cobra supports Persistent Flags which if defined here will be global for your application

	rootCmd.PersistentFlags().StringVarP(&flagCluster, "cluster", "c", "", "AWS ECS Cluster")
	rootCmd.PersistentFlags().StringVarP(&flagService, "service", "s", "", "Service in an ECS cluster")
	rootCmd.MarkPersistentFlagRequired("cluster")
}

func initConfig() {
	cwd, err := os.Getwd()

	handleError(err)

	viper.AddConfigPath(cwd + "/.ufo")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("ufo config not found...")
		fmt.Println("Initializing ufo config.")

		createDirectory(filepath.Join(cwd, configDir))

		f, err := createConfig(filepath.Join(cwd, configPath))

		handleError(err)

		defer f.Close()

		fmt.Fprint(f, configTemplate)

		fmt.Println("ufo config initialized")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("Unable to decode config into struct, %v", err)
		os.Exit(1)
	}

	ufoCfg = UFO.Config{
		Region:  &cfg.Region,
		Profile: &cfg.Profile,
	}
}
