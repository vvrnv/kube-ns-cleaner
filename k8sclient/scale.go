package k8sclient

import (
	"context"

	log "github.com/rs/zerolog/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ScaleStatefulSet(clientset *kubernetes.Clientset, namespace string, replicas int32) {
	statefulSetsClient := clientset.AppsV1().StatefulSets(namespace)
	listOptions := metav1.ListOptions{}
	statefulSets, err := statefulSetsClient.List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal().Msgf("Failed to get StatefulSets: %v", err)
	}
	for _, statefulSet := range statefulSets.Items {
		statefulSet.Spec.Replicas = &replicas
		_, err := statefulSetsClient.Update(context.TODO(), &statefulSet, metav1.UpdateOptions{})
		if err != nil {
			log.Fatal().Msgf("Failed to scale down StatefulSet '%s': %v", statefulSet.Name, err)
		}
	}
}

func ScaleDeployment(clientset *kubernetes.Clientset, namespace string, replicas int32) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	listOptions := metav1.ListOptions{}
	deployments, err := deploymentsClient.List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal().Msgf("Failed to get Deployments: %v", err)
	}
	for _, deployment := range deployments.Items {
		deployment.Spec.Replicas = &replicas
		_, err := deploymentsClient.Update(context.TODO(), &deployment, metav1.UpdateOptions{})
		if err != nil {
			log.Fatal().Msgf("Failed to scale down Deployment '%s': %v", deployment.Name, err)
		}
	}
}

// delete daemonset because cannot be scaled to zero
// a patch of resources is possible, but to return it back, you need to patch the resource again
func ScaleDaemonSet(clientset *kubernetes.Clientset, namespace string, replicas int32) {
	daemonSetsClient := clientset.AppsV1().DaemonSets(namespace)
	listOptions := metav1.ListOptions{}
	daemonSets, err := daemonSetsClient.List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal().Msgf("Failed to get DaemonSets: %v", err)
	}
	for _, daemonSet := range daemonSets.Items {
		err := daemonSetsClient.Delete(context.TODO(), daemonSet.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Fatal().Msgf("Failed to delete DaemonSet '%s': %v", daemonSet.Name, err)
		}
	}
}
