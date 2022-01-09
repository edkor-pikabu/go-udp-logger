package main

import (
	"fmt"
	"github.com/enriquebris/goconcurrentqueue"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"os"
	"server/handlers"
	"server/helpers/config"
	"server/helpers/db"
	"time"
)

func main() {
	conf := config.New()
	dbConn := db.New(conf)
	handler := handlers.New(dbConn)
	queue := goconcurrentqueue.NewFIFO()

	for i:=1; i<10; i++ {
		go func(){
			var messages []string
			ticker := time.NewTicker(30 * time.Second)
			for {
				select {
				case <-ticker.C:
					chunk := messages[:]
					messages = []string{}
					handler.Handle(chunk)
				default:
					tmp, err := queue.Dequeue()
					if err == nil {
						messages = append(messages, fmt.Sprintf("%v", tmp))
					}
					if len(messages) >= 100 {
						chunk := messages[0:100]
						messages = []string{}
						handler.Handle(chunk)
					}
				}

			}
		}()
	}

	serverConn := runServer(conf)

	defer dbConn.Close()
	defer serverConn.Close()

	buf := make([]byte, 1024)
	var message string

	for {
		n, _, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		message = string(buf[0:n])
		fmt.Println(message)
		queue.Enqueue(message)
	}
}

func runServer(config *config.Config) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", config.AppPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println("Start listening " + config.AppPort)
	return serverConn
}