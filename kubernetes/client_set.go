package kubernetes

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func NewClientSet() (*kubernetes.Clientset, error) {
	if viper.GetBool("kubeconfig.inCluster") {
		return NewInclusterClientSet()
	}

	if !viper.GetBool("kubeconfig.inCluster") {
		return NewOutOfClusterClientSet()
	}

	return nil, errors.New("no configuration provided")
}

func NewOutOfClusterClientSet() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig.path"))
	if err != nil {
		return nil, err
	}

	// return the clientset
	return kubernetes.NewForConfig(config)
}

func NewInclusterClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	// return the clientset
	return kubernetes.NewForConfig(config)
}
