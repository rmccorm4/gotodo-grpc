package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/rmccorm4/gotodo-grpc/todo"
	"google.golang.org/grpc"
)

const (
	dbPath       = "db.pb"
	sizeOfLength = 8
)

var endianness = binary.LittleEndian

type length int64

// Implement TaskServiceServer interface
type taskServer struct {
}

func (taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {
	task := &todo.Task{
		Text: text.Text,
		Done: false,
	}

	bs, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("Couldn't encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, endianness, length(len(bs))); err != nil {
		return nil, fmt.Errorf("Couldn't encode length of message: %v", err)
	}
	_, err = f.Write(bs)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", dbPath, err)
	}
	return task, nil
}

func (ts taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error) {
	bs, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read file %s: %v", dbPath, err)
	}

	var tasks todo.TaskList
	for {
		if len(bs) == 0 {
			// Return TaskList for Client to process
			return &tasks, nil
		} else if len(bs) < sizeOfLength {
			return nil, fmt.Errorf("Bytes missing length header, only %d bytes remaining in byte slice.", len(bs))
		}

		// Get first 4 bytes containing the length of our message
		var length int64
		if err := binary.Read(bytes.NewReader(bs[:sizeOfLength]), endianness, &length); err != nil {
			return nil, fmt.Errorf("Couldn't decode message length: %v", err)
		}
		// Remove length bytes from the beginning so we can read current task
		bs = bs[sizeOfLength:]

		var task todo.Task
		if err := proto.Unmarshal(bs[:length], &task); err != nil {
			return nil, fmt.Errorf("Couldn't read task: %v", err)
		}

		// Remove message bytes from the beginning so we can read next length/task
		bs = bs[length:]
		tasks.Tasks = append(tasks.Tasks, &task)
	}
}

func main() {
	var tasks taskServer
	srv := grpc.NewServer()
	todo.RegisterTaskServiceServer(srv, tasks)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Couldn't listen to :8888: %v", err)
	}

	log.Println("Listening on port 8888...")
	log.Fatal(srv.Serve(l))
}
