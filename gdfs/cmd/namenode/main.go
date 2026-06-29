package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"gdfs/internal/cluster"
	"gdfs/internal/namenode"
)

func main() {
	addr := flag.String("addr", ":9000", "listen address")
	flag.Parse()

	node, err := namenode.NewNameNode(namenode.NewMetadataStore())
	if err != nil {
		log.Fatal(err)
	}

	server := namenode.NewHTTPServer(node)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node.StartHealthChecker(ctx, 5*time.Second, cluster.DefaultHealthConfig())

	log.Printf("starting namenode addr=%s", *addr)

	if err := http.ListenAndServe(*addr, server); err != nil {
		log.Fatal(err)
	}
}
