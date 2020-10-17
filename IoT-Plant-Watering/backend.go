package main

import (
    MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"errors"
	"time"
	"os"
	"database/sql"
	_ "github.com/lib/pq"

)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = ""
	dbname = "sampledb"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	t := time.Now()
	min := t.Minute()
	fmt.Println(topic)
	fmt.Println(string(payload))	
}

func main() {
	c := make(chan)


	db := dbConnect()
	
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	mqttConnect()


}

func dbConnect () *DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func dbInsert(db *DB) {
	stmt := `INSERT INTO plant (age, email, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err = db.Exec(stmt, 24, "jhoydich3@gmail.com", "Jeremiah", "Hoydich")
	if err != nil {
		panic(err)
	}

}

func mqttConnect() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	opts.SetClientID("Device-sub")
	opts.SetDefaultPublishHandler(f)

	//New MQTT connection 
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.SubscribeMultiple(["esp/test"],0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}

func Collector() {

}