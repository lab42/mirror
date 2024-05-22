/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/lab42/mirror/kubernetes"
	"github.com/lab42/mirror/reflector"
	"github.com/lab42/mirror/resource"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reflect",
	Short: "Secrets and config reflector for kubernetes.",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := kubernetes.NewClientSet()
		cobra.CheckErr(err)

		// Create resources
		cmResource := &resource.ConfigMapResource{Clientset: clientset}
		secretResource := &resource.SecretResource{Clientset: clientset}

		// Create a Reflector for ConfigMaps
		cmReflectorConfig := reflector.ReflectorConfig{
			ClientSet:    clientset,
			Resource:     cmResource,
			ResourceType: reflector.ConfigMap,
		}
		cmReflector := reflector.NewReflector(cmReflectorConfig)

		// Create a Reflector for Secrets
		secretReflectorConfig := reflector.ReflectorConfig{
			ClientSet:    clientset,
			Resource:     secretResource,
			ResourceType: reflector.Secret,
		}
		secretReflector := reflector.NewReflector(secretReflectorConfig)

		// Run the reflectors
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			log.Info().Str("type", "configMap").Msg("starting reflector")
			defer wg.Done()
			if err := cmReflector.Run(ctx); err != nil {
				log.Fatal().Err(err).Str("type", "configMap").Msg("error running reflector")
				cancel() // Cancel the context to stop the other reflector
			}
		}()

		go func() {
			log.Info().Str("type", "secret").Msg("starting reflector")
			defer wg.Done()
			if err := secretReflector.Run(ctx); err != nil {
				log.Fatal().Err(err).Str("type", "secret").Msg("error running reflector")
				cancel() // Cancel the context to stop the other reflector
			}
		}()

		// Wait for the reflectors to finish
		wg.Wait()
		log.Info().Msg("reflectors stopped running")
	},
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
		log.Info().Str("path", viper.ConfigFileUsed()).Str("type", "configuration").Msg("configuration loaded")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info().Str("path", viper.ConfigFileUsed()).Str("type", "configuration").Msg("configuration changed")
		if err := viper.ReadInConfig(); err == nil {
			log.Info().Str("path", viper.ConfigFileUsed()).Str("type", "configuration").Msg("configuration reloaded")
		}
	})
}
