package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"net"
	"os"
	"server/config"
	"strings"
	"time"
)

type Message struct {
	Name string
	Group string
	Data string
}

var db *sql.DB

func init() {
    if err := godotenv.Load(); err != nil {
        fmt.Println("No .env file found")
    }
}

func main() {
	conf := config.New()

	var connectErr error
	db, connectErr = sql.Open("mysql", conf.MysqlDsn)
	if connectErr != nil {
		panic(connectErr)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()


	messageChn := make(chan string)
	go func() {
		var messages []string
		for {
			select {
				case tmp := <- messageChn:
					messages = append(messages, tmp)
					if len(messages) >= 5 {
						chunk := messages[0:5]
						messages = []string{}
						flush(chunk)
					}
				case <-time.After(10 * time.Second):
					chunk := messages[:]
					messages = []string{}
					flush(chunk)
			}

		}
	}()

	serverAddr, err := net.ResolveUDPAddr("udp", conf.AppPort);
	if err != nil {
		fmt.Println(err);
		os.Exit(0);
	}

	serverConn, err := net.ListenUDP("udp", serverAddr);
	if err != nil {
		fmt.Println(err);
		os.Exit(0);
	}

	defer serverConn.Close();

	buf := make([]byte, 1024);

	var message string;

	for {
		n, _, err := serverConn.ReadFromUDP(buf);
		message = string(buf[0:n]);
		messageChn <- message;
		if err != nil {
			fmt.Println(err);
		}
	}
}

func flush(messages []string) {
	if len(messages) == 0 {
		return
	}

	valueStrings := make([]string, 0, len(messages))
	valueArgs := make([]interface{}, 0, len(messages) * 3)

	var message Message
	for _, value := range messages {
		err := json.Unmarshal([]byte(value), &message)
		if err != nil {
			fmt.Println("Parsing error", err)
			continue
		}

		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, message.Name)
		valueArgs = append(valueArgs, message.Group)
		valueArgs = append(valueArgs, message.Data)
		fmt.Println(message)
	}

	ctx := context.Background()
	SQL := fmt.Sprintf(
		"INSERT INTO logs (record_name, group_name, data) VALUES %s",
		strings.Join(valueStrings, ","),
	)
	_, err := db.ExecContext(ctx, SQL, valueArgs...)
	if err != nil {
		fmt.Println(err)
	}
}