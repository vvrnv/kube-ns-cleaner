package k8sclient

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vvrnv/kube-ns-cleaner/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func K8sClient() {
	/* LOCAL DEBUG CODE
	// connect to k8s outside for local debugging
	kubeconfig := "/Users/voronov2/.kube/config"
	k8sConfing, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// get current k8s context
	checkConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	// check current context
	currentContext := checkConfig.CurrentContext
	if currentContext != "docker-desktop" && currentContext != "dev" {
		log.Fatal().Msgf("Can't work with %s context!", currentContext)
	}
	*/ //LOCAL DEBUG CODE

	// in cluster config
	k8sConfing, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create a kubernetes client
	clientset, err := kubernetes.NewForConfig(k8sConfing)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// get list of namespaces
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	excludedNamespaces := config.Config.ExcludedNamespaces
	scalingLifeTime := config.Config.ScalingLifeTime
	deletingLifeTime := config.Config.DeleteingLifeTime

	// print ns list with active time in hours
	for _, ns := range namespaceList.Items {
		if contains(excludedNamespaces, ns.Name) {
			continue
		}
		duration := time.Since(ns.CreationTimestamp.Time)
		hours := int(duration.Hours())
		// scaling namespaces
		if scalingLifeTime <= hours && hours < deletingLifeTime {
			log.Info().Msgf("Namespace: %s, Active: %d hours, Action: will be scaled to 0", ns.Name, hours)
			scaleResources := []string{"statefulset", "deployment", "daemonset"}
			for _, resourceType := range scaleResources {
				switch resourceType {
				case "statefulset":
					ScaleStatefulSet(clientset, ns.Name, 0)
				case "deployment":
					ScaleDeployment(clientset, ns.Name, 0)
				case "daemonset":
					ScaleDaemonSet(clientset, ns.Name, 0)
				}
			}
			log.Info().Msgf("Resources of %s namespace successfully scaled down to 0 replicas.", ns.Name)
		}
		// deleting namespaces
		if hours >= deletingLifeTime {
			log.Info().Msgf("Namespace: %s, Active: %d hours, Action: will be deleted", ns.Name, hours)
			err = clientset.CoreV1().Namespaces().Delete(context.TODO(), ns.Name, metav1.DeleteOptions{})
			if err != nil {
				log.Error().Msgf("Failed to delete namespace: %v", err)
			}
			log.Info().Msgf("Namespace '%s' deleted successfully", ns.Name)
		}
	}
}

// func for checking namespaces
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
