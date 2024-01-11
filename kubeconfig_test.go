/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package k8s

import (
	"fmt"
	"testing"
)

func TestGetKubeConfig(t *testing.T) {
	clusterConfig, _ := GetKubeConfig("/home/sthings/.kube/pve-cd43")

	ns := GetK8sNamespaces(clusterConfig)
	fmt.Println(ns)
}
