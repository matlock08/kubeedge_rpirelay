package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/types"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/kubeedge/kubeedge/edge/pkg/devicetwin/dttype"
)

//var deviceID string = "traffic-light-instance-01"
var deviceID string
var mqtturl string
var modelName string

var ch1_wpi_num int64 = 25
var ch2_wpi_num int64 = 28
var ch3_wpi_num int64 = 29

var CONFIG_MAP_PATH = "/opt/kubeedge/deviceProfile.json"

const (
	DeviceETPrefix            = "$hw/events/device/"
	TwinETUpdateSuffix        = "/twin/update"
	TwinETUpdateDetalSuffix   = "/twin/update/delta"
	DeviceETStateUpdateSuffix = "/state/update"
	TwinETCloudSyncSuffix     = "/twin/cloud_updated"
	TwinETGetResultSuffix     = "/twin/get/result"
	TwinETGetSuffix           = "/twin/get"
)

const (
	CH1_STATE = "ch1"
	CH2_STATE = "ch2"
	CH3_STATE = "ch3"
)

const (
	CH1pinNumberConfig = "ch1-pin-number"
	CH2pinNumberConfig = "ch2-pin-number"
	CH3pinNumberConfig = "ch3-pin-number"
)

var Client MQTT.Client
var onceClient sync.Once

func parseFlag() {
	flag.StringVar(&deviceID, "device", "relay-instance-01", "device id name, default is relay-instance-01 ")
	flag.StringVar(&mqtturl, "mqtturl", "tcp://127.0.0.1:1883", "mqtt url default is tcp://127.0.0.1:1883")
	flag.StringVar(&modelName, "modelname", "relay-model", "device model name , default is relay-model")
	flag.Parse()
}

func loadConfigMap() error {

	readConfigMap := &types.DeviceProfile{}
	jsonFile, err := ioutil.ReadFile(CONFIG_MAP_PATH)
	if err != nil {
		log.Fatalf("readfile %v error %v\n", CONFIG_MAP_PATH, err)
		return err
	}
	err = json.Unmarshal(jsonFile, readConfigMap)
	if err != nil {
		log.Fatalf("unmarshal error %v", err)
		return err
	}

	for _, deviceModel := range readConfigMap.DeviceModels {
		if strings.ToUpper(deviceModel.Name) == strings.ToUpper(modelName) {
			for _, property := range deviceModel.Properties {
				name := strings.ToUpper(property.Name)
				if name == strings.ToUpper(CH1pinNumberConfig) {
					if num, ok := property.DefaultValue.(float64); !ok {
						log.Fatalf("get ch1 pin number error %v", property.DefaultValue)
						return errors.New(" Error in reading ch1 pin number from config map")
					} else {
						ch1_wpi_num = int64(num)
					}

				}
				if name == strings.ToUpper(CH2pinNumberConfig) {
					if num, ok := property.DefaultValue.(float64); !ok {
						log.Fatalf("get ch2 pin number error ")
						return errors.New(" Error in reading ch2 pin number from config map")
					} else {
						ch2_wpi_num = int64(num)
					}
				}
				if name == strings.ToUpper(CH3pinNumberConfig) {
					if num, ok := property.DefaultValue.(float64); !ok {
						log.Fatalf("get ch3 pin number error ")
						return errors.New(" Error in reading ch3 pin number from config map")
					} else {
						ch3_wpi_num = int64(num)
					}
				}
			}
		}
	}
	fmt.Printf("Get wpi pin number from configmap: ch1 %d ch2 %d ch3 %d\n",
		ch1_wpi_num, ch2_wpi_num, ch3_wpi_num)

	SetOutput(ch1_wpi_num)
	SetOutput(ch2_wpi_num)
	SetOutput(ch3_wpi_num)
	return nil
}

func InitCLient() MQTT.Client {
	fmt.Println("init client ...")
	fmt.Println("DeviceID", deviceID)
	fmt.Println("mqtturl", mqtturl)

	onceClient.Do(func() {
		opts := MQTT.NewClientOptions().AddBroker(mqtturl).SetClientID(deviceID + "-client").SetCleanSession(true)
		opts = opts.SetKeepAlive(10)
		opts = opts.SetOnConnectHandler(func(c MQTT.Client) {
			topic := DeviceETPrefix + deviceID + TwinETUpdateDetalSuffix
			fmt.Println("Connected trying to subscribe: ", topic)
			if token := c.Subscribe(topic, 0, OperateUpdateDetalSub); token.Wait() && token.Error() != nil {
				fmt.Println("subscribe error: ", token.Error())
				//os.Exit(1)
			}
		})
		opts = opts.SetConnectionLostHandler(func(c MQTT.Client, err error) {
			fmt.Println("connection lost: ", err.Error())
			//os.Exit(1)
		})

		//opts.SetAutoReconnect(true)
		//opts.SetMaxReconnectInterval(10 * time.Second)

		Client = MQTT.NewClient(opts)
	})
	return Client
}

func OperateUpdateDetalSub(c MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Receive msg topic %s %v\n\n", msg.Topic(), string(msg.Payload()))
	current := &dttype.DeviceTwinUpdate{}
	if err := json.Unmarshal(msg.Payload(), current); err != nil {
		fmt.Printf("unmarshl receive msg DeviceTwinUpdate{} to error %v\n", err)
		return
	}
	value := *(current.Twin[CH1_STATE].Expected.Value)
	if RelayState(ch1_wpi_num) != value {
		if err := Set(ch1_wpi_num, value); err != nil {
			fmt.Printf("Set CH1 to %v error %v", value, err)
		}
	}

	value = *(current.Twin[CH2_STATE].Expected.Value)
	if RelayState(ch2_wpi_num) != value {
		if err := Set(ch2_wpi_num, value); err != nil {
			fmt.Printf("Set CH2 to %v error %v", value, err)
		}
	}
	value = *(current.Twin[CH3_STATE].Expected.Value)
	if RelayState(ch3_wpi_num) != value {
		if err := Set(ch3_wpi_num, value); err != nil {
			fmt.Printf("Set CH3 to %v error %v", value, err)
		}
	}
}

func CreateActualDeviceStatus(actch1, actch2, actch3 string) dttype.DeviceTwinUpdate {
	act := dttype.DeviceTwinUpdate{}
	actualMap := map[string]*dttype.MsgTwin{
		CH1_STATE: {
			Actual:   &dttype.TwinValue{Value: &actch1},
			Metadata: &dttype.TypeMetadata{Type: "Updated"}},
		CH2_STATE: {
			Actual:   &dttype.TwinValue{Value: &actch2},
			Metadata: &dttype.TypeMetadata{Type: "Updated"}},
		CH3_STATE: {
			Actual:   &dttype.TwinValue{Value: &actch3},
			Metadata: &dttype.TypeMetadata{Type: "Updated"}},
	}
	act.Twin = actualMap
	return act
}

func RelayState(number int64) string {
	s, err := State(number)
	if err != nil {
		log.Fatalf("get Channel %d State  error %v", number, err)
	}
	switch s[0] {
	case '1':
		return "OFF"
	case '0':
		return "ON"
	}
	return UNKNOW
}

func UpdateActualDeviceStatus() {
	//r .y. g

	deviceTwinUpdate := DeviceETPrefix + deviceID + TwinETUpdateSuffix
	for {
		act := CreateActualDeviceStatus(RelayState(ch1_wpi_num), RelayState(ch2_wpi_num), RelayState(ch3_wpi_num))

		//twinUpdateBody, err := json.MarshalIndent(act, "", "	")
		twinUpdateBody, err := json.Marshal(act)
		if err != nil {
			log.Fatalf("Error:  %v", err)
		}
		token := Client.Publish(deviceTwinUpdate, 1, false, twinUpdateBody)
		if token.Wait() && token.Error() != nil {
			log.Fatalf("client.publish() Error in device twin update is %v", token.Error())
		}

		//fmt.Printf("update deviceTwin %++v\n", string(twinUpdateBody))

		time.Sleep(time.Second * 3)
	}

}

//DeviceStateUpdate is the structure used in updating the device state
type DeviceStateUpdate struct {
	State string `json:"state,omitempty"`
}

/*
func ChangeDeviceState(state string) {
	fmt.Println("Changing the state of the device to online")
	var deviceStateUpdateMessage DeviceStateUpdate
	deviceStateUpdateMessage.State = state
	stateUpdateBody, err := json.Marshal(deviceStateUpdateMessage)
	if err != nil {
		log.Fatalf("Error:   %v", err)
	}
	deviceStatusUpdate := DeviceETPrefix + deviceID + DeviceETStateUpdateSuffix
	token := Client.Publish(deviceStatusUpdate, 0, false, stateUpdateBody)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("client.publish() Error in device state update  is  %v", token.Error())
	}
}
*/

//getTwin function is used to get the device twin details from the edge
/*
func GetTwin(updateMessage dttype.DeviceTwinUpdate) {
	getTwin := DeviceETPrefix + deviceID + TwinETGetSuffix
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	token := Client.Publish(getTwin, 0, false, twinUpdateBody)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("client.publish() Error in device twin get  is ", token.Error())
	}
}
*/
