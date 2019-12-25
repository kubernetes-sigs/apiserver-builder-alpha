package bindings_test

import (
	"encoding/binary"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/canonical/go-dqlite/internal/bindings"
	"github.com/canonical/go-dqlite/internal/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_Create(t *testing.T) {
	_, cleanup := newNode(t)
	defer cleanup()
}

func TestNode_Start(t *testing.T) {
	dir, cleanup := newDir(t)
	defer cleanup()

	server, err := bindings.NewNode(1, "1", dir)
	require.NoError(t, err)
	defer server.Close()

	err = server.SetBindAddress("@")
	require.NoError(t, err)

	err = server.Start()
	require.NoError(t, err)

	conn, err := net.Dial("unix", server.GetBindAddress())
	require.NoError(t, err)
	conn.Close()

	assert.True(t, strings.HasPrefix(server.GetBindAddress(), "@"))

	err = server.Stop()
	require.NoError(t, err)
}

func TestNode_Start_Inet(t *testing.T) {
	dir, cleanup := newDir(t)
	defer cleanup()

	server, err := bindings.NewNode(1, "1", dir)
	require.NoError(t, err)
	defer server.Close()

	err = server.SetBindAddress("127.0.0.1:9000")
	require.NoError(t, err)

	err = server.Start()
	require.NoError(t, err)

	conn, err := net.Dial("tcp", server.GetBindAddress())
	require.NoError(t, err)
	conn.Close()

	err = server.Stop()
	require.NoError(t, err)
}

func TestNode_Leader(t *testing.T) {
	_, cleanup := newNode(t)
	defer cleanup()

	conn := newClient(t)

	// Make a Leader request
	buf := makeClientRequest(t, conn, protocol.RequestLeader)
	assert.Equal(t, uint8(1), buf[0])

	require.NoError(t, conn.Close())
}

// func TestNode_Heartbeat(t *testing.T) {
// 	server, cleanup := newNode(t)
// 	defer cleanup()

// 	listener, cleanup := newListener(t)
// 	defer cleanup()

// 	cleanup = runNode(t, server, listener)
// 	defer cleanup()

// 	conn := newClient(t, listener)

// 	// Make a Heartbeat request
// 	makeClientRequest(t, conn, bindings.RequestHeartbeat)

// 	require.NoError(t, conn.Close())
// }

// func TestNode_ConcurrentHandleAndClose(t *testing.T) {
// 	server, cleanup := newNode(t)
// 	defer cleanup()

// 	listener, cleanup := newListener(t)
// 	defer cleanup()

// 	acceptCh := make(chan error)
// 	go func() {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			acceptCh <- err
// 		}
// 		server.Handle(conn)
// 		acceptCh <- nil
// 	}()

// 	conn, err := net.Dial("unix", listener.Addr().String())
// 	require.NoError(t, err)

// 	require.NoError(t, conn.Close())

// 	assert.NoError(t, <-acceptCh)
// }

// Create a new Node object for tests.
func newNode(t *testing.T) (*bindings.Node, func()) {
	t.Helper()

	dir, dirCleanup := newDir(t)

	server, err := bindings.NewNode(1, "1", dir)
	require.NoError(t, err)

	err = server.SetBindAddress("@test")
	require.NoError(t, err)

	require.NoError(t, server.Start())

	cleanup := func() {
		require.NoError(t, server.Stop())
		server.Close()
		dirCleanup()
	}

	return server, cleanup
}

// Create a new client network connection, performing the handshake.
func newClient(t *testing.T) net.Conn {
	t.Helper()

	conn, err := net.Dial("unix", "@test")
	require.NoError(t, err)

	// Handshake
	err = binary.Write(conn, binary.LittleEndian, protocol.VersionLegacy)
	require.NoError(t, err)

	return conn
}

// Perform a client request.
func makeClientRequest(t *testing.T, conn net.Conn, kind byte) []byte {
	t.Helper()

	// Number of words
	err := binary.Write(conn, binary.LittleEndian, uint32(1))
	require.NoError(t, err)

	// Type, flags, extra.
	n, err := conn.Write([]byte{kind, 0, 0, 0})
	require.NoError(t, err)
	require.Equal(t, 4, n)

	n, err = conn.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0}) // Unused single-word request payload
	require.NoError(t, err)
	require.Equal(t, 8, n)

	// Read the response
	conn.SetDeadline(time.Now().Add(250 * time.Millisecond))
	buf := make([]byte, 64)
	_, err = conn.Read(buf)
	require.NoError(t, err)

	return buf
}

// Return a new temporary directory.
func newDir(t *testing.T) (string, func()) {
	t.Helper()

	dir, err := ioutil.TempDir("", "dqlite-replication-test-")
	assert.NoError(t, err)

	cleanup := func() {
		_, err := os.Stat(dir)
		if err != nil {
			assert.True(t, os.IsNotExist(err))
		} else {
			assert.NoError(t, os.RemoveAll(dir))
		}
	}

	return dir, cleanup
}
