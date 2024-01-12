/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package k8s

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validJobManifest = `apiVersion: batch/v1
kind: Job
metadata:
  name: node-app-job
spec:
  template:
    spec:
      containers:
        - name: node-app-job
          image: alpine
          command: ["echo", "Welcome to my Node app"]
      restartPolicy: Never
`

	invalidJobManifest = `apiVersion: batch/v1
kind: Job
metadata:
  name: node-app-job
spec:
  template:
    spec:
         containers:
        - name: node-app-job
           image: alpine
          command: ["echo", "Welcome to my Node app"]
       restartPolicy: Never
`

	yamlDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
    name: my-complex-app-my-complex-app
    labels:
    app: my-complex-app
    chart: my-complex-app-0.2.0
    release: my-complex-app
    heritage: Tiller
spec:
    replicas: 1
    template:
    metadata:
        labels:
        app: my-complex-app
        release: my-complex-app
    spec:
        containers:
        - name: my-complex-app
        image: "nginx:stable"
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
        livenessProbe:
        httpGet:
            path: /
            port: 80
        readinessProbe:
        httpGet:
            path: /
            port: 80
        resources: {}
`
)

func TestVerifyYamlJobDefinition(t *testing.T) {
	assert := assert.New(t)

	// TEST THE VALID YAML
	valid, _ := VerifyYamlJobDefinition(validJobManifest)
	assert.Equal(valid, true)

	// TEST THE INVALID YAML
	invalid, err := VerifyYamlJobDefinition(invalidJobManifest)
	fmt.Println(err)
	assert.Equal(invalid, false)

}

func TestYAMLtoObject(t *testing.T) {
	YAMLtoObject(yamlDeployment)

}
