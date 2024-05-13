/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lab42/reflect/kubernetes"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start reflector",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("starting reflector")
		cs, _ := kubernetes.NewClientSet()
		watcher, _ := cs.CoreV1().Secrets(v1.NamespaceAll).Watch(context.Background(), metav1.ListOptions{})

		watcher.ResultChan()
		for event := range watcher.ResultChan() {
			secret := event.Object.(*corev1.Secret)

			switch event.Type {
			case watch.Added:
				// fmt.Printf("\n%s: %s\n", secret.ObjectMeta.Name, secret.GetAnnotations())
				annotations := secret.GetAnnotations()
				val, ok := annotations["reflect.lab42.io/to-namespaces"]
				// If the key exists
				if ok && len(val) > 0 {
					fmt.Println(val)
					namespaces := strings.Split(val, ",")
					for _, namespace := range namespaces {
						if namespace == secret.Namespace {
							log.Error().Err(errors.New("cannot overwrite self")).Str("namespace", namespace).Msg("")
						}

						reflectedSecret := secret
						reflectedSecret.SetNamespace(secret.Namespace)
						reflectedSecret.Annotations["reflect.lab42.io/from-namespaces"] = secret.Namespace
						delete(reflectedSecret.Annotations, "reflect.lab42.io/to-namespaces")
					}
				}
			case watch.Modified:
				fmt.Printf("Service %s/%s modified", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)
			case watch.Deleted:
				fmt.Printf("Service %s/%s deleted", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
