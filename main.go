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
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get working directory:", err)
	}
	log.Println("Working directory:", wd)

	err = loadDB()
	if err != nil {
		log.Println("Could not load database:", err)
	} else {
		log.Println("Database loaded successfully")
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
				fmt.Fprintf(conn, "-ERR wrong number of arguments for '%s' command\r\n", strings.ToLower(args[0]))
				continue
			}
			key, value := args[1], args[2]
			DB[key] = value
			if err := persistDB(); err != nil {
				log.Println("Failed to persist DB:", err)
			}
			fmt.Fprint(conn, "+OK\r\n")
		case "GET":
			if len(args) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for '%s' command\r\n", strings.ToLower(args[0]))
				continue
			}
			key := args[1]
			val, ok := DB[key]
			if !ok {
				fmt.Fprint(conn, "$-1\r\n")
			} else {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(val), val)
			}
		case "DELETE":
			if len(args) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for '%s' command\r\n", strings.ToLower(args[0]))
				continue
			}
			key := args[1]
			if _, ok := DB[key]; ok {
				delete(DB, key)
				if err := persistDB(); err != nil {
					log.Println("Failed to persist DB:", err)
				}
				fmt.Fprint(conn, ":1\r\n")
			} else {
				fmt.Fprint(conn, ":0\r\n")
			}
		default:
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", strings.ToLower(args[0]))
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
