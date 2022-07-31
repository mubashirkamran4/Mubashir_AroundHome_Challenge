package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	hostname = "127.0.0.1:3306"
	dbname   = "around_home"
	username = "root"
	password = "root"
)

type partner struct {
	name string
	experienced_in []string
	latitude float32
	longitude float32
	operating_radius_latitude float32
	operating_radius_longitude float32
	rating int
}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	//defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

func createPartnerTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Partner(partner_id int primary key auto_increment, name text, 
        experienced_in JSON, latitude float, longitude float, operating_radius_latitude float,
        operating_radius_longitude float, rating int)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating partner table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func multipleInsert(db *sql.DB, partners []partner) error {
	query := "INSERT INTO partner(name, experienced_in, latitude, longitude, operating_radius_latitude, " +
				"operating_radius_longitude, rating) VALUES "
	var inserts []string
	var params []interface{}
	for _, v := range partners {
		inserts = append(inserts, "(?, ?)")
		params = append(params, v.name, v.experienced_in, v.latitude, v.longitude, v.operating_radius_latitude,
			v.operating_radius_longitude, v.rating)
	}
	queryVals := strings.Join(inserts, ",")
	query = query + queryVals
	log.Println("query is", query)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		log.Printf("Error %s when inserting row into partners table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d partners created simultaneously", rows)
	return nil
}

func main() {
	db, err := dbConnection()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()
	log.Printf("Successfully connected to database")
	err = createPartnerTable(db)
	if err != nil {
		log.Printf("Create partner table failed with error %s", err)
		return
	}

	p1 := partner{
		name:  "partner1",
		latitude: 98.9,
		longitude: 10.2,
		operating_radius_latitude: 10.5,
		operating_radius_longitude: 4.5,
		rating: 1,
		experienced_in: []string{"wood"},
	}
	p2 := partner{
		name:  "partner2",
		latitude: 76.9,
		longitude: 1.2,
		operating_radius_latitude: 7.5,
		operating_radius_longitude: 2.5,
		rating: 2,
		experienced_in: []string{"tiles"},
	}
	p3 := partner{
		name:  "partner3",
		latitude: 67.9,
		longitude: 1.6,
		operating_radius_latitude: 105.5,
		operating_radius_longitude: 28.5,
		rating: 3,
		experienced_in: []string{"flooring"},

	}
	p4 := partner{
		name:  "partner4",
		latitude: 68.9,
		longitude: 2.2,
		operating_radius_latitude: 125.5,
		operating_radius_longitude: 55.5,
		rating: 4,
		experienced_in: []string{"wood", "tiles"},
	}
	p5 := partner{
		name:  "partner5",
		latitude: 198.9,
		longitude: 28.5,
		operating_radius_latitude: 7.5,
		operating_radius_longitude: 2.5,
		rating: 5,
		experienced_in: []string{"wood", "tiles", "flooring"},
	}


	err = multipleInsert(db, []partner{p1, p2, p3, p4, p5})
	if err != nil {
		log.Printf("Multiple insert failed with error %s", err)
		return
	}
}