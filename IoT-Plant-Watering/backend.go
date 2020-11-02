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
	//time of message reception, only care about minute
	t := time.Now().Format("02 Jan 06 15:04 MST")
	
	//splitting the topic, getting the subtopic and device id
	l := strings.Split(topic, "/")
	id := l[0]
	top := l[1]
	
	

	//sql statement to insert into database
	stmt := `INSERT INTO plantdata (dtype, payload, dev-id, time) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(stmt, top, payload, id, t)
	if err != nil {
		panic(err)
	}

}


//function connecting to the broker and subscribing to the necessary topics
func mqttConnect() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://broker.hivemq.com:1883")
	opts.SetClientID("Device-sub")
	opts.SetDefaultPublishHandler(f)

	//New MQTT connection 
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribing to device
	if token := c.Subscribe("jsh/#", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
