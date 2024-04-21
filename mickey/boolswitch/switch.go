package boolswitch

import (
	"sync"
	"sync/atomic"
)

type Switch interface {
	Enable()
	Disable()
	Enabled() bool
	AddEnableCallback(f func())
	AddDisableCallback(f func())
}

type boolSwitch struct {
	v int32

	disableCallbackMu *sync.Mutex
	disableCallback   func()

	enableCallbackMu *sync.Mutex
	enableCallback   func()
}

func (t *boolSwitch) AddEnableCallback(f func()) {
	t.enableCallbackMu.Lock()
	t.enableCallback = f
	t.enableCallbackMu.Unlock()
}

func (t *boolSwitch) AddDisableCallback(f func()) {
	t.disableCallbackMu.Lock()
	t.disableCallback = f
	t.disableCallbackMu.Unlock()
}

func (t *boolSwitch) Enable() {
	atomic.StoreInt32(&t.v, 1)
	t.enableCallbackMu.Lock()
	t.enableCallback()
	t.enableCallbackMu.Unlock()
}

func (t *boolSwitch) Disable() {
	atomic.StoreInt32(&t.v, 0)
	t.disableCallbackMu.Lock()
	t.disableCallback()
	t.disableCallbackMu.Unlock()
}

func (t *boolSwitch) Enabled() bool {
	return atomic.LoadInt32(&t.v) == 1
}

func New(isOn bool) Switch {
	var v int32
	if isOn {
		v = 1
	}
	return &boolSwitch{
		v:                 v,
		disableCallbackMu: &sync.Mutex{},
		disableCallback:   func() {},
		enableCallbackMu:  &sync.Mutex{},
		enableCallback:    func() {},
	}
}

func NewEnabled() Switch {
	return New(true)
}

func NewDisabled() Switch {
	return New(false)
}
