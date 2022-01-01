package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
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

	serverAddr, err := net.ResolveUDPAddr("udp", ":10001");
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
	for _, value := range messages {
		fmt.Println(value)
	}
}