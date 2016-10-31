package iface

import (
	"github.com/Supernomad/nexus/nexusd/common"
	"syscall"
)

// Socket is a generic multique socket
type Socket struct {
	queues []int
	log    *common.Logger
	cfg    *common.Config
}

// Open the socket
func (sock *Socket) Open() error {
	for i := 0; i < sock.cfg.NumWorkers; i++ {
		var queue int
		var err error

		if !sock.cfg.ReuseFDS {
			queue, err = createSocket()
			if err != nil {
				return err
			}
			err = initSocket(queue, sock.cfg.ListenAddress)
			if err != nil {
				return err
			}
		} else {
			queue = 3 + sock.cfg.NumWorkers + i
		}
		sock.queues[i] = queue
	}
	return nil
}

// Close the socket
func (sock *Socket) Close() error {
	for i := 0; i < len(sock.queues); i++ {
		if err := syscall.Close(sock.queues[i]); err != nil {
			return err
		}
	}
	return nil
}

// GetFDs will return the underlying queue fds
func (sock *Socket) GetFDs() []int {
	return sock.queues
}

// Read a packet from the socket
func (sock *Socket) Read(queue int) (*common.Packet, error) {
	buf := make([]byte, 65536)
	n, src, err := syscall.Recvfrom(sock.queues[queue], buf, 0)
	if err != nil {
		return nil, err
	}
	return common.NewPacket(buf, src, n), nil
}

// Write a packet to the socket
func (sock *Socket) Write(queue int, packet *common.Packet) error {
	return syscall.Sendto(sock.queues[queue], packet.Data[:packet.Length], 0, packet.DestinationAddress)
}

func newSocket(log *common.Logger, cfg *common.Config) *Socket {
	queues := make([]int, cfg.NumWorkers)
	return &Socket{queues: queues, log: log, cfg: cfg}
}

func initSocket(queue int, sa syscall.Sockaddr) error {
	err := syscall.SetsockoptInt(queue, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return err
	}

	err = syscall.Bind(queue, sa)
	if err != nil {
		return err
	}

	return nil
}

func createSocket() (int, error) {
	return syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
}
