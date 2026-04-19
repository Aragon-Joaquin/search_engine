package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"search_engine/internal/blobs"
	"search_engine/internal/repository"
	"search_engine/internal/utils"

	"github.com/charmbracelet/ssh"
)

var (
	systemThreads = runtime.NumCPU()
	maxTimeout    = time.Second * 3

	HOST         = utils.GetEnv(utils.ENV_HOST)
	PORT         = utils.GetEnv(utils.ENV_PORT)
	KEYHOST_PATH = utils.GetEnv(utils.ENV_KEYHOST)

	start time.Time
)

type Application struct {
	rep    *repository.Repository
	server *ssh.Server
}

var app *Application

func init() {
	start = time.Now()
}

func main() {
	if HOST == "" || PORT == "" {
		panic("host and/or port are not defined in the .env")
	}

	if KEYHOST_PATH == "" {
		panic("keyhost path is not defined. the default should be .ssh/id_ed25519")
	}

	// loads all blobs again into the redisDB
	flagLoadBlobs := flag.Bool("l", false, "requires a bool")
	flag.Parse()

	if flagLoadBlobs != nil && *flagLoadBlobs {
		log.Println("USED FLAG -l - LOADING ALL ./data/* BLOBS TO REDIS")
		loadBlobs()
		log.Println("UPLOAD FINISHED IN: ", time.Since(start))
		return
	}

	// NOTE: else, we boot up the ssh server
	s, err := initServer()
	if err != nil {
		panic(err)
	}

	app = &Application{
		rep:    repository.CreateRepostory(DBRedis),
		server: s,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Starting SSH server", "host", HOST, "port", PORT)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalln("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Println("Closing SSH Server...")
}

func loadBlobs() {
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()
	blobList := blobs.LoadBlobsFromFolder()

	var wg sync.WaitGroup
	waitChan := make(chan any)

	go func() {
		for _, blob := range blobList.Blobs {
			wg.Go(func() {
				if err := DBRedis.AddZSort(ctx, blob); err != nil {
					log.Println("Error in one of the blobs while trying to load it to redis: ", err)
				}
			})
		}
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		return
	case <-ctx.Done():
		panic("timeout'ed while loading all blobs from local to redis")
	}
}
