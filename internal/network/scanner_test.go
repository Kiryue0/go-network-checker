package network

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestScanPort_Open(t *testing.T) {
	// test server aç
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	// port numarasını al
	port := ln.Addr().(*net.TCPAddr).Port

	ctx := context.Background()
	result := ScanPort(ctx, "127.0.0.1", port, 2*time.Second)

	if !result.IsOpen {
		t.Errorf("expected port to be open")
	}
}

func TestScanPort_Closed(t *testing.T) {
	ctx := context.Background()
	result := ScanPort(ctx, "127.0.0.1", 19999, 2*time.Second)

	if result.IsOpen {
		t.Errorf("expected port to be closed")
	}
}
