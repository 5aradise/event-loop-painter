package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	curr screen.Texture

	state *textureState

	mq messageQueue

	stop chan struct{}
}

var (
	size             = image.Pt(800, 800)
	MessageQueueSize = 1 << 10
)

func NewLoop() *Loop {
	return &Loop{
		mq:   messageQueue{make(chan any, MessageQueueSize)},
		stop: make(chan struct{}),
	}
}

func (l *Loop) Start(s screen.Screen) {
	l.curr, _ = s.NewTexture(size)

	l.state = newTextureState()

loop:
	for {
		switch msg := l.mq.pull().(type) {
		case Operation:
			update := msg.Do(l.state)
			if update {
				l.state.set(l.curr)
				l.Receiver.Update(l.curr)
				l.curr, _ = s.NewTexture(size)
				l.state = newTextureState()
			}

		case closeSignal:
			break loop

		default:
			panic("message in messageQueue not Operation or closeSignal")
		}
	}
	close(l.stop)
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

type closeSignal struct{}

func (l *Loop) StopAndWait() {
	l.mq.push(closeSignal{})
	<-l.stop
}

type messageQueue struct {
	buf chan any
}

func (mq *messageQueue) push(v any) {
	mq.buf <- v
}

func (mq *messageQueue) pull() any {
	return <-mq.buf
}
