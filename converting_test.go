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
