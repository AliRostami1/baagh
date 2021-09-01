package rf

// Closer is something that can be closed.
type Closer interface {
	// Close closes the thing.
	Close() error
}

// CodeTransmitter defines the interface for a rf code transmitter.
type CodeTransmitter interface {
	Closer
	// Transmit transmits a code using given protocol and pulse length.
	//
	// This method returns immediately. The code is transmitted in the background.
	// If you need to ensure that a code has been fully transmitted, wait for the
	// returned channel to be closed.
	Transmit(code uint64, protocol Protocol, pulseLength uint) <-chan struct{}
}

// CodeReceiver defines the interface for a rf code receiver.
type CodeReceiver interface {
	Closer
	// Receive blocks until there is a result on the receive channel.
	Receive() <-chan ReceiveResult
}

// ReceiveResult contains information about a detected code sent by an rf code
// transmitter.
type ReceiveResult struct {
	// Code is the detected code.
	Code uint64

	// BitLength is the detected bit length.
	BitLength uint

	// PulseLength is the detected pulse length.
	PulseLength int64

	// Protocol is the detected protocol. The protocol is 1-indexed.
	Protocol int
}
