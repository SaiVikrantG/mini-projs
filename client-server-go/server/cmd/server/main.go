package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/SaiVikrantG/server/internal/handlers"
	"github.com/SaiVikrantG/server/internal/router"
	"github.com/SaiVikrantG/server/internal/server"
)

const DefaultContextTimeout = 2

func main() {
	port := flag.Int("port", 28333, "Port to start the server on")
	rootDir := flag.String("dir", "./", "Root directory to serve files from")

	flag.Parse()

	r := router.NewRouter()
	r.Handle("GET", "/file1.txt", &handlers.FileHandler{RootDir: *rootDir})
	r.Handle("GET", "/file2.html", &handlers.FileHandler{RootDir: *rootDir})

	srv := server.ServerInit(*port, r)
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
