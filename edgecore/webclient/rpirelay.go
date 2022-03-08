package main

import (
	"encoding/json"
	"context"
	"fmt"
	
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

type deviceItem struct {
    Device string
	Model string
}

type deviceBox struct {
    Items []deviceItem
}

func (box *deviceBox) AddItem(item deviceItem) []deviceItem {
	box.Items = append(box.Items, item)
    return box.Items
}

func getDevices(client dynamic.Interface) deviceBox  {
	//  Create a GVR which represents an Istio Virtual Service.
	relayServiceGVR := schema.GroupVersionResource{
		Group:    "devices.kubeedge.io",
        Version:  "v1alpha2",
        Resource: "devices",
	}
	//var result = make(map[string]string)
	result := deviceBox{}

	listOptions := metav1.ListOptions{}

	//  List all of the Virtual Services.
    devices, err := client.Resource(relayServiceGVR).Namespace("default").List(context.TODO(), listOptions)
    
    if err != nil {
		panic(err.Error())
	}	

    for _, device := range devices.Items {
		var deviceModel = device.UnstructuredContent()["spec"].(map[string]interface{})["deviceModelRef"].(map[string]interface{})["name"].(string)
		fmt.Printf("\nDevice: %s DeviceModel: %s \n", device.GetName() , deviceModel )

		item := deviceItem{ Device: device.GetName(), Model: deviceModel }  
		result.AddItem( item )
    }

	return result
}

func getDevice(client dynamic.Interface, deviceName string) (string, string)  {
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
		var deviceModel = device.UnstructuredContent()["spec"].(map[string]interface{})["deviceModelRef"].(map[string]interface{})["name"].(string)
		fmt.Printf("\nDevice: %s DeviceModel: %s \n", device.GetName() , deviceModel )

		return device.GetName(), deviceModel
    }

	return "", ""
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

func getDHTState(client dynamic.Interface, deviceName string) (string, string, string )  {
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
		temperature := twins.([]interface{})[0].(map[string]interface{})["reported"].(map[string]interface{})["value"].(string)
		humidity := twins.([]interface{})[1].(map[string]interface{})["reported"].(map[string]interface{})["value"].(string)
		
		return deviceName, temperature, humidity
    }

	return "", "", ""
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