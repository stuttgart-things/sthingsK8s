/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package k8s

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDynamicResourcesFromTemplatet(t *testing.T) {
	assert := assert.New(t)

	namespace := "default"
	clusterConfig, _ := GetKubeConfig("/home/sthings/.kube/pve-cd43")
	validManifest := `apiVersion: v1
kind: ConfigMap
metadata:
  name: game-config-1
data:
  enemies: aliens
  lives: "5"
`

	invalidManifest := `apiVersion: v1
kind: ConfigMap
metadata:
  name: game-config-1
stringData:
  enemies: raiders
`

	fmt.Println(validManifest)
	fmt.Println(invalidManifest)

	resource1Created, err1 := CreateDynamicResourcesFromTemplate(clusterConfig, []byte(validManifest), namespace)
	fmt.Println(resource1Created, err1)
	assert.Equal(resource1Created, true)

	resource2Created, err2 := CreateDynamicResourcesFromTemplate(clusterConfig, []byte(invalidManifest), namespace)
	fmt.Println(resource2Created, err2)
	assert.Equal(resource2Created, false)
}
