package mocks

import (
	"fmt"
	"net"
	"time"

	"github.com/puppetlabs/lumogon/utils"
)

// MockNetConn TODO
type MockNetConn struct {
	ReadFn             func(b []byte) (n int, err error)
	WriteFn            func(b []byte) (n int, err error)
	CloseFn            func() error
	LocalAddrFn        func() net.Addr
	RemoteAddrFn       func() net.Addr
	SetDeadlineFn      func(t time.Time) error
	SetReadDeadlineFn  func(t time.Time) error
	SetWriteDeadlineFn func(t time.Time) error
}

// Read TODO
func (c MockNetConn) Read(b []byte) (n int, err error) {
	if c.ReadFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		fmt.Println("[MockNetConn]  - len(b): ", len(b))
		return c.ReadFn(b)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// Write TODO
func (c MockNetConn) Write(b []byte) (n int, err error) {
	if c.WriteFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		fmt.Println("[MockNetConn]  - b: ", b)
		return c.WriteFn(b)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// Close TODO
func (c MockNetConn) Close() error {
	if c.CloseFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		return c.CloseFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// LocalAddr TODO
func (c MockNetConn) LocalAddr() net.Addr {
	if c.LocalAddrFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		return c.LocalAddrFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// RemoteAddr TODO
func (c MockNetConn) RemoteAddr() net.Addr {
	if c.RemoteAddrFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		return c.RemoteAddrFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// SetDeadline TODO
func (c MockNetConn) SetDeadline(t time.Time) error {
	if c.SetDeadlineFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		fmt.Println("[MockNetConn]  - t: ", t)
		return c.SetDeadlineFn(t)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// SetReadDeadline TODO
func (c MockNetConn) SetReadDeadline(t time.Time) error {
	if c.SetReadDeadlineFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		fmt.Println("[MockNetConn]  - t: ", t)
		return c.SetReadDeadlineFn(t)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// SetWriteDeadline TODO
func (c MockNetConn) SetWriteDeadline(t time.Time) error {
	if c.SetWriteDeadlineFn != nil {
		fmt.Println("[MockNetConn] In ", utils.CurrentFunctionName())
		fmt.Println("[MockNetConn]  - t: ", t)
		return c.SetWriteDeadlineFn(t)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}
