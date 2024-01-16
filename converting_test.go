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

	yamlPipelineRun = `
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: simulate-stagetime-pipelinerun-25
  namespace: tektoncd
spec:
  timeouts:
    tasks: 20s
    finally: 10s
    pipeline: 30s
  params:
    - name: gitRepoUrl
      default: 'https://github.com/stuttgart-things/stageTime-server.git'
    - name: gitRevision
      default: main
    - name: gitWorkspaceSubdirectory
      default: stageTime
    - name: scriptPath
      default: tests/prime.sh
    - name: scriptTimeout
      default: "15s"
  taskRunTemplate:
    podTemplate:
      securityContext:
        fsGroup: 65532
  workspaces:
    - name: source
      volumeClaimTemplate:
        spec:
          storageClassName: openebs-hostpath
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 20Mi
  pipelineRef:
    resolver: git
    params:
      - name: url
        value: "https://github.com/stuttgart-things/stuttgart-things.git"
      - name: revision
        value: main
      - name: pathInRepo
        value: stageTime/pipelines/simulate-stagetime-pipelineruns.yaml
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

func TestConvertYAMLtoObject(t *testing.T) {
	assert := assert.New(t)

	validDeployment, _ := ConvertYAMLtoDeployment(yamlDeployment)
	assert.Equal(validDeployment, true)
}

func TestConvertYAMLtoPipelineRun(t *testing.T) {
	assert := assert.New(t)

	validPipelineRun, _, err := ConvertYAMLtoPipelineRun(yamlPipelineRun)
	fmt.Println(err)
	assert.Equal(validPipelineRun, true)
}
