package server

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name         string
		port         string
		addr         string
		hasHandler   bool
		idleTimeout  time.Duration
		readTimeout  time.Duration
		writeTimeout time.Duration
	}{
		{
			name:         "invalid port",
			port:         "invalid",
			addr:         ":0",
			hasHandler:   true,
			idleTimeout:  time.Minute,
			readTimeout:  10 * time.Second,
			writeTimeout: 30 * time.Second,
		},
		{
			name:         "valid port",
			port:         "8080",
			addr:         ":8080",
			hasHandler:   true,
			idleTimeout:  time.Minute,
			readTimeout:  10 * time.Second,
			writeTimeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPort := os.Getenv("PORT")
			os.Setenv("PORT", tt.port)
			defer os.Setenv("PORT", originalPort)

			server := NewServer()

			if diff := cmp.Diff(tt.addr, server.Addr); diff != "" {
				t.Errorf("Test %q failed (Addr) (-want +got):\n%s", tt.name, diff)
			}

			if server.Handler == nil {
				t.Errorf("Test %q failed (Handler): server.Handler is nil\n", tt.name)
			}

			if diff := cmp.Diff(tt.idleTimeout, server.IdleTimeout); diff != "" {
				t.Errorf("Test %q failed (IdleTimeout) (-want +got):\n%s", tt.name, diff)
			}

			if diff := cmp.Diff(tt.readTimeout, server.ReadTimeout); diff != "" {
				t.Errorf("Test %q failed (ReadTimeout) (-want +got):\n%s", tt.name, diff)
			}

			if diff := cmp.Diff(tt.writeTimeout, server.WriteTimeout); diff != "" {
				t.Errorf("Test %q failed (WriteTimeout) (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}
