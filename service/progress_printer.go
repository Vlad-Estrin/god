package service

import (
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	tm "github.com/muesli/termenv"
)

type ProgressPrinter interface {
	Start()
	Update(fileName string, progress int)
	Stop(wg *sync.WaitGroup)
}

type ProgressPrinterImpl struct {
	fileProgress *sync.Map
	ticker       *time.Ticker
	out          *tm.Output
	completed    []string
}

func NewPrinter() ProgressPrinter {
	return &ProgressPrinterImpl{
		fileProgress: &sync.Map{},
		ticker:       time.NewTicker(time.Second / 10),
		out:          tm.NewOutput(os.Stdout),
		completed:    []string{},
	}
}

func (p *ProgressPrinterImpl) Start() {
	p.out.AltScreen()
	startLine := 1
	fileLine := sync.Map{}

	for range p.ticker.C {
		p.fileProgress.Range(func(key, value interface{}) bool {
			keyStr := key.(string)

			_, ok := fileLine.Load(keyStr)
			if !ok {
				fileLine.Store(keyStr, startLine)
				startLine++
			}

			lineToPrint, ok := fileLine.Load(keyStr)
			if ok {
				p.out.MoveCursor(lineToPrint.(int), 0)
			}

			fmt.Printf("\r%s: %d%%", key, value)

			return true
		})
	}
}

func (p *ProgressPrinterImpl) Update(fileName string, progress int) {
	p.fileProgress.Store(fileName, progress)
}
func (p *ProgressPrinterImpl) Stop(wg *sync.WaitGroup) {
	for {
		if p.canBeStopped() {
			p.out.ExitAltScreen()
			fmt.Printf("\n- Completed! Downloaded %d file(s)\n", len(p.completed))
			p.ticker.Stop()
			wg.Done()

			return
		}
	}
}

func (p *ProgressPrinterImpl) canBeStopped() bool {
	result := true
	p.fileProgress.Range(func(key, value interface{}) bool {
		keyStr := key.(string)

		if value != 100 {
			result = false
			return false
		} else {
			if !slices.Contains(p.completed, keyStr) {
				p.completed = append(p.completed, keyStr)
			}

			return true
		}
	})
	return result
}
