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
	"github.com/google/uuid"
	_ "github.com/lib/pq"

)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "P3rmaS0rt"
	dbname = "sampleDB"
)

var db *sql.DB


var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := float64(msg.Payload())
	
	fmt.Println(topic)
	fmt.Println(string(payload))	
}

func main() {
	db = dbConnect()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	mqttConnect()	
}

func dbConnect () *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

//Plant check will see if the plant in the topic is in the db table listing all plants
func plantCheck(db *sql.DB, topic string, payload float64) {
	var id uuid.UUID
	l := strings.Split(topic, "/")
	user := l[0]
	plant := l[1]
	t := l[2]

	nameqry := `SELECT dev-id FROM plantlist WHERE plantName=$1 and username=$2;`
	nameinsert := `INSERT INTO plantlist (plantname, username, dev-id) VALUES ($1, $2, $3);`

	
	fmt.Println(time.Now())
	row := db.QueryRow(nameqry, plant, user)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		id = uuid.New()
		_, err = db.Exec(nameinsert, plant, user, id)
		dbInsert(db, t, payload, id)
	case nil:
		dbInsert(db, t, payload, dev-id)
		fmt.Println(time.Now())
	default:
		panic(err)
	 }



}

func dbInsert(db *sql.DB, topic string, payload float64) {
	stmt := `INSERT INTO plantdata (dtype, payload, dev-id) VALUES ($1, $2, $3)`
	_, err := db.Exec(stmt, topic, payload)
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

	if token := c.Subscribe("jhoy/plant/temp", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

