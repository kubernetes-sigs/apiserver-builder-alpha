package protocol_test

import (
	"context"
	"testing"
	"time"

	"github.com/canonical/go-dqlite/internal/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestProtocol_Heartbeat(t *testing.T) {
// 	c, cleanup := newProtocol(t)
// 	defer cleanup()

// 	request, response := newMessagePair(512, 512)

// 	protocol.EncodeHeartbeat(&request, uint64(time.Now().Unix()))

// 	makeCall(t, c, &request, &response)

// 	servers, err := protocol.DecodeNodes(&response)
// 	require.NoError(t, err)

// 	assert.Len(t, servers, 2)
// 	assert.Equal(t, client.Nodes{
// 		{ID: uint64(1), Address: "1.2.3.4:666"},
// 		{ID: uint64(2), Address: "5.6.7.8:666"}},
// 		servers)
// }

// Test sending a request that needs to be written into the dynamic buffer.
func TestProtocol_RequestWithDynamicBuffer(t *testing.T) {
	p, cleanup := newProtocol(t)
	defer cleanup()

	request, response := newMessagePair(64, 64)

	protocol.EncodeOpen(&request, "test.db", 0, "test-0")

	makeCall(t, p, &request, &response)

	id, err := protocol.DecodeDb(&response)
	require.NoError(t, err)

	request.Reset()
	response.Reset()

	sql := `
CREATE TABLE foo (n INT);
CREATE TABLE bar (n INT);
CREATE TABLE egg (n INT);
CREATE TABLE baz (n INT);
`
	protocol.EncodeExecSQL(&request, uint64(id), sql, nil)

	makeCall(t, p, &request, &response)
}

func TestProtocol_Prepare(t *testing.T) {
	c, cleanup := newProtocol(t)
	defer cleanup()

	request, response := newMessagePair(64, 64)

	protocol.EncodeOpen(&request, "test.db", 0, "test-0")

	makeCall(t, c, &request, &response)

	db, err := protocol.DecodeDb(&response)
	require.NoError(t, err)

	request.Reset()
	response.Reset()

	protocol.EncodePrepare(&request, uint64(db), "CREATE TABLE test (n INT)")

	makeCall(t, c, &request, &response)

	_, stmt, params, err := protocol.DecodeStmt(&response)
	require.NoError(t, err)

	assert.Equal(t, uint32(0), stmt)
	assert.Equal(t, uint64(0), params)
}

/*
func TestProtocol_Exec(t *testing.T) {
	client, cleanup := newProtocol(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	db, err := client.Open(ctx, "test.db", "volatile")
	require.NoError(t, err)

	stmt, err := client.Prepare(ctx, db.ID, "CREATE TABLE test (n INT)")
	require.NoError(t, err)

	_, err = client.Exec(ctx, db.ID, stmt.ID)
	require.NoError(t, err)
}

func TestProtocol_Query(t *testing.T) {
	client, cleanup := newProtocol(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	db, err := client.Open(ctx, "test.db", "volatile")
	require.NoError(t, err)

	start := time.Now()

	stmt, err := client.Prepare(ctx, db.ID, "CREATE TABLE test (n INT)")
	require.NoError(t, err)

	_, err = client.Exec(ctx, db.ID, stmt.ID)
	require.NoError(t, err)

	_, err = client.Finalize(ctx, db.ID, stmt.ID)
	require.NoError(t, err)

	stmt, err = client.Prepare(ctx, db.ID, "INSERT INTO test VALUES(1)")
	require.NoError(t, err)

	_, err = client.Exec(ctx, db.ID, stmt.ID)
	require.NoError(t, err)

	_, err = client.Finalize(ctx, db.ID, stmt.ID)
	require.NoError(t, err)

	stmt, err = client.Prepare(ctx, db.ID, "SELECT n FROM test")
	require.NoError(t, err)

	_, err = client.Query(ctx, db.ID, stmt.ID)
	require.NoError(t, err)

	_, err = client.Finalize(ctx, db.ID, stmt.ID)
	require.NoError(t, err)

	fmt.Printf("time %s\n", time.Since(start))
}
*/

func newProtocol(t *testing.T) (*protocol.Protocol, func()) {
	t.Helper()

	address, serverCleanup := newNode(t, 0)

	store := newStore(t, []string{address})

	connector := newConnector(t, store)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	client, err := connector.Connect(ctx)

	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		serverCleanup()
	}

	return client, cleanup
}

// Perform a client call.
func makeCall(t *testing.T, p *protocol.Protocol, request, response *protocol.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err := p.Call(ctx, request, response)
	require.NoError(t, err)
}

// Return a new message pair to be used as request and response.
func newMessagePair(size1, size2 int) (protocol.Message, protocol.Message) {
	message1 := protocol.Message{}
	message1.Init(size1)

	message2 := protocol.Message{}
	message2.Init(size2)

	return message1, message2
}
