package client_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	dqlite "github.com/canonical/go-dqlite"
	"github.com/canonical/go-dqlite/client"
	"github.com/canonical/go-dqlite/internal/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Leader(t *testing.T) {
	node, cleanup := newNode(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := client.New(ctx, node.BindAddress())
	require.NoError(t, err)
	defer client.Close()

	leader, err := client.Leader(context.Background())
	require.NoError(t, err)

	assert.Equal(t, leader.ID, uint64(1))
	assert.Equal(t, leader.Address, "1")
}

func TestClient_Dump(t *testing.T) {
	node, cleanup := newNode(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := client.New(ctx, node.BindAddress())
	require.NoError(t, err)
	defer client.Close()

	// Open a database and create a test table.
	request := protocol.Message{}
	request.Init(4096)

	response := protocol.Message{}
	response.Init(4096)

	protocol.EncodeOpen(&request, "test.db", 0, "volatile")

	p := client.Protocol()
	err = p.Call(ctx, &request, &response)
	require.NoError(t, err)

	db, err := protocol.DecodeDb(&response)
	require.NoError(t, err)

	request.Reset()
	response.Reset()

	protocol.EncodeExecSQL(&request, uint64(db), "CREATE TABLE foo (n INT)", nil)

	err = p.Call(ctx, &request, &response)
	require.NoError(t, err)

	request.Reset()
	response.Reset()

	files, err := client.Dump(ctx, "test.db")
	require.NoError(t, err)

	require.Len(t, files, 2)
	assert.Equal(t, "test.db", files[0].Name)
	assert.Equal(t, 4096, len(files[0].Data))

	assert.Equal(t, "test.db-wal", files[1].Name)
	assert.Equal(t, 8272, len(files[1].Data))
}

func TestClient_Cluster(t *testing.T) {
	node, cleanup := newNode(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := client.New(ctx, node.BindAddress())
	require.NoError(t, err)
	defer client.Close()

	servers, err := client.Cluster(context.Background())
	require.NoError(t, err)

	assert.Len(t, servers, 1)
	assert.Equal(t, servers[0].ID, uint64(1))
	assert.Equal(t, servers[0].Address, "1")
}

func newNode(t *testing.T) (*dqlite.Node, func()) {
	t.Helper()
	dir, dirCleanup := newDir(t)

	node, err := dqlite.New(uint64(1), "1", dir, dqlite.WithBindAddress("@"))
	require.NoError(t, err)

	err = node.Start()
	require.NoError(t, err)

	cleanup := func() {
		require.NoError(t, node.Close())
		dirCleanup()
	}

	return node, cleanup
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
