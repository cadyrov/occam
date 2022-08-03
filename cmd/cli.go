package cmd

import (
	"context"
	"github.com/cadyrov/occam/internal"
	"github.com/cadyrov/occam/providers/origin"
	"github.com/cadyrov/occam/providers/storage"
	"github.com/spf13/cobra"
	"io"
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

	ctx := context.Background()

	stopWaitGroup := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	log.Debug().Interface("values", cnf).Msg("config")

	precision := cnf.Project.PrecisionSecond
	if precision <= 0 {
		precision = 60
	}

	keeper := storage.New(&log, cnf.Project.Shift, precision)
	if keeper == nil {
		panic("keeper cant be created")
	}

	var w io.Writer
	if cnf.Project.Output != "" {
		f, err := os.Create(cnf.Project.Output)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		w = f
	} else {
		w = os.Stdout
	}

	origins := make([]internal.PriceStreamSubscriber, 0, 100)

	for i := 0; i < 100; i++ {
		mo := origin.New(&log)
		mo.Start(ctx)

		origins = append(origins, mo)
	}

	srv := internal.New(&log, keeper, origins, precision)

	if srv == nil {
		panic("service cant be created")
	}

	stopWaitGroup.Add(1)

	go func() {
		defer stopWaitGroup.Done()

		srv.Run(ctx, w)
	}()

	stopWaitGroup.Add(1)

	SigTermHandler(func() {
		defer stopWaitGroup.Done()

		log.Info().Msg("server shutdown")

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
