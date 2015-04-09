package receiver

import (
	"io"
	"log"
	"net"
	"os"
	"path"
	"sync"
	"time"
)

// NewReceiver is a Receiver factory
func NewReceiver(quit chan struct{}, path string) Receiver {
	return Receiver{Quit: quit, path: path}
}

// Receiver implements a simple TCP service to receive Oopses
type Receiver struct {
	Wg   sync.WaitGroup
	Quit chan struct{}
	path string
}

// Run start the Receiver service
func (r *Receiver) Run(laddr *net.TCPAddr) {
	log.Println("Running")
	l, err := net.Listen("tcp", laddr.String())
	if err != nil {
		// FIXME this exits the process calling os.Exit(1),
		// it would be nice to have only one exit point but this will be used as goroutine.
		// What's the best way to handle this fatal error?
		log.Fatal(err)
	}
	defer func() {
		l.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	// FIXME what kind of back-pressure mechanism can be used to warm the dispatchers when under heavy load?
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		select {
		case <-r.Quit:
			return
		default:
		}
		r.Wg.Add(1)
		go r.ReceiveOops(conn)
	}
}

// ReceiverOops is the TPC connection handler that handles the reception of oopses
func (r *Receiver) ReceiveOops(c net.Conn) {
	log.Println("Receiving")
	defer func() {
		r.Wg.Done()
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	f, err := os.Create(path.Join(r.path, string(time.Now().UnixNano())+".oops"))
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	// FIXME what happens if there is an error when receiving the Oops?
	// How do we notify the dispatcher?
	_, err = io.Copy(f, c)
	if err != nil {
		log.Printf("error handling request: %s", err)
	}
}

// Stop stops the receptions of oops
func (r *Receiver) Stop() error {
	r.Wg.Wait()
	return nil
}
