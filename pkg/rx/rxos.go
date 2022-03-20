package rx

import (
	"github.com/reactivex/rxgo/v2"
	"os"
	"os/signal"
	"syscall"
)

func OsSignalTermProducer() rxgo.Observable {
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	ch := make(chan rxgo.Item)
	go func() {
		<-wait
		ch <- rxgo.Of("done")
	}()
	return rxgo.FromChannel(ch)
}
