package main

import (
	"encoding/hex"
	"fmt"
	"os"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"

	"github.com/go-redis/redis"
	"strconv"

	"encoding/json"
	"encoding/base64"
	"github.com/TheThingsNetwork/ttn/core/types"
)

const (
	sdkClientName = "my-amazing-app"
)

type NetworkMetaData struct {
	Time types.JSONTime `json:"time"`
}

type MessageMetadata struct {
	Network NetworkMetaData `json:"network"`

}

type Message struct {
	Device string `json:"device"`
	Payload string `json:"payload"`
	Metada MessageMetadata `json:"metadata"`
}



func main() {
	log := apex.Stdout() // We use a cli logger at Stdout
	ttnlog.Set(log)      // Set the logger as default for TTN

	// We get the application ID and application access key from the environment
	appID := os.Getenv("TTN_APP_ID")
	appAccessKey := os.Getenv("TTN_APP_ACCESS_KEY")

	// Create a new SDK configuration for the public community network
	config := ttnsdk.NewCommunityConfig(sdkClientName)
	config.ClientVersion = "2.0.5" // The version of the application


	// Create a new SDK client for the application
	client := config.NewClient(appID, appAccessKey)

	// Make sure the client is closed before the function returns
	// In your application, you should call this before the application shuts down
	defer client.Close()

	// Manage devices for the application.
	devices, err := client.ManageDevices()
	if err != nil {
		log.WithError(err).Fatal("my-amazing-app: could not get device manager")
	}

	// List the first 10 devices
	deviceList, err := devices.List(10, 0)
	if err != nil {
		log.WithError(err).Fatal("my-amazing-app: could not get devices")
	}
	log.Info("my-amazing-app: found devices")
	for _, device := range deviceList {
		fmt.Printf("- %s", device.DevID)
	}



	// Start Publish/Subscribe client (MQTT)
	pubsub, err := client.PubSub()
	if err != nil {
		log.WithError(err).Fatal("my-amazing-app: could not get application pub/sub")
	}

	// Make sure the pubsub client is closed before the function returns
	// In your application, you should call this before the application shuts down
	defer pubsub.Close()



	all_devices := pubsub.AllDevices()
	defer all_devices.Close()
	uplink, err := all_devices.SubscribeUplink()
	if( err != nil ){
		fmt.Println("Error");
	}

	db_client := redis.NewClient(&redis.Options{
		Addr:	"localhost:6379",
		Password:	"",
		DB:	0,
	})


	for message := range uplink {
		hexPayload := hex.EncodeToString(message.PayloadRaw)
		log.WithField("data", hexPayload).Info("my-amazing-app: received uplink")

		n := NetworkMetaData{message.Metadata.Time}
		meta := MessageMetadata{n}
		m := Message{
			message.HardwareSerial,
			base64.StdEncoding.EncodeToString(message.PayloadRaw),
			meta,
		}

		data, err := json.Marshal(m)
		if(err != nil){
			panic("We have a json probel")
		}

		if err := db_client.LPush("lora:rx:" + strconv.Itoa(int(message.FPort)), data).Err(); err != nil {
			panic(err)
		}


	}

	err = all_devices.UnsubscribeUplink()
	if err != nil {
		log.WithError(err).Fatal("my-amazing-app: could not unsubscribe from uplink")
	}

}
