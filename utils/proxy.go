package utils

import (
	"io"
	"net"
	"sync"
	"zhtcloud/utils/logger"
)

func StartTCPListener(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Error starting TCP listener on %s: %v", addr, err)
		return
	}
	defer listener.Close()

	var connA, connB net.Conn

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Error accepting connection: %v", err)
			continue
		}

		if connA == nil {
			connA = conn
			continue
		} else if connB == nil {
			connB = conn

			// local variables to avoid closure issues
			localConnA := connA
			localConnB := connB

			// ensure close connections only once
			var once sync.Once
			closeConns := func() {
				if localConnA != nil {
					localConnA.Close()
					connA = nil
				}
				if localConnB != nil {
					localConnB.Close()
					connB = nil
				}
			}

			go func() {
				defer once.Do(closeConns)
				io.Copy(localConnB, localConnA)
			}()

			go func() {
				defer once.Do(closeConns)
				io.Copy(localConnA, localConnB)
			}()

			continue
		} else {
			logger.Warning("Too many connections, closing new one")
			conn.Close()
		}
	}
}
