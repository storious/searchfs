package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"gdfs/internal/cluster"
	"gdfs/internal/datanode"
	"gdfs/internal/namenode"
)

func main() {
	var (
		id           = flag.String("id", "node-1", "datanode id")
		addr         = flag.String("addr", ":9001", "listen address")
		root         = flag.String("root", "data/datanode", "storage root")
		namenodeAddr = flag.String("namenode", "", "namenode address")
	)
	flag.Parse()

	store := datanode.NewLocalBlockStore(*root)

	node, err := datanode.NewDataNode(
		datanode.NodeID(*id),
		*addr,
		store,
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registryAddr := "http://localhost" + *addr

	startHeartbeat(
		ctx,
		*namenodeAddr,
		node,
		registryAddr,
		5*time.Second,
	)

	server := datanode.NewHTTPServer(node)

	log.Printf("starting datanode id=%s addr=%s root=%s", *id, *addr, *root)

	if err := http.ListenAndServe(*addr, server); err != nil {
		log.Fatal(err)
	}
}
func startHeartbeat(
	ctx context.Context,
	namenodeAddr string,
	node *datanode.DataNode,
	addr string,
	interval time.Duration,
) {
	if namenodeAddr == "" {
		return
	}

	client := namenode.NewHTTPClient(namenodeAddr)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		send := func() {
			stats, err := node.Stats()
			if err != nil {
				log.Printf("collect stats failed: %v", err)
				return
			}

			err = client.Heartbeat(ctx, cluster.Heartbeat{
				ID:       cluster.DataNodeID(node.ID),
				Addr:     addr,
				Capacity: stats.Capacity,
				Used:     stats.Used,
			})
			if err != nil {
				log.Printf("heartbeat failed: %v", err)
				return
			}

			log.Printf("heartbeat sent id=%s used=%d capacity=%d", node.ID, stats.Used, stats.Capacity)
		}

		send()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				send()
			}
		}
	}()
}
