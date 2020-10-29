	package main

//gets new data 


import (
    MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	_ "errors"
	"time"
	"os"
	"database/sql"
	"strings"
	"strconv"
	_ "github.com/lib/pq"

)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = ""
	dbname = "sampleDB"
)

var db *sql.DB


var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	s := string(payload)
	p, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println(topic)
	go dbInsert(db, topic, p)		
}

func main() {
	db = dbConnect()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	mqttConnect()
	
	for {
		time.Sleep(time.Millisecond * 500)
	}
}

func dbConnect () *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}


func dbInsert(db *sql.DB, topic string, payload float64) {
	l := strings.Split(topic, "/")
	
	t := l[1]


	stmt := `INSERT INTO plantdata (dtype, payload) VALUES ($1, $2)`
	_, err := db.Exec(stmt, t, payload)
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

	if token := c.Subscribe("jsh/#", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
