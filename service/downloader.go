package service

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type Downloader interface {
	Download(URI string) error
}

type DownloaderImpl struct {
	fileManager     FileManager
	progressTracker ProgressTracker
}

func NewDownloader(fileManager FileManager, progressTracker ProgressTracker) Downloader {
	return &DownloaderImpl{
		fileManager:     fileManager,
		progressTracker: progressTracker,
	}
}

func (d *DownloaderImpl) Download(URI string) error {
	sync.OnceFunc(func() {
		slog.Info("Downloading...")
	})

	file, err := d.fileManager.Create(URI)
	defer file.Close()

	resp, err := http.Get(URI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resourceSize := resp.Header.Get("Content-Length")

	if resourceSize == "" {
		d.fileManager.Delete(URI)
		return errors.New(fmt.Sprintf("\"%s\" was skipped. Invalid URI", URI))
	}

	d.progressTracker.SetTotalSize(resourceSize)
	d.progressTracker.SetFileName(file.Name())

	_, err = io.Copy(io.MultiWriter(file, d.progressTracker), resp.Body)
	if err != nil {
		return err
	}

	return nil
}
