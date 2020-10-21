package generate

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/turbo/torrent"
)

func SeedSnapshots(dir string) error {
	client := torrent.New(dir, torrent.SnapshotMode{
		Headers:  true,
		Bodies:   true,
		State:    false,
		Receipts: false,
	}, true)

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	db := ethdb.NewLMDB().Path(dir + "/tmpdb").MustOpen()
	err := client.Run(ethdb.NewObjectDatabase(db, ethdb.DefaultStateBatchSize))
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			for _, t := range client.Cli.Torrents() {
				log.Info("Snapshot stats", "snapshot", t.Name(), "active peers", t.Stats().ActivePeers, "seeding", t.Seeding())
			}
			if common.IsCanceled(ctx) {
				ticker.Stop()
				return
			}
		}
	}()

	<-ctx.Done()
	return nil
}
