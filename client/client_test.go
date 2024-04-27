package client

import (
	"context"
	"fmt"
	"redis/server"
	"strconv"
	"testing"
)

func TestClient_Set(t *testing.T) {
	srv := server.New(server.Config{})
	go func() {
		err := srv.Start()
		if err != nil {
			t.Fatalf("error while starting server: %e", err)
		}
	}()

	for i := 0; i < 5; i++ {
		go func(i int) {
			client, err := New("localhost:8080")
			if err != nil {
				t.Fatalf("error while creating client %e", err)
			}

			key := fmt.Sprintf("foo_%d", i)
			value := fmt.Sprintf("bar_%d", i)

			if err := client.Set(context.Background(), key, value); err != nil {
				t.Fatalf("error while using SET: %e", err)
			}
			str, err := client.Get(context.Background(), key)
			if err != nil {
				t.Fatalf("error while using GET: %e", err)
			}

			if str != value {
				t.Fatalf("got %q, wanted %q", str, value)
			}

			err = client.Close()
			if err != nil {
				t.Fatalf("error while using GET: %e", err)
			}
		}(i)

	}

	if len(srv.Clients) != 0 {
		t.Fatalf("not all client were disconnected: %d", len(srv.Clients))
	}

	srv.Stop()
}

func TestClientInt_Set(t *testing.T) {
	srv := server.New(server.Config{})
	go func() {
		err := srv.Start()
		if err != nil {
			t.Fatalf("error while starting server: %e", err)
		}
	}()

	for i := 0; i < 5; i++ {
		go func(i int) {
			client, err := New("localhost:8080")
			if err != nil {
				t.Fatalf("error while creating client %e", err)
			}

			key := fmt.Sprintf("foo_%d", i)
			value := i

			if err := client.Set(context.Background(), key, value); err != nil {
				t.Fatalf("error while using SET: %e", err)
			}
			strInt, err := client.Get(context.Background(), key)
			if err != nil {
				t.Fatalf("error while using GET: %e", err)
			}

			valInt, err := strconv.Atoi(strInt)

			if err != nil {
				t.Fatalf("error while using strconv.Atoi: %e", err)
			}

			if valInt != value {
				t.Fatalf("got %q, wanted %q", valInt, value)
			}

			err = client.Close()
			if err != nil {
				t.Fatalf("error while using GET: %e", err)
			}
		}(i)

	}

	if len(srv.Clients) != 0 {
		t.Fatalf("not all client were disconnected: %d", len(srv.Clients))
	}
	srv.Stop()
}
