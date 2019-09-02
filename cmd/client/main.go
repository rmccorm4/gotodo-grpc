package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rmccorm4/gotodo-grpc/todo"
	"google.golang.org/grpc"
)

func add(ctx context.Context, client todo.TaskServiceClient, text string) error {
	_, err := client.Add(ctx, &todo.Text{Text: text})
	if err != nil {
		return fmt.Errorf("Couldn't add task to TaskList: %v", err)
	}

	fmt.Printf("Task added: [%s]\n", text)
	return nil
}

func list(ctx context.Context, client todo.TaskServiceClient) error {
	list, err := client.List(ctx, &todo.Void{})
	if err != nil {
		return fmt.Errorf("Couldn't fetch task list: %v", err)
	}

	for _, task := range list.Tasks {
		if task.Done {
			fmt.Printf("✔️")
		} else {
			fmt.Printf("❌")
		}
		fmt.Printf(" %s\n", task.Text)
	}

	return nil
}

func main() {
	// Create gRPC Client and connect to gRPC Server
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err)
	}
	client := todo.NewTaskServiceClient(conn)

	// Parse CLI args
	flag.Parse()
	if flag.NArg() < 1 {
		subcommands := []string{"add", "list"}
		fmt.Println("Missing subcommand:", subcommands)
		os.Exit(1)
	}

	switch cmd := flag.Arg(0); cmd {
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	case "list":
		err = list(context.Background(), client)
	default:
		err = fmt.Errorf("Unknown subcommand: %s", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
