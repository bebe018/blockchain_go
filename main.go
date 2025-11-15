package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 使用 := 宣告
	addr := ":2112"
	os.Setenv("NODE_ID", "3000")

	metricsManager := NewMetricsManager(addr)

	metricsErrCh := make(chan error, 1)

	osSignalCh := make(chan os.Signal, 1)
	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("Starting metrics server on %s ...\n", addr)
	go metricsManager.RunMetricsServer(metricsErrCh)

	cli := CLI{}
	eventCh := cli.StartListener()

	for {
		select {
		case event := <-eventCh:
			fmt.Printf("main function receives event: command [%s], success: %t, message: %s\n", event.Command, event.Success, event.Message)
			if event.Command == "quit" || event.Command == "exit" {
				goto Shutdown
			}

		case err := <-metricsErrCh:
			fmt.Printf("Metrics processing error: %v\n", err)
			goto Shutdown

		case sig := <-osSignalCh:
			fmt.Println("Receive signal inter\n", sig)
			goto Shutdown
		}
	}

Shutdown:
	fmt.Println("--- Graceful shutdown executing ---")

	if err := metricsManager.GracefulShutdownMetricsServer(context.Background()); err != nil {
		fmt.Printf("Shutdown Metrics Server fail: %v\n", err)
	}

	fmt.Println("function has been terminated")
}
