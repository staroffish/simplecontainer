package process

import (
	"os"
	"os/signal"
	"syscall"
)

func ContianerProcess(name string) {
	sigCh := make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigCh:
			if sig == syscall.SIGINT {
				initProcess()
			} else {
				Cleanup(name)
				os.Exit(-1)
			}
		}
	}
}

func initProcess() {

}

func Cleanup(name string) {

}
