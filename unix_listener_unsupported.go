// +build !linux,!freebsd

package dkvolume

import (
	"errors"
	"net"
)

var (
	errOnlySupportedOnLinuxAndFreeBSD = errors.New("unix socket creation is only supported on linux and freebsd")
)

func newUnixListener(
	volumeDriverName string,
	group string,
) (net.Listener, string, error) {
	return nil, "", errOnlySupportedOnLinuxAndFreeBSD
}
