package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Hamaiz/go-rest-eg/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

// DBConn - connects to postgresql database
func DBConn() (*pgxpool.Pool, error) {
	// connecting to databbase
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	// error handling
	if err != nil {
		return nil, err
	}

	log.Println("connecting to database...")

	return conn, err
}

// DeleteAccount - delete accounts with expiry
func DeleteAccount() {
	log.Println("connected to server")

	// loading dotenv
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found: %v", err)
	}
	switch os.Getenv("APP_ENV") {
	case "development":
		godotenv.Load(".env.development")
	case "production":
		godotenv.Load(".env.production")
	}

	// making connection
	conn, err := DBConn()
	if err != nil {
		log.Println("an error occured: ", err)
	}

	// time tick
	c := time.Tick(30 * time.Minute)
	for _ = range c {
		// get email token struct
		ets := make([]model.EmailDbToken, 0)

		rows, err := conn.Query(context.Background(), "SELECT * FROM addition WHERE confirmed=false")
		switch {
		case err == pgx.ErrNoRows:
			log.Println("there are no unverified users")
		case err != nil:
			log.Println("an error occured")
		}

		defer rows.Close()

		for rows.Next() {
			et := model.EmailDbToken{}
			err := rows.Scan(&et.Confirmed, &et.Expires, &et.Token, &et.Account_id)

			if err != nil {
				log.Println("an error occured")
			}

			ets = append(ets, et)
		}

		for _, e := range ets {
			ex := e.Expires.String()
			exp, _ := time.Parse("2006-01-02 15:04:05", ex)
			exps := time.Now().Local().After(exp)

			id := e.Account_id

			if exps {
				_, err = conn.Exec(context.Background(), "DELETE FROM addition WHERE account_id=$1", id)

				if err != nil {
					log.Println("error occured")
				}

				_, err = conn.Exec(context.Background(), "DELETE FROM account WHERE id=$1", id)

				if err != nil {
					log.Println("error occured")
				}

			}
		}
	}

}
