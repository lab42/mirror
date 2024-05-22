package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
		// Create clientset
		clientset, err := kubernetes.NewClientSet()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

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

		// Create a context with cancel function
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(2)

		// Function to run reflector and handle errors
		runReflector := func(r *reflector.Reflector, resourceType string) {
			defer wg.Done()
			log.Info().Str("type", resourceType).Msg("starting reflector")
			if err := r.Run(ctx); err != nil {
				log.Fatal().Err(err).Str("type", resourceType).Msg("error running reflector")
				cancel() // Cancel the context to stop the other reflector
			}
		}

		// Run the reflectors in separate goroutines
		go runReflector(cmReflector, "configMap")
		go runReflector(secretReflector, "secret")

		// Setup signal handling to gracefully shutdown on SIGINT or SIGTERM
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for a termination signal
		go func() {
			sig := <-sigChan
			log.Info().Str("signal", sig.String()).Msg("received signal")
			log.Info().Msg("initiating shutdown")
			cancel() // Cancel the context to stop the reflectors
		}()

		// Wait for the reflectors to finish
		wg.Wait()
		log.Info().Msg("reflectors exited")
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
