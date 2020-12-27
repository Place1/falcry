package webhook

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Broadcaster struct {
	channels sync.Map
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		channels: sync.Map{},
	}
}

func (b *Broadcaster) Channel(ctx context.Context) FalcoEventChannel {
	logrus.Debug("broadcast channel created")

	// We allow a small buffer so that slight differences
	// in channel processing time don't block up the system.
	// If a reader is very slow then we'll still have an issue
	// but the goal is to keep the system performant enough
	// to avoid this case.
	channel := make(FalcoEventChannel, 2)
	b.channels.Store(channel, true)

	go func() {
		logrus.Debug("broadcast channel deleted")
		<-ctx.Done()
		b.channels.Delete(channel)
	}()

	return channel
}

func (b *Broadcaster) Send(event *FalcoEvent) {
	b.channels.Range(func(key interface{}, value interface{}) bool {
		channel := key.(FalcoEventChannel)
		channel <- event
		return true
	})
}
