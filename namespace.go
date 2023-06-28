/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetK8sNamespaces(kubeConfig *rest.Config) []string {

	var allNamespaces []string

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < len(namespaces.Items); i++ {
		allNamespaces = append(allNamespaces, namespaces.Items[i].Name)
	}

	return allNamespaces
}

func CreateNamespace(kubeConfig *rest.Config, namespaceName string) {

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}

	nsName := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	clientset.CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})

}
