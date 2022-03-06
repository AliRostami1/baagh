package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	overwriteConfigFile string
	daemonAddr          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "baagh-cli",
	Short: "baagh-cli is the interface to talk to baaghd",
	Long: `baagh-cli makes it easy to manipulate GPIO pins 
	using the baaghd (baagh daemon) application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&overwriteConfigFile, "config", "", "config file (default is $XDG_CONFIG_HOME/baagh/cli | $HOME/.config/baagh/cli)")
	rootCmd.PersistentFlags().StringVarP(&daemonAddr, "daemon", "d", "http://127.0.0.1:8080", "URL that daemon is listening on, defaults to 127.0.0.1:8080")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if overwriteConfigFile == "" {
		// use default config directory if the passed
		// overwriteConfigPath is empty string
		userConfigPath, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "config path failed: %v", err)
		}
		// path whould be $XDG_CONFIG_HOME/baagh/cli | $HOME/.config/baagh/cli
		viper.AddConfigPath(path.Join(userConfigPath, "baagh"))
		// config name is cli without extentions
		viper.SetConfigName("cli")
		// file format is yaml | JSON | toml
		viper.SetConfigType("yaml")
		viper.SetConfigType("json")
		viper.SetConfigType("toml")
	} else {
		// but if it's not empty then set the path
		// for the config file
		viper.SetConfigFile(overwriteConfigFile)
	}

	viper.SetEnvPrefix("BAAGH")
	viper.AutomaticEnv()

	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "config path failed: %v", err)
		}
	}
}
