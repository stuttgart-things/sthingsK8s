/*
Copyright © 2022 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"io/ioutil"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ConvertKubeconfigFileToByteArray(kubeconfigPath string) []byte {

	kubeconfigByte, err := ioutil.ReadFile(kubeconfigPath)

	if err != nil {
		panic(err)
	}

	return kubeconfigByte
}

func CreateRestConfig(kubeconfig []byte) *rest.Config {

	// create clientcmd.clientconfig
	config, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create restconfig
	restconfig, err := config.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	return restconfig
}

func CreateInClusterRestConfig() *rest.Config {

	restconfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	return restconfig
}

func GetKubeConfig(kubeConfigPath string) (clusterConfig *rest.Config, clusterConnection string) {

	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
		clusterConnection = "inside"
		clusterConfig = CreateInClusterRestConfig()

	} else {
		clusterConnection = "outside"
		clusterConfig = CreateRestConfig([]byte(sthingsBase.ReadFileToVariable(kubeConfigPath)))
	}

	return
}
