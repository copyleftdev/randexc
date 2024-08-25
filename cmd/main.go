package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"randexc"
)

func main() {
	maxDuration := flag.String("max-duration", "1m", "Maximum duration for random execution (e.g., 1s, 5m, 2h)")
	actionMessage := flag.String("message", "Action executed!", "Message to print when action is executed")
	flag.Parse()

	executor, err := randexc.New(*maxDuration)
	if err != nil {
		log.Fatalf("Failed to create executor: %v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	fmt.Printf("Waiting up to %s to execute the action...\n", *maxDuration)
	err = executor.Execute(ctx, func() error {
		fmt.Println(*actionMessage)
		return nil
	})

	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	asyncResultChan := executor.ExecuteAsync(ctx, func() error {
		time.Sleep(time.Second)
		fmt.Println("Async action executed!")
		return nil
	})

	asyncResult := <-asyncResultChan
	if asyncResult.Error != nil {
		log.Printf("Async execution failed: %v", asyncResult.Error)
	} else {
		fmt.Printf("Async execution completed. Start: %v, End: %v\n", asyncResult.StartTime, asyncResult.EndTime)
	}
}
