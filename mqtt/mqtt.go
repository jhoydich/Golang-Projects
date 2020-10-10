package main

import(
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"errors"
	"time"
	"os"
)

qosErr := errors.New("QOS needs to be from 0 to 2")
var flag bool = false

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	if strings.Compare(string(payload), "\n") > 0 {
		fmt.Printf("Topic: %s\n", topic)
		fmt.Printf("MSG: %s\n", payload)
	}
}

func main() {

	//Options for the MQTT connection
	opts := MQTT.NewClientOptions().AddBroker("tcp://mqtt.eclipse.org:1883")
	opts.SetClientID("Device-sub")
	opts.SetDefaultPublishHandler(f)

	//New MQTT connection 
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}


}

//Func to publish 
func Publish(topic string, msg string, qos uint16, c Client) error {
	if qos > 2 {
		return qosError
	}
	token := c.Publish(topic, qos, false, msg) 
	token.Wait()

	time.Sleep(time.Second * 1)
	return nil
}

//Subscribe function
func Subscribe(topic string, msg string, qos uint16, c Client) error {
	if qos > 2 {
		return qosError
	}
	if token := c.Publish(topic, qos, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	

	return nil
}