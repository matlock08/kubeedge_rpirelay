package main

import (
	"encoding/json"
	"context"
	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/dynamic"
	
)

var (

	//  Name of the VirtualService to weight, and the two weight values.
	deviceName = "relay-instance-01"
	ch1            = string(50)
	ch2            = string(50)
	ch3            = string(50)
)

//  patchStringValue specifies a patch operation for a string.
type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func getRelays(client dynamic.Interface) []string  {
	//  Create a GVR which represents an Istio Virtual Service.
	relayServiceGVR := schema.GroupVersionResource{
		Group:    "devices.kubeedge.io",
        Version:  "v1alpha2",
        Resource: "devices",
	}
	var result []string

	listOptions := metav1.ListOptions{}

	//  List all of the Virtual Services.
    devices, err := client.Resource(relayServiceGVR).Namespace("default").List(context.TODO(), listOptions)
    
    if err != nil {
		panic(err.Error())
	}	

    for _, device := range devices.Items {
		
		result = append(result, device.GetName() )
		
    }

	return result
}

func getRelayState(client dynamic.Interface, deviceName string) (string, string, string, string)  {
	//  Create a GVR which represents an Istio Virtual Service.
	relayServiceGVR := schema.GroupVersionResource{
		Group:    "devices.kubeedge.io",
        Version:  "v1alpha2",
        Resource: "devices",
	}

	listOptions := metav1.ListOptions{
		FieldSelector: "metadata.name=" + deviceName,
	}

	//  List all of the Virtual Services.
    devices, err := client.Resource(relayServiceGVR).Namespace("default").List(context.TODO(), listOptions)
    
    if err != nil {
		panic(err.Error())
	}	

    for _, device := range devices.Items {
		twins := device.UnstructuredContent()["status"].(map[string]interface{})["twins"]
		deviceName := device.GetName()
		ch1Value := twins.([]interface{})[0].(map[string]interface{})["desired"].(map[string]interface{})["value"].(string)
		ch2Value := twins.([]interface{})[1].(map[string]interface{})["desired"].(map[string]interface{})["value"].(string)
		ch3Value := twins.([]interface{})[2].(map[string]interface{})["desired"].(map[string]interface{})["value"].(string)
		//fmt.Printf("\nDevice: %s\n[CH1]: %s\n[CH2]: %s\n[CH3]: %s\n", device.GetName(), ch1Value, ch2Value, ch3Value)
		return deviceName, ch1Value, ch2Value, ch3Value
    }

	return "", "", "", ""
}

func setRelayState(client dynamic.Interface, deviceName string, ch1 string, ch2 string, ch3 string) error {
	//  Create a GVR which represents an Istio Virtual Service.
	virtualServiceGVR := schema.GroupVersionResource{
		Group:    "devices.kubeedge.io",
        Version:  "v1alpha2",
        Resource: "devices",
	}

	//  Weight the two routes - 50/50.
	patchPayload := make([]patchStringValue, 3)
	patchPayload[0].Op = "replace"
	patchPayload[0].Path = "/status/twins/0/desired/value"
	patchPayload[0].Value = ch1

	patchPayload[1].Op = "replace"
	patchPayload[1].Path = "/status/twins/1/desired/value"
	patchPayload[1].Value = ch2

	patchPayload[2].Op = "replace"
	patchPayload[2].Path = "/status/twins/2/desired/value"
	patchPayload[2].Value = ch3


	patchBytes, _ := json.Marshal(patchPayload)

	//  Apply the patch to the 'service2' service.
	_, err := client.Resource(virtualServiceGVR).Namespace("default").Patch(context.TODO(), deviceName, types.JSONPatchType, patchBytes, metav1.PatchOptions{})

	return err
}
/*
func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

    dynamicClient, _ := dynamic.NewForConfig(config)

	deviceName, ch1Value, ch2Value, ch3Value := getRelayState(dynamicClient, "relay-instance-01")

	fmt.Printf("\nDevice: %s\n[CH1]: %s\n[CH2]: %s\n[CH3]: %s\n", deviceName, ch1Value, ch2Value, ch3Value)

	err = setRelayState(dynamicClient, "relay-instance-01", "OFF", "OFF", "OFF")
	if err != nil {
		panic(err.Error())
	}

	
}
*/