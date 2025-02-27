/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// cordonNode marks the node as unschedulable
func CordonNode(clientset *kubernetes.Clientset, nodeName string) error {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	node.Spec.Unschedulable = true // Mark node as unschedulable

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Node %s is now cordoned\n", nodeName)
	return nil
}

func UncordonNode(clientset *kubernetes.Clientset, nodeName string) error {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	node.Spec.Unschedulable = false // Mark node as schedulable

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Node %s is now uncordoned (schedulable again)\n", nodeName)
	return nil
}
