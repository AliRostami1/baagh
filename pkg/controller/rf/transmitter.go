package rf

import (
	"sync/atomic"
	"time"

	"github.com/AliRostami1/baagh/pkg/controller/core"
)

const (
	// DefaultTransmissionCount defines how many times a code should be
	// transmitted in a row by default.
	DefaultTransmissionCount = 10

	transmissionChanLen = 32
	bitLength           = 24
)

type transmission struct {
	code        uint64
	protocol    Protocol
	pulseLength uint
	done        chan struct{}
}

// Transmitter can serialize and transmit rf codes.
type Transmitter struct {
	item              core.Item
	transmission      chan transmission
	closed            int32
	transmissionCount int
	// delay can be replaced in tests to make timing predictable.
	delay func(time.Duration)
}

// NewTransmitter creates a Transmitter which attaches to the chip's pin at
// offset.
func NewTransmitter(chip string, offset int, options ...TransmitterOption) (*Transmitter, error) {
	i, err := core.RequestItem(chip, offset, core.AsOutput(core.StateInactive))
	if err != nil {
		return nil, err
	}

	return NewPinTransmitter(i, options...), nil
}

// NewPinTransmitter creates a *Transmitter that sends on pin.
func NewPinTransmitter(item core.Item, options ...TransmitterOption) *Transmitter {
	t := &Transmitter{
		item:              item,
		transmission:      make(chan transmission, transmissionChanLen),
		closed:            0,
		transmissionCount: DefaultTransmissionCount,
		delay:             delay,
	}

	for _, option := range options {
		option(t)
	}

	if t.transmissionCount <= 0 {
		t.transmissionCount = 1
	}

	go t.watch()

	return t
}

// Transmit transmits a code using given protocol and pulse length.
//
// This method returns immediately. The code is transmitted in the background.
// If you need to ensure that a code has been fully transmitted, wait for the
// returned channel to be closed.
func (t *Transmitter) Transmit(code uint64, protocol Protocol, pulseLength uint) <-chan struct{} {
	done := make(chan struct{})

	if atomic.LoadInt32(&t.closed) == 1 {
		close(done)
		return done
	}

	t.transmission <- transmission{
		code:        code,
		protocol:    protocol,
		pulseLength: pulseLength,
		done:        done,
	}

	return done
}

// transmit performs the acutal transmission of the remote control code.
func (t *Transmitter) transmit(trans transmission) {
	defer close(trans.done)

	for i := 0; i < t.transmissionCount; i++ {
		for j := bitLength - 1; j >= 0; j-- {
			if trans.code&(1<<uint64(j)) > 0 {
				t.send(trans.protocol.One, trans.pulseLength)
			} else {
				t.send(trans.protocol.Zero, trans.pulseLength)
			}
		}
		t.send(trans.protocol.Sync, trans.pulseLength)
	}
}

// Close stops started goroutines and closes the gpio pin.
func (t *Transmitter) Close() error {
	atomic.StoreInt32(&t.closed, 0)
	close(t.transmission)
	return t.item.Close()
}

// watch listens on a channel and processes incoming transmissions.
func (t *Transmitter) watch() {
	for trans := range t.transmission {
		t.transmit(trans)
	}
}

// send sends a sequence of high and low pulses on the gpio pin.
func (t *Transmitter) send(pulses HighLow, pulseLength uint) {
	t.item.SetState(core.StateActive)
	t.delay(time.Microsecond * time.Duration(pulseLength*pulses.High))
	t.item.SetState(core.StateInactive)
	t.delay(time.Microsecond * time.Duration(pulseLength*pulses.Low))
}
