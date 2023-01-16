package common

import (
	"context"
	"errors"
	"github.com/lunny/log"
	"sync"
	"time"
)

var NotUsed = Pack{}
var eventMap = sync.Map{}

type Event struct {
	//监听者
	waiters []chan *Pack
	//结果
	result *Pack
	//上下文控制
	ctxBg     context.Context
	ctxCancel context.CancelFunc
}

type Pack struct {
	Content interface{}
	Err     error
}

func NewEvent(key string) *Event {
	e := &Event{}
	e.Reset()
	return e
}

func AddEvent(key string, e *Event) {
	eventMap.Store(key, e)
}

func RemoveEvent(key string) {
	eventMap.Delete(key)
}

func GetEvent(key string) *Event {
	event, ok := eventMap.Load(key)
	if ok {
		return event.(*Event)
	}
	return nil
}

func (e *Event) AddWaiter() *chan *Pack {
	//等待者
	resultChan := make(chan *Pack, 0)
	e.waiters = append(e.waiters, resultChan)
	return &resultChan
}

// 等待结果
func (e *Event) Wait(waiter *chan *Pack, timeout time.Duration) (*Pack, error) {
	if e.result == &NotUsed {
		ctx, cancel := context.WithTimeout(e.ctxBg, time.Second)
		defer cancel()

		//等待
		select {
		case result := <-*waiter:
			return result, nil
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				return nil, errors.New("context.Canceled")
			}
			//log.Warnf("event wait timeout")
			return nil, errors.New("event wait timeout")
		}
	} else {
		return e.result, nil
	}
}

// 发送结果
func (e *Event) Send(result *Pack) error {
	//防止发送多次
	if e.result != &NotUsed {
		return errors.New("event is used")
	}

	ctx, cancel := context.WithTimeout(e.ctxBg, time.Second*1000)
	defer cancel()

	for _, resultChan := range e.waiters {
		select {
		case resultChan <- result:
		case <-ctx.Done():
			//log.Warnf("Event.Send %p resultChan=%d", e, len(resultChan))
			return errors.New("event wait timeout")
		}
	}
	e.result = result
	return nil
}

// 发送结果
func Send(key string, result *Pack) {
	e := GetEvent(key)
	if e != nil {
		err := e.Send(result)
		if err != nil {
			log.Warn(err.Error())
		}
	}
}

// 重置
func (e *Event) Reset() {
	if e.ctxBg != nil {
		e.ctxCancel()
	}

	e.ctxBg, e.ctxCancel = context.WithCancel(context.Background())
	e.waiters = nil
	e.result = &NotUsed
}
