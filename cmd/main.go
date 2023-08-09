package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/raj3k/BlazeDB/blazedb"
	"github.com/raj3k/BlazeDB/internal/proto"
	"github.com/raj3k/BlazeDB/internal/utils"
	"go.etcd.io/bbolt"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	BUCKET    = "default"
)

type RedisLiteServer struct {
	data  map[string]string
	mutex sync.Mutex
}

func NewRedisLiteServer() *RedisLiteServer {
	return &RedisLiteServer{
		data: make(map[string]string),
	}
}

func (rls *RedisLiteServer) processCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}

	rls.mutex.Lock()
	defer rls.mutex.Unlock()

	switch parts[0] {
	case "get":
		if len(parts) != 2 {
			return "ERROR: Invalid argument for GET command"
		}
		val, found := rls.data[parts[1]]
		if !found {
			return "(nil)"
		}
		return "$3\r\n" + val
	case "set":
		if len(parts) != 3 {
			return "ERROR: Invalid argument for SET command"
		}
		rls.data[parts[1]] = parts[2]
		return "+OK"
	case "exists":
		if len(parts) != 2 {
			return "ERROR: Invalid argument for GET command"
		}
		_, found := rls.data[parts[1]]

		if !found {
			return "0"
		}
		return "1"
	case "ping":
		return "+pong"
	default:
		return "ERROR: Unknown command"
	}
}

func (rls *RedisLiteServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 4096)
		readLen, err := conn.Read(buffer)

		if err != nil {
			log.Println("Error reading from connection: ", err)
			return
		}

		r := proto.NewReader(strings.NewReader(string(buffer[:readLen])))

		cmd, err := r.ReadReply()

		if err != nil {
			log.Println("Error reading reply: ", err)
		}

		response := ""

		if utils.IsSliceOfInterface(cmd) {
			iSlice := cmd.([]interface{})
			stringSlice := make([]string, len(iSlice))

			for i, v := range iSlice {
				// Convert each element to a string representation
				stringSlice[i] = fmt.Sprintf("%v", v)
			}

			resultString := strings.Join(stringSlice, " ")

			response = rls.processCommand(resultString)
		}

		if utils.IsInterface(cmd) {
			response = rls.processCommand(cmd.(string))
		}

		_, err = conn.Write([]byte(response + "\r\n"))
		if err != nil {
			log.Println("Error writing to connection: ", err)
			return
		}
	}
}

func main() {
	// rls := NewRedisLiteServer()

	// listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	// if err != nil {
	// 	log.Fatal("Error creating listener: ", err)
	// }
	// defer listener.Close()

	// log.Println("Redis-lite server is now listening on port ", CONN_PORT)

	// for {
	// 	conn, err := listener.Accept()

	// 	if err != nil {
	// 		log.Println("Error accepting connection: ", err)
	// 		continue
	// 	}
	// 	go rls.handleConnection(conn)
	// }

	db, err := blazedb.New()
	if err != nil {
		log.Fatal(err)
	}

	k := []byte("1")
	v := []byte("Hello World!")

	db.Set(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET))
		if err != nil {
			return err
		}
		return b.Put(k, v)
	})

	db.Get(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET))
		v1 := b.Get([]byte(k))
		fmt.Println(string(v1))
		return nil
	})

	// fmt.Println(result)

	db.DropDatabase("blaze")
}
