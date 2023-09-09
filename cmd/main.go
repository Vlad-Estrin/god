package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"god/service"
)

func main() {
	wg := &sync.WaitGroup{}
	var errs []error

	validator := service.NewValidator()

	args := os.Args
	resources, err := validator.GetValid(args[1:])
	if err != nil {
		errs = append(errs, err)
	}

	fileManager := service.NewFileManager()
	progressPrinter := service.NewPrinter()

	go progressPrinter.Start()

	mu := sync.Mutex{}
	for i := 0; i < len(resources); i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			progressTracker := service.NewProgressTracker(progressPrinter)
			downloader := service.NewDownloader(fileManager, progressTracker)
			err = downloader.Download(resources[idx])
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	wg.Add(1)
	progressPrinter.Stop(wg)
	wg.Wait()

	for _, err = range errs {
		slog.Error(fmt.Sprintf("%s\n", err.Error()))
	}
}
