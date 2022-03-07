package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var dynamicClient dynamic.Interface

func GetRelays(c echo.Context) error {

	result := getRelays(dynamicClient)

	return c.JSON(http.StatusOK, result)
}

// http://localhost:1323/device/relay-instance-01
func GetRelayState(c echo.Context) error {
	deviceName := c.Param("data")

	_, ch1Value, ch2Value, ch3Value := getRelayState(dynamicClient, deviceName)

	return c.JSON(http.StatusOK, map[string]string{
		"device": deviceName,
		"ch1":    ch1Value,
		"ch2":    ch2Value,
		"ch3":    ch3Value})
}

func PathRelayState(c echo.Context) error {
	type DeviceRelay struct {
		Device   string `json:"device"`
		CH1Value string `json:"ch1"`
		CH2Value string `json:"ch2"`
		CH3Value string `json:"ch3"`
	}
	deviceRelay := DeviceRelay{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&deviceRelay)

	if err != nil {
		fmt.Printf("Failed reading the request body %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	fmt.Printf("\nDevice: %s\n[CH1]: %s\n[CH2]: %s\n[CH3]: %s\n", deviceRelay.Device, deviceRelay.CH1Value, deviceRelay.CH2Value, deviceRelay.CH3Value)

	err = setRelayState(dynamicClient, deviceRelay.Device, deviceRelay.CH1Value, deviceRelay.CH2Value, deviceRelay.CH3Value)

	if err != nil {
		fmt.Printf("Failed executing set relay state %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"device": deviceRelay.Device,
		"ch1":    deviceRelay.CH1Value,
		"ch2":    deviceRelay.CH2Value,
		"ch3":    deviceRelay.CH3Value})
}

func main() {
	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()

	// use the current context in kubeconfig
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, _ = dynamic.NewForConfig(config)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/device", GetRelays)

	e.GET("/device/:data", GetRelayState)

	e.POST("/device", PathRelayState)

	e.Logger.Fatal(e.Start(":1323"))
}
