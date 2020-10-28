package main

import(
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	_ "errors"
	"time"
	"os"
	"sync"
)

type plant struct {
	t time.Time
	id string
	topic string
	payload string
}

var p plant

var wg sync.WaitGroup

var ch chan string = make(chan string, 20)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	wg.Add(1)
	topic := string(msg.Topic())
	payload := string(msg.Payload())
	go dataSend(ch, topic, payload)
}

func main() {
	go chanDisplay(ch)
	go chanDisplay(ch)
	go dump(ch)
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
	go chanDisplay(ch)
	
	for {
		time.Sleep(time.Millisecond * 200)
	}
}

func chanDisplay(ch <- chan string) {
	
	x := <-ch
	if x != ""{
		fmt.Println(x)
	}
	
}

func dump(ch chan string) {
	 
	 for {
		 time.Sleep(time.Millisecond * 500)
		 ch<- ""
	 }
}

func dataSend(ch chan<- string, topic string, payload string) {
	defer wg.Done()
	ch <- payload
	ch <- topic
}