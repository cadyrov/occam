package cmd

import (
	"context"
	"github.com/cadyrov/occam/internal"
	"github.com/cadyrov/occam/providers/storage"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// nolint: gochecknoglobals
var (
	cliCmd = &cobra.Command{
		Use:   "cli",
		Short: "Run app as daemon",
		Long:  `Run app as daemon`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

// nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
}

func run() {
	log := initLogger()

	log.Debug().Interface("values", cnf).Msg("config")

	keeper := storage.New(&log, cnf.Project.Shift)
	if keeper == nil {
		panic("keeper cant be created")
	}

	srv := internal.New(&log, keeper)

	if srv == nil {
		panic("service cant be created")
	}

	ctx := context.Background()

	stopWaitGroup := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	stopWaitGroup.Add(1)

	go func() {
		defer stopWaitGroup.Done()

		srv.Run(ctx)
	}()

	stopWaitGroup.Add(1)

	SigTermHandler(func() {
		defer stopWaitGroup.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown")
		}

		cancel()
	})

	stopWaitGroup.Wait()
}

const maxTermChanLen = 10

func SigTermHandler(stopFunc func()) {
	termCh := make(chan os.Signal, maxTermChanLen)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-termCh

		stopFunc()
	}()
}
