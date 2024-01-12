/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"fmt"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"gopkg.in/yaml.v2"
	v1Apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"

	"k8s.io/client-go/kubernetes/scheme"
)

func VerifyYamlJobDefinition(jobManifest string) (bool, error) {

	job := &batch.Job{}
	err := yaml.Unmarshal([]byte(jobManifest), job)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ConvertYAMLtoDeployment(yamlString string) (bool, *v1Apps.Deployment) {

	fmt.Println(yamlString)
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(yamlString), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	fmt.Println(obj)

	deployment := obj.(*v1Apps.Deployment)
	fmt.Printf("%#v\n", deployment)

	return true, deployment
}

func ConvertYAMLtoPipelineRun(yamlString string) (bool, *v1.PipelineRun) {

	pipelienRun := &v1.PipelineRun{}
	fmt.Println(yamlString)

	err := yaml.Unmarshal([]byte(yamlString), pipelienRun)
	if err != nil {
		fmt.Printf("%#v", err)
		return false, nil
	}
	fmt.Printf("%#v\n", pipelienRun)

	return true, pipelienRun
}
