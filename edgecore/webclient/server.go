package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
	//"path/filepath"
	//"k8s.io/client-go/util/homedir"
	//"flag"
)

var dynamicClient dynamic.Interface

func GetDevices(c echo.Context) error {

	result := getDevices(dynamicClient)

	return c.JSON(http.StatusOK, result)
}

// http://localhost:1323/device/relay-instance-01
func GetDeviceState(c echo.Context) error {
	deviceName := c.Param("data")

	// TODO evaludate model to select
	_, model := getDevice(dynamicClient, deviceName)

	if model == "relay-model" {
		_, ch1Value, ch2Value, ch3Value := getRelayState(dynamicClient, deviceName)
		return c.JSON(http.StatusOK, map[string]string{
			"device": deviceName,
			"ch1":    ch1Value,
			"ch2":    ch2Value,
			"ch3":    ch3Value})
	}

	if model == "dht-model" {
		_, temperature, humidity := getDHTState(dynamicClient, deviceName)
		return c.JSON(http.StatusOK, map[string]string{
			"device":      deviceName,
			"temperature": temperature,
			"humidity":    humidity})
	}

	return c.JSON(http.StatusInternalServerError, map[string]string{
		"device": "eror",
	})

}

func UpdateDeviceState(c echo.Context) error {
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

	// Todo validate how to override correct values
	_, deviceRelay.CH1Value, deviceRelay.CH2Value, deviceRelay.CH3Value = getRelayState(dynamicClient, deviceRelay.Device)

	return c.JSON(http.StatusOK, map[string]string{
		"device": deviceRelay.Device,
		"ch1":    deviceRelay.CH1Value,
		"ch2":    deviceRelay.CH2Value,
		"ch3":    deviceRelay.CH3Value})
}

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, _ = dynamic.NewForConfig(config)

	e := echo.New()

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "valid-key", nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/device", GetDevices)

	e.GET("/device/:data", GetDeviceState)

	e.POST("/device", UpdateDeviceState)

	e.Logger.Fatal(e.Start(":1323"))
}
