package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ledgerwatch/turbo-geth/log"
)

func RootContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		ch := make(chan os.Signal, 1)
		defer close(ch)

		signal.Notify(ch,
			os.Interrupt,
			os.Kill,
			syscall.CTRL_C_EVENT,
			syscall.CTRL_BREAK_EVENT,
			syscall.CTRL_CLOSE_EVENT,
			syscall.CTRL_LOGOFF_EVENT,
			syscall.CTRL_SHUTDOWN_EVENT,
			syscall.PROCESS_TERMINATE)
		defer signal.Stop(ch)

		select {
		case <-ch:
			log.Info("Got interrupt, shutting down...")
		case <-ctx.Done():
		}
	}()
	return ctx
}
