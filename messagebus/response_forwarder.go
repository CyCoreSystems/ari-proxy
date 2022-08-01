package messagebus

import (
	"sync"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
)

type responseForwarder struct {
	closed   bool
	count    int
	expected int
	fwdChan  chan *proxy.Response

	mu sync.Mutex
}

func (f *responseForwarder) Forward(o *proxy.Response) {

	f.mu.Lock()
	defer f.mu.Unlock()

	f.count++

	if f.closed {
		return
	}

	// always send up reply, so we can track errors.
	select {
	case f.fwdChan <- o:
	default:
	}

	if f.count >= f.expected {
		f.closed = true
		close(f.fwdChan)
	}
}
