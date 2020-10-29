//Plant check will see if the plant in the topic is in the db table listing all plants
package main


func main() {

}


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
