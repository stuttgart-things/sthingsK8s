/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"fmt"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/batch/v1"

	v1Apps "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func VerifyYamlJobDefinition(jobManifest string) (bool, error) {

	job := &v1.Job{}
	err := yaml.Unmarshal([]byte(jobManifest), job)
	if err != nil {
		return false, err
	}

	return true, nil
}

func YAMLtoObject(yamlString string) {

	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(yamlString), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	deployment := obj.(*v1Apps.Deployment)

	fmt.Printf("%#v\n", deployment)
}
