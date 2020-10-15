package main

import(
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	_ "errors"
	"time"
	"os"
)


var flag bool = false

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	//fmt.Println(string(payload))
	if topic == "esp/test" && string(payload) == "request" {
		if flag == true {
			token := client.Publish("esp/led", 0, false, "1")
			flag = false
			token.Wait()
		} else {
			token := client.Publish("esp/led", 0, false, "0")
			flag = true
			token.Wait()
		}
		//fmt.Println(flag)
	}
	
}

func main() {
	//qosErr := errors.New("QOS needs to be from 0 to 2")
	//Options for the MQTT connection
	opts := MQTT.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	opts.SetClientID("Device-sub")
	opts.SetDefaultPublishHandler(f)

	//New MQTT connection 
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe("esp/test",0,nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	

	for {
		time.Sleep(time.Second * 1)
	}
}

/*Func to publish 
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
*/