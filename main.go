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
			fmt.Printf("📢 主程式監聽到事件: 命令 [%s], 成功: %t, 訊息: %s\n", event.Command, event.Success, event.Message)
			if event.Command == "quit" || event.Command == "exit" {
				goto Shutdown
			}

		case err := <-metricsErrCh:
			fmt.Printf("❌ 嚴重錯誤：Metrics Server 運行失敗: %v\n", err)
			goto Shutdown

		case sig := <-osSignalCh:
			fmt.Printf("\n🚨 收到操作系統信號 (%v)，準備執行優雅關閉。\n", sig)
			goto Shutdown
		}
	}

Shutdown:
	fmt.Println("--- 執行優雅關閉程序 ---")

	if err := metricsManager.GracefulShutdownMetricsServer(context.Background()); err != nil {
		fmt.Printf("❌ 關閉 Metrics Server 失敗: %v\n", err)
	}

	fmt.Println("👋 程式已終止。")
}
