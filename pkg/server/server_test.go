package server

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestGetGRPCConn(t *testing.T) {
	cases := []struct {
		name   string
		target func(string) string
	}{
		{name: "IP address", target: func(addr string) string { return addr }},
		{name: "hostname", target: func(addr string) string {
			_, port, _ := net.SplitHostPort(addr)
			return net.JoinHostPort("localhost", port)
		}},
		{name: "explicit DNS", target: func(addr string) string { return "dns:///" + addr }},
		{name: "explicit passthrough", target: func(addr string) string { return "passthrough:///" + addr }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			server := grpc.NewServer()
			grpc_health_v1.RegisterHealthServer(server, health.NewServer())
			t.Cleanup(server.Stop)
			go func() { _ = server.Serve(listener) }()
			conn, err := getGRPCConn(c.target(listener.Addr().String()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = conn.Close() })
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			response, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{}, grpc.WaitForReady(true))
			if err != nil {
				t.Fatal(err)
			}
			if response.Status != grpc_health_v1.HealthCheckResponse_SERVING {
				t.Fatalf("unexpected health status: %v", response.Status)
			}
		})
	}
}
