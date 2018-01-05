package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	UFO "gitlab.fuzzhq.com/Web-Ops/ufo/ufo"
)

// Config
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
	Use:   "ufo",
	Short: "Ufo is an ecs deployment tool",
	Long: `A longer description that spans multiple
	lines and likely ... application`,
	Version: "0.0.1",
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

const configTemplate = `{
	"profile": "default",
	"region": "us-east-1",
	"repo": "default.dkr.ecr.us-west-1.amazonaws.com/default",
	"clusters": [
		{
			"name": "dev",
			"branch": "dev",
			"region": "us-west-1",
			"services": ["api", "queue"],
			"dockerfile": "Dockerfile.local"
		}
	],
	"tasks": [
		{
			"name": "migrate",
			"command": "php artisan migrate"
		}
	]
}
`

const gitIgnoreString = `
# UFO Config
.ufo/
`

const (
	configPath = "/.ufo/config.json"
	configDir  = "/.ufo/"
	configFile = "/config.json"
)

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

		updateGitIgnore(cwd)

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

func createDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating %s\n", path)
		os.Mkdir(path, os.ModePerm)
	}
}

func createConfig(path string) (*os.File, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil, ErrConfigFileAlreadyExists
	}

	fmt.Printf("Creating %s\n", path)

	f, err := os.Create(path)

	if err != nil {
		return nil, ErrCouldNotCreateConfig
	}

	return f, nil
}

func updateGitIgnore(path string) error {
	gitIgnore := filepath.Join(path, "/.gitignore")

	if _, err := os.Stat(gitIgnore); os.IsNotExist(err) {
		return ErrNoGitIgnore
	}

	// Open the file with read-write privileges so that we can check
	// if .ufo is already ignored and if not, write to .gitignore
	f, err := os.OpenFile(gitIgnore, os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		return ErrCouldNotOpenGitIgnore
	}

	defer f.Close()

	b := make([]byte, 1024)
	_, err = f.Read(b)

	fmt.Printf(".gitignore: %s", string(b))

	if strings.Contains(string(b), gitIgnoreString) {
		fmt.Println(".ufo is already ignored")
	} else {
		fmt.Println("Updating .gitignore")
		_, err = f.WriteString(gitIgnoreString)
	}

	return err
}
