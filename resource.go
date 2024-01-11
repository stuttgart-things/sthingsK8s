/*
Copyright © 2022 Patrick Hermann patrick.hermann@sva.de
*/

package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"

	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type K8sResource struct {
	Name       string `mapstructure:"name"`
	Kind       string `mapstructure:"kind"`
	ApiVersion string `mapstructure:"api-version"`
}

var prGroupResource = schema.GroupVersionResource{Group: "tekton.dev", Resource: "pipelineruns"}

func CreateRestMapperAndDynamicInterface(kubeconfig *rest.Config) (meta.RESTMapper, dynamic.Interface) {

	c, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := restmapper.GetAPIGroupResources(c.Discovery())
	if err != nil {
		log.Fatal(err)
	}

	mapper := restmapper.NewDiscoveryRESTMapper(resources)

	dynamicREST, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	return mapper, dynamicREST

}

func GetK8sResourcesByKind(kubeconfig *rest.Config, apiVersion, kind, namespace string, debug bool) []K8sResource {

	mapper, dynamicREST := CreateRestMapperAndDynamicInterface(kubeconfig)

	unstructuredObj := &unstructured.Unstructured{}
	unstructuredObj.SetAPIVersion(apiVersion)
	unstructuredObj.SetKind(kind)
	unstructuredObj.SetNamespace(namespace)

	resourceREST := getDynamicResourceInterface(unstructuredObj, mapper, dynamicREST)

	resourceList := getExistenceOfUnstructedObject(resourceREST)

	if debug {

		if len(resourceList) != 0 {
			fmt.Println("Found resources!")
		} else {
			fmt.Println("No resources found!!")
		}
	}

	return resourceList

}

func getDynamicResourceInterface(unstructuredObj *unstructured.Unstructured, mapper meta.RESTMapper, dynamicREST dynamic.Interface) dynamic.ResourceInterface {

	var resourceREST dynamic.ResourceInterface
	//unstructuredObj.SetNamespace("hello")

	mapping, err := mapper.RESTMapping(
		unstructuredObj.GetObjectKind().GroupVersionKind().GroupKind(),
		unstructuredObj.GetObjectKind().GroupVersionKind().Version)
	if err != nil {
		log.Fatal(err)
	}

	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {

		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace("default")
		}
		resourceREST =
			dynamicREST.
				Resource(mapping.Resource).
				Namespace(unstructuredObj.GetNamespace())
	} else {
		resourceREST = dynamicREST.Resource(mapping.Resource)
	}

	return resourceREST
}

func getExistenceOfUnstructedObject(resourceREST dynamic.ResourceInterface) []K8sResource {

	var resourceList []K8sResource

	unstructuredList, err := resourceREST.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range unstructuredList.Items {

		// fmt.Println(unstructuredList.Items[i])

		resourceList = append(resourceList, K8sResource{unstructuredList.Items[i].GetName(), unstructuredList.Items[i].GetKind(), unstructuredList.Items[i].GetAPIVersion()})

	}
	return resourceList

}

func CreateDynamicResourcesFromTemplate(kubeconfig *rest.Config, templatedResource []byte, namespace string) (bool, error) {

	mapper, dynamicREST := CreateRestMapperAndDynamicInterface(kubeconfig)
	objectsInYAML := bytes.Split(templatedResource, []byte("---"))

	for _, objectInYAML := range objectsInYAML {

		unstructuredObj, resourceREST := getUnstructedObjectAndDynamicResourceInterface(objectInYAML, mapper, dynamicREST, namespace)

		if len(getExistenceOfUnstructedObject(resourceREST)) < 0 {

			fmt.Println("CREATING RESOURCE ...")

			if _, err := resourceREST.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
				log.Fatal(err)
			} else {
				return true, err
			}

		} else {

			fmt.Println("PATCHING RESOURCE ...")

			data, err := json.Marshal(unstructuredObj)
			if err != nil {
				log.Fatal(err)
			}

			forceConflicts := false

			if _, err := resourceREST.Patch(context.Background(), unstructuredObj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
				FieldManager: "sthings-k8s",
				Force:        &forceConflicts,
			}); err != nil {
				return false, err
			}

		}
	}
	return true, nil

}

func getUnstructedObjectAndDynamicResourceInterface(objectInYAML []byte, mapper meta.RESTMapper, dynamicREST dynamic.Interface, namespace string) (*unstructured.Unstructured, dynamic.ResourceInterface) {

	runtimeObject, groupVersionAndKind, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(objectInYAML, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(groupVersionAndKind.GroupKind())
	fmt.Println(groupVersionAndKind.Version)

	unstructuredObj := runtimeObject.(*unstructured.Unstructured)
	unstructuredObj.SetNamespace(namespace)

	resourceREST := getDynamicResourceInterface(unstructuredObj, mapper, dynamicREST)

	return unstructuredObj, resourceREST
}
