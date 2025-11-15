package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsManager 結構：用於封裝和管理 HTTP 伺服器
type MetricsManager struct {
	// metricsServer 變數現在是結構的私有成員，只能透過結構方法存取
	server *http.Server
}

// NewMetricsManager 構造函式
func NewMetricsManager(addr string) *MetricsManager {
	// 使用 := 宣告
	s := &http.Server{
		Addr:         addr,
		Handler:      promhttp.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 使用 := 宣告
	manager := &MetricsManager{
		server: s,
	}
	return manager
}

var (
	blocksMined = promauto.NewCounter(prometheus.CounterOpts{
		Name: "blocks_mined_total",
		Help: "Total number of blocks mined by this node",
	})
	chainHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chain_block_height",
		Help: "Current blockchain height",
	})
)

func (m *MetricsManager) RunMetricsServer(errCh chan<- error) {
	fmt.Printf("✅ Metrics server listening on %s\n", m.server.Addr)

	// 使用 := 宣告，並檢查伺服器是否已經存在
	if m.server == nil {
		errCh <- fmt.Errorf("metrics server instance is nil")
		return
	}

	// 在 Goroutine 中啟動伺服器
	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// 只有當錯誤不是 "Server closed" 時，才發送錯誤到 channel
		errCh <- fmt.Errorf("metrics server error: %w", err)
	}
	// 注意：http.ErrServerClosed 在優雅關閉時會被觸發，這不是一個真正的錯誤
}

func (m *MetricsManager) GracefulShutdownMetricsServer(ctx context.Context) error {
	if m.server == nil {
		return nil
	}

	fmt.Println("⏳ Shutting down metrics server gracefully...")

	// 使用 := 宣告
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Shutdown() 現在直接作用於結構內的 m.server
	if err := m.server.Shutdown(shutdownCtx); err != nil {
		// 使用 := 宣告
		return fmt.Errorf("metrics server shutdown error: %w", err)
	}

	fmt.Println("✅ Metrics server shut down.")
	return nil
}
