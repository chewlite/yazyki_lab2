package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"protector"
	"strconv"
	"strings"
)

func main() {

	args := os.Args[1:]

	switch args[0] {
	case "server_mode":
		port := args[1]
		limit := args[2]
		startServer(port, limit)
	case "client_mode":
		host := args[1]
		port := args[2]
		startClient(host, port)
	}

}

func startServer(port string, limit string) {

	fmt.Println("Launching server...")

	connectionsLimit, _ := strconv.Atoi(limit)
	ln, _ := net.Listen("tcp", ":"+port)

	connected := 0

	for {
		if connected >= connectionsLimit {
			fmt.Println("connections limit reached")
			os.Exit(0)
		} else {
			conn, _ := ln.Accept()
			connected++
			go workWithConnection(conn, &connected)
		}
	}

}

func startClient(host string, port string) {
	con, _ := net.Dial("tcp", host+":"+port)
	hashStr := protector.GetHashStr()
	firstKey := protector.GetSessionKey()
	log.Println("hash: " + hashStr + ", first key: " + firstKey)
	writeMessage(con, hashStr+" "+firstKey)
	keyFromServer := readMessage(con)
	nextKey := protector.NextSessionKey(hashStr, firstKey)
	fmt.Println("generated key: " + nextKey)
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Println("key from server: " + keyFromServer + ", key in client: " + nextKey)
		if keyFromServer != nextKey {
			break
		}
		nextKey = protector.NextSessionKey(hashStr, nextKey)
		fmt.Println("generated key: " + nextKey)
		fmt.Print("write message: ")
		userMessage, _ := reader.ReadString('\n')
		userMessage = userMessage[:len(userMessage)-1]
		writeMessage(con, nextKey+" "+userMessage)
		messageFromServer := readMessage(con)
		messageParts := strings.Split(messageFromServer, " ")
		keyFromServer = messageParts[0]
		nextKey = protector.NextSessionKey(hashStr, nextKey)
		fmt.Println("generated key: " + nextKey)
	}
}

func workWithConnection(con net.Conn, connected *int) {
	defer closeConnection(con, connected)
	message := readMessage(con)
	messageParts := strings.Split(message, " ")
	startString := messageParts[0]
	firstKey := messageParts[1]
	fmt.Println("first key: " + firstKey)
	nextKey := protector.NextSessionKey(startString, firstKey)
	writeMessage(con, nextKey)

	var key string
	var innerMessage string
	var newKey string
	for {
		message = readMessage(con)
		messageParts = strings.Split(message, " ")
		key = messageParts[0]
		innerMessage = messageParts[1]
		newKey = protector.NextSessionKey(startString, key)
		fmt.Println("generated key: " + newKey)
		writeMessage(con, newKey+" "+innerMessage)
	}
}

func closeConnection(con net.Conn, connected *int) {
	con.Close()
	*connected--
}

func readMessage(con net.Conn) string {
	message, _ := bufio.NewReader(con).ReadString('\n')
	message = message[:len(message)-1]
	fmt.Println("incoming message: " + message)
	return message
}

func writeMessage(con net.Conn, message string) {
	fmt.Println("outgoing message: " + message)
	con.Write([]byte(message + "\n"))
}
