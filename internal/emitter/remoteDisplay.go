// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

package emitter

import (
	"fmt"
	"log"
	"net/rpc"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rokath/trice/pkg/msg"
)

// RemoteDisplay is transferring to a remote display object.
type RemoteDisplay struct {
	Err    error
	Cmd    string      // remote server executable
	Params string      // remote server additional parameters (despite "ds -ipa a.b.c.d -ipp nnnn")
	IPAddr string      // IP addr
	IPPort string      // IP port
	PtrRPC *rpc.Client // PtrRPC is a pointer for remote calls valid after a successful rpc.Dial()
}

// NewRemoteDisplay creates a connection to a remote Display and implements the Linewriter inteface.
// It accepts 0 to 4 string arguments. More arguments are ignored.
// For not given parameters default values are taken. The parameters are in the following order.
// args[0] (exe), is a programm started to create a remote server instance if not already running.
// If the remote server is already running on ips:ipp than a start of a 2nd instace is not is possible. This is silently ignored.
// If args[0] is empty, a running display server is assumed and a connection is established if possible.
// args[1] (params) contains additional remote display (=trice) command parameters.
// This value is used only if the remote server gets started.
// args[2] (ipa) is the IP address to be used to connect to the remote display.
// args[3] (ipp) is the IP port to be used to connect to the remote display.
func NewRemoteDisplay(args ...string) *RemoteDisplay {
	args = append(args, "", "", "", "") // make sure to have at least 4 elements in args.
	p := &RemoteDisplay{
		Err:    nil,
		Cmd:    "",        // per default assume remote display is already active.
		Params: "-lg off", // no logfile as default when remote display is launched.
		IPAddr: IPAddr,
		IPPort: IPPort,
		PtrRPC: nil,
	}
	if "" != args[0] {
		p.Cmd = args[0]
	}
	if "" != args[1] {
		p.Params = args[1]
	}
	if "" != args[2] {
		p.IPAddr = args[2]
	}
	if "" != args[3] {
		p.IPPort = args[3]
	}
	if "" != p.Cmd {
		p.startServer()
	}
	p.Connect()
	return p
}

// ErrorFatal ends in osExit(1) if p.Err not nil.
func (p *RemoteDisplay) ErrorFatal() {
	if nil == p.Err {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	log.Fatal(p.Err, filepath.Base(file), line)
}

// writeLine is implementing the Linewriter interface for RemoteDisplay.
func (p *RemoteDisplay) writeLine(line []string) {
	p.ErrorFatal()
	p.Err = p.PtrRPC.Call("Server.WriteLine", line, nil) // TODO: Change to "Server.WriteLine"
}

// startServer starts a display server with the filename exe (if not already running).
func (p *RemoteDisplay) startServer() {
	var cmd *exec.Cmd
	s := strings.Fields("ds -ipa " + p.IPAddr + " -ipp " + p.IPPort + " " + p.Params)
	cmd = exec.Command(p.Cmd, s...) // ... expands slice into individual string arguments
	go func() {
		p.Err = cmd.Run()
		p.ErrorFatal()
	}()
	time.Sleep(1000 * time.Millisecond)
}

// Connect is called by the client and tries to dial.
// On success PtrRpc is valid afterwards and the output is re-directed.
// Otherwise an error code is stored inside remotDisplay.
func (p *RemoteDisplay) Connect() {
	addr := p.IPAddr + ":" + p.IPPort
	if nil != p.PtrRPC {
		if Verbose {
			fmt.Println("already connected", p.PtrRPC)
		}
		return
	}
	if Verbose {
		fmt.Println("dialing " + addr + " ...")
	}
	p.PtrRPC, p.Err = rpc.Dial("tcp", addr)
	msg.FatalOnErr(p.Err)
	//p.ErrorFatal()
	if Verbose {
		fmt.Println("...remoteDisplay @ " + addr + " connected.")
	}
}

// ScShutdownRemoteDisplayServer connects to a client to send shutdown message to display server.
// 0==timeStamp -> no timestamp in shutdown message, use for tests.
// 1==timeStamp -> timestamp in shutdown message, use normally.
// It accepts 0 to 2 string arguments as args. More arguments are ignored.
// For not given parameters default values are taken. The parameters are in the following order.
// args[0] (ipa) is the IP address to be used to connect to the remote display.
// args[1] (ipp) is the IP port to be used to connect to the remote display.
func ScShutdownRemoteDisplayServer(timeStamp int64, args ...string) error {
	args = append(args, "", "") // make sure to have at least 2 elements in args.
	p := NewRemoteDisplay("", "", args[0], args[1])
	if nil == p.PtrRPC {
		p.Connect()
	}
	p.stopServer(timeStamp)
	return p.Err
}

// StopServer sends signal to display server to quit.
// `ts` is used as flag. If 1 shutdown message is with timestamp (default usage), if 0 shutdown message is without timestamp (for testing).
func (p *RemoteDisplay) stopServer(ts int64) {
	if Verbose {
		fmt.Println("sending Server.Shutdown...")
	}
	p.Err = p.PtrRPC.Call("Server.Shutdown", []int64{ts}, nil) // if 1st param nil -> gob: cannot encode nil value
	msg.FatalOnErr(p.Err)
}
