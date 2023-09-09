package service

import (
	"strconv"
)

type ProgressTracker interface {
	Write(data []byte) (int, error)
	SetTotalSize(size string)
	SetFileName(fileName string)
}

type ProgressTrackerImpl struct {
	progressPrinter ProgressPrinter
	fileName        string
	fileSize        int
	downloaded      int
}

func NewProgressTracker(progressPrinter ProgressPrinter) ProgressTracker {
	return &ProgressTrackerImpl{
		progressPrinter: progressPrinter,
	}
}

func (p *ProgressTrackerImpl) Write(data []byte) (int, error) {
	p.downloaded += len(data)
	progress := p.downloaded * 100 / p.fileSize
	p.progressPrinter.Update(p.fileName, progress)
	return len(data), nil
}

func (p *ProgressTrackerImpl) SetTotalSize(size string) {
	p.fileSize, _ = strconv.Atoi(size)
}

func (p *ProgressTrackerImpl) SetFileName(fileName string) {
	p.fileName = fileName
}
