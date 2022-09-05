package core

import (
	"golang.org/x/sys/unix"
	"syscall"
)

type sockOptions struct {
	// SOL_SOCKET opts
	keepAlive int
	rcvBuf    int
	sndBuf    int
	rcvLowAt  int
	sndLowAt  int
	ReuseAddr bool
	ReusePort bool
	// IPPROTO_TCP opts
	NoDelay bool
}

func defaultSockOptions() *sockOptions {
	return &sockOptions{
		keepAlive: 0,
		rcvBuf:    0,
		sndBuf:    0,
		rcvLowAt:  0,
		sndLowAt:  0,
		ReuseAddr: false,
		ReusePort: false,
		NoDelay:   false,
	}
}

func (opt *sockOptions) Control() func(fd uintptr) {
	bool2int := func(value bool) int {
		if value {
			return 1
		} else {
			return 0
		}
	}

	return func(fd uintptr) {
		_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEADDR, bool2int(opt.ReuseAddr))
		_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, bool2int(opt.ReusePort))
	}
}

// SockOption configures how we set up the socket.
type SockOption interface {
	apply(options *sockOptions)
}

// funcDialOption wraps a function that modifies sockOptions into an
// implementation of the SockOption interface.
type funcSockOption struct {
	f func(options *sockOptions)
}

func (fdo *funcSockOption) apply(do *sockOptions) {
	fdo.f(do)
}

func WithReuseAddr(reuseAddr bool) SockOption {
	return &funcSockOption{f: func(options *sockOptions) {
		options.ReuseAddr = reuseAddr
	}}
}

func WithReusePort(reusePort bool) SockOption {
	return &funcSockOption{f: func(options *sockOptions) {
		options.ReusePort = reusePort
	}}
}
