// Strip the path off for security reasons.
// Read the data from the named file.
// Determine the type of data in the file, HTML or text.
// Build an HTTP response packet with the file data in the payload.
// Send that HTTP response back to the client.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/SaiVikrantG/server/internal/server"
)

const DefaultContextTimeout = 2

func main() {
	port := flag.Int("port", 28333, "Port to start the server on")

	flag.Parse()

	srv := server.ServerInit(*port)
	if err := srv.ServerStart(); err != nil {
		fmt.Println("Server failed to start with the following error: ", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go srv.ServerListen(ctx)

	<-ctx.Done()

	shutDownctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)
	defer cancel()

	if err := srv.ServerStop(shutDownctx); err != nil {
		fmt.Println("Server failed to stop with error: ", err)
	}
}
