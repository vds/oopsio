package receiver_test

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/vds/oops"
	"github.com/vds/oopsio/receiver"
)

// func TestRun(t *testing.T) {
// 	a := "127.0.0.1:5678"
// 	q := make(chan struct{})
// 	r := receiver.NewReceiver(q, "")
// 	laddr, err := net.ResolveTCPAddr("tcp", a)
// 	if nil != err {
// 		t.Fatalf(err.Error())
// 	}
// 	go r.Run(laddr)
// 	runtime.Gosched()
// 	conn, err := net.Dial("tcp", a)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	defer conn.Close()
// 	ra := conn.RemoteAddr()
// 	if ra.String() != a {
// 		t.Fatalf("Remote address shoud be: %s, but is: %s", a, ra)
// 	}
// }

func TestReceiveOops(t *testing.T) {
	d, err := ioutil.TempDir("/tmp/", "oops")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(d)
		if nil != err {
			t.Fatalf(err.Error())
		}
	}()

	fs, err := ioutil.ReadDir(d)
	if nil != err {
		t.Fatalf(err.Error())
	}
	if len(fs) != 0 {
		t.Fatalf("Oops directory not empty.")
	}

	a := "127.0.0.1:5678"
	laddr, err := net.ResolveTCPAddr("tcp", a)
	if nil != err {
		t.Fatalf(err.Error())
	}

	q := make(chan struct{})
	r := receiver.NewReceiver(q, d)

	r.Wg.Add(1)
	go func() {
		l, err := net.Listen("tcp", laddr.String())
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go r.ReceiveOops(conn)
	}()

	runtime.Gosched()
	o := oops.Oops{}
	o.SetError(errors.New("oops"), false)
	b, err := o.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	cconn, err := net.Dial("tcp", a)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cconn.Write(b)
	if err != nil {
		t.Fatal(err)
	}
	cconn.Close()

	r.Wg.Wait()
	fs, err = ioutil.ReadDir(d)
	if nil != err {
		t.Fatalf(err.Error())
	}
	if len(fs) != 1 {
		t.Fatalf("Oops directory empty.")
	}
}
