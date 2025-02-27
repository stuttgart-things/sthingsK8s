/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package k8s

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// RestartAllDeployments restarts all deployments in every namespace
func RestartAllDeployments(clientset *kubernetes.Clientset) error {
	// Get all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ns := range namespaces.Items {
		namespace := ns.Name
		fmt.Printf("Restarting deployments in namespace: %s\n", namespace)

		// Get all deployments in the namespace
		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error getting deployments in namespace %s: %v", namespace, err)
			continue
		}

		for _, deploy := range deployments.Items {
			fmt.Printf("Restarting deployment: %s/%s\n", namespace, deploy.Name)

			// Add an annotation to trigger a rolling restart
			if deploy.Spec.Template.Annotations == nil {
				deploy.Spec.Template.Annotations = make(map[string]string)
			}
			deploy.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

			_, err := clientset.AppsV1().Deployments(namespace).Update(context.TODO(), &deploy, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("Failed to restart deployment %s/%s: %v", namespace, deploy.Name, err)
			}
		}
	}

	return nil
}

// evictPods evicts all non-DaemonSet and non-mirror pods from the node
func EvictPods(clientset *kubernetes.Clientset, nodeName string) error {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		// Skip DaemonSet pods
		if isDaemonSetPod(pod) {
			continue
		}

		// Create eviction
		eviction := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			DeleteOptions: &metav1.DeleteOptions{
				GracePeriodSeconds: int64Ptr(30), // Set a grace period
			},
		}

		err := clientset.PolicyV1().Evictions(eviction.Namespace).Evict(context.TODO(), eviction)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			log.Printf("Failed to evict pod %s: %v\n", pod.Name, err)
		} else {
			fmt.Printf("Evicted pod: %s/%s\n", pod.Namespace, pod.Name)
		}

		time.Sleep(1 * time.Second) // Avoid API throttling
	}

	return nil
}

// isDaemonSetPod checks if a pod is managed by a DaemonSet
func isDaemonSetPod(pod v1.Pod) bool {
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "DaemonSet" {
			return true
		}
	}
	return false
}

// Helper function to get pointer to int64
func int64Ptr(i int64) *int64 {
	return &i
}
