/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reflect",
	Short: "Secrets and config reflector for kubernetes.",
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(
		rootCmd.Execute(),
	)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reflect.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".reflect" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".reflect")
	}

	viper.SetDefault("logLevel", "debug")

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("path", viper.ConfigFileUsed()).Msg("configuration loaded")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Warn().Str("path", viper.ConfigFileUsed()).Msg("configuration changed")
		if err := viper.ReadInConfig(); err == nil {
			log.Info().Str("path", viper.ConfigFileUsed()).Msg("configuration reloaded")
		}
	})
}
