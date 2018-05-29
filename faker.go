package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/manveru/faker"

	_ "github.com/lib/pq"
)

type xuserfaker struct {
	UserName   string
	Email      string
	Password   string
	Name       string
	Phone      string
	Gender     int
	Bio        string
	Credit     int
	Language   int
	Country    int
	TimeZone   int
	LastIP     string
	Updatetime int
	Createtime int
}

type conn struct {
	db   *sql.DB
	name string
}

func genXuserFaker(fake *faker.Faker, r *rand.Rand) xuserfaker {
	// println(fake.Name())  //> "Adriana Crona"
	// println(fake.Email()) //> charity.brown@fritschbotsford.biz
	u := xuserfaker{
		UserName:   fake.UserName(),
		Email:      fake.Email(),
		Password:   fake.Name(),
		Name:       fake.Name(),
		Phone:      fake.PhoneNumber(),
		Gender:     r.Intn(3),
		Bio:        fake.Sentence(r.Intn(10), true),
		Credit:     0,
		Language:   13,
		Country:    207,
		TimeZone:   28800,
		LastIP:     fake.IPv4Address().String(),
		Updatetime: 1527496777,
		Createtime: 1527496777,
	}
	fmt.Println(u)
	return u
}
func connectDB(dbinfo string, name string) conn {
	var err error
	c := conn{name: name}
	c.db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}
func insertDB(c *conn, u xuserfaker) error {
	// var id
	result, err := c.db.Exec(`INSERT INTO xuser 
		(username, email, salted_password, name, phone, gender, bio, credit, language, country, timezone, last_ip, createtime, updatetime) 
		VALUES 
		('robby', 'robby@gmail.com', 'salted', 'robby', '+886-911111111', 1, 'hi', 0, 13, 207, 28800, '123.194.188.0', 1527496777, 1527496777)`)
	log.Println("xuser result", result, err)
	// rows, err := c.db.Query("SELECT name FROM xuser;")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	name := ""
	// 	if err := rows.Scan(&name); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(name)
	// }
	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	return nil
}
func queryXuser() error {
	// var db *sql.DB
	// if err := connectDB(db); err != nil {
	// 	log.Fatalln("db connection", err)
	// }
	// var comments []string
	// err := db.QueryRow(`SELECT comments from users WHERE id=$1`, id).Scan(pq.Array(&comments))
	// if err != nil {
	// 	return err
	// }
	// log.Println(id, comments)
	return nil
}
func startInsertXuserfaker() error {
	fake, err := faker.New("en")
	if err != nil {
		log.Fatalln("faker", err)
	}
	r := rand.New(rand.NewSource(99))
	c := connectDB("host=localhost port=31497 user=postgres password=mysecretpassword sslmode=disable", "PgSQL")
	defer c.db.Close()
	log.Println(c.db)
	for i := 0; i < 1; i++ {
		u := genXuserFaker(fake, r)
		if err := insertDB(&c, u); err != nil {
			log.Println("insert", err)
		}
	}
	return nil
}
