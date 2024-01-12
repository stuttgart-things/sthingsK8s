/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/batch/v1"
)

func VerifyYamlJobDefinition(jobManifest string) (bool, error) {

	job := &v1.Job{}
	err := yaml.Unmarshal([]byte(jobManifest), job)
	if err != nil {
		return false, err
	}

	return true, nil
}
