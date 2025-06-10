package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var DB = make(map[string]string)

const dbFile = "redis.db"

func main() {
	err := loadDB()
	if err != nil {
		log.Println("Could not load database:", err)
	}

	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on port :6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		if len(args) == 0 {
			continue
		}
		cmd := strings.ToUpper(args[0])

		switch cmd {
		case "SET":
			if len(args) != 3 {
				fmt.Fprintln(conn, "Error: Wrong number of arguments for 'SET'")
				continue
			}
			key, value := args[1], args[2]
			DB[key] = value
			if err := persistDB(); err != nil {
				log.Println("Failed to persist DB:", err)
			}
			fmt.Fprintln(conn, "OK")
		case "GET":
			if len(args) != 2 {
				fmt.Fprintln(conn, "Error: Wrong number of arguments for 'GET'")
				continue
			}
			key := args[1]
			val, ok := DB[key]
			if !ok {
				fmt.Fprintln(conn, "(nil)")
			} else {
				fmt.Fprintln(conn, val)
			}
		case "DELETE":
			if len(args) != 2 {
				fmt.Fprintln(conn, "Error: Wrong number of arguments for 'DELETE'")
				continue
			}
			key := args[1]
			delete(DB, key)
			if err := persistDB(); err != nil {
				log.Println("Failed to persist DB:", err)
			}
			fmt.Fprintln(conn, "OK")
		default:
			fmt.Fprintln(conn, "Error: Unknown command")
		}
	}
}

func persistDB() error {
	file, err := os.Create(dbFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(DB)
}

func loadDB() error {
	file, err := os.Open(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, so we start with an empty DB.
		}
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&DB)
}
