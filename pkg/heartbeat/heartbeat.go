package heartbeat

import (
	"fmt"
	"os"
	"time"
)

type Heartbeat interface {
	Run()
	Down()

	// TestFreeze only for test reboot function
	TestFreeze()
}

type heartbeat struct {
	serviceName   string
	fileDir       string
	writeInterval time.Duration
	down          chan struct{}

	// only for test reboot function
	_freeze chan struct{}
}

func New(serviceName, fileDir string, writeInterval time.Duration) Heartbeat {
	return &heartbeat{
		serviceName:   serviceName,
		fileDir:       fileDir,
		writeInterval: writeInterval,
		down:          make(chan struct{}),
		_freeze:       make(chan struct{}),
	}
}

func (h *heartbeat) Run() {
	if _, err := os.Stat(h.fileDir); os.IsNotExist(err) {
		if err = os.Mkdir(h.fileDir, 0777); err != nil {
			panic(err)
		}
	}

	logsFilePath := fmt.Sprintf("%s/%s.txt", h.fileDir, h.serviceName)

	writeHeartbeat := func() {
		ticker := time.NewTicker(h.writeInterval)

		for {
			select {
			case <-h.down:
				return
			case <-ticker.C:
				_ = os.WriteFile(logsFilePath, []byte("0"), 0666)
			case <-h._freeze:
				time.Sleep(120 * time.Second)
			}
		}
	}

	go writeHeartbeat()
}

func (h *heartbeat) Down() {
	h.down <- struct{}{}
}

func (h *heartbeat) TestFreeze() {
	h._freeze <- struct{}{}
}
