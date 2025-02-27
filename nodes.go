/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetNodesByRole fetches nodes by role from the cluster
func GetNodesByRole(clientset *kubernetes.Clientset, role string) ([]v1.Node, error) {
	labelSelector := fmt.Sprintf("node-role.kubernetes.io/%s", role)
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}
