package rest

import (
	"net/http"
	"strconv"

	"adapter/log"

	"github.com/gin-gonic/gin"
	"adapter/config"
	"adapter/device"
)

func StartHttpServer(errChan chan error, port int) {
	go func() {
		r := LoadRestRoutes()
		errChan <- r.Run(":" + strconv.Itoa(port))
	}()
}

func LoadRestRoutes() *gin.Engine {
	r := gin.New()
	r.GET("/api/v1/ping", pingHandler)
	r.POST("/api/v1/device", restAddNewDevice)
	r.GET("/api/v1/device", restGetAllDevices)
	r.GET("/api/v1/device/delete", restRemoveDevice)
	r.POST("/api/v1/device/update", restUpdateDevice)
	r.GET("/api/v1/config", restGetConfig)
	r.PUT("/api/v1/config", restSetConfig)
	return r
}

func pingHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	c.String(http.StatusOK, "pong")
}

func restAddNewDevice(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	var d device.DeviceInfo
	err := c.BindJSON(&d)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logclient.Log.Println(err)
		return
	}

	err = device.CreateDevice(&d, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logclient.Log.Println(err)
		return
	}
	c.String(http.StatusOK, "OK")
}

func restGetAllDevices(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	dl := device.GetAllDevices()
	c.JSON(http.StatusOK, dl)
}

func restRemoveDevice(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	t := c.Query("type")
	n := c.Query("name")
	err := device.RemoveDevice(t, n)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logclient.Log.Println(err)
		return
	}
	c.String(http.StatusOK, "OK")
}

func restUpdateDevice(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	n := c.Query("name")
	t := c.Query("type")
	var d device.DeviceInfo
	err := c.BindJSON(&d)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logclient.Log.Println(err)
		return
	}
	err = device.UpdateDevice(t, n, &d)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logclient.Log.Println(err)
		return
	}
	c.String(http.StatusOK, "OK")
}

func restGetConfig(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	cfg := config.GetGloablConfig()
	c.JSON(http.StatusOK, cfg)
}

func restSetConfig(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	var cfg config.AdapterConfig
	err := c.BindJSON(&cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logclient.Log.Println(err)
		return
	}
	config.SetNewConfig(cfg)
	c.String(http.StatusOK, "OK")
}
