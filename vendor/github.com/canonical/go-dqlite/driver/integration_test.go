package driver_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	dqlite "github.com/canonical/go-dqlite"
	"github.com/canonical/go-dqlite/client"
	"github.com/canonical/go-dqlite/driver"
	"github.com/canonical/go-dqlite/internal/logging"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_DatabaseSQL(t *testing.T) {
	db, _, cleanup := newDB(t)
	defer cleanup()

	tx, err := db.Begin()
	require.NoError(t, err)

	_, err = tx.Exec(`
CREATE TABLE test  (n INT, s TEXT);
CREATE TABLE test2 (n INT, t DATETIME DEFAULT CURRENT_TIMESTAMP)
`)
	require.NoError(t, err)

	stmt, err := tx.Prepare("INSERT INTO test(n, s) VALUES(?, ?)")
	require.NoError(t, err)

	_, err = stmt.Exec(int64(123), "hello")
	require.NoError(t, err)

	require.NoError(t, stmt.Close())

	_, err = tx.Exec("INSERT INTO test2(n) VALUES(?)", int64(456))
	require.NoError(t, err)

	require.NoError(t, tx.Commit())

	tx, err = db.Begin()
	require.NoError(t, err)

	rows, err := tx.Query("SELECT n, s FROM test")
	require.NoError(t, err)

	for rows.Next() {
		var n int64
		var s string

		require.NoError(t, rows.Scan(&n, &s))

		assert.Equal(t, int64(123), n)
		assert.Equal(t, "hello", s)
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())

	rows, err = tx.Query("SELECT n, t FROM test2")
	require.NoError(t, err)

	for rows.Next() {
		var n int64
		var s time.Time

		require.NoError(t, rows.Scan(&n, &s))

		assert.Equal(t, int64(456), n)
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())

	require.NoError(t, tx.Rollback())

	require.NoError(t, db.Close())
}

func TestIntegration_Error(t *testing.T) {
	db, _, cleanup := newDB(t)
	defer cleanup()

	_, err := db.Exec("CREATE TABLE test (n INT, UNIQUE (n))")
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO test (n) VALUES (1)")
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO test (n) VALUES (1)")
	if err, ok := err.(driver.Error); ok {
		assert.Equal(t, int(sqlite3.ErrConstraintUnique), err.Code)
		assert.Equal(t, "UNIQUE constraint failed: test.n", err.Message)
	} else {
		t.Fatalf("expected diver error, got %+v", err)
	}
}

func TestIntegration_ConfigMultiThread(t *testing.T) {
	_, _, cleanup := newDB(t)
	defer cleanup()

	err := dqlite.ConfigMultiThread()
	assert.EqualError(t, err, "SQLite is already initialized")
}

func TestIntegration_LargeQuery(t *testing.T) {
	db, _, cleanup := newDB(t)
	defer cleanup()

	tx, err := db.Begin()
	require.NoError(t, err)

	_, err = tx.Exec("CREATE TABLE test (n INT)")
	require.NoError(t, err)

	stmt, err := tx.Prepare("INSERT INTO test(n) VALUES(?)")
	require.NoError(t, err)

	for i := 0; i < 512; i++ {
		_, err = stmt.Exec(int64(i))
		require.NoError(t, err)
	}

	require.NoError(t, stmt.Close())

	require.NoError(t, tx.Commit())

	tx, err = db.Begin()
	require.NoError(t, err)

	rows, err := tx.Query("SELECT n FROM test")
	require.NoError(t, err)

	columns, err := rows.Columns()
	require.NoError(t, err)

	assert.Equal(t, []string{"n"}, columns)

	count := 0
	for i := 0; rows.Next(); i++ {
		var n int64

		require.NoError(t, rows.Scan(&n))

		assert.Equal(t, int64(i), n)
		count++
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())

	assert.Equal(t, count, 512)

	require.NoError(t, tx.Rollback())

	require.NoError(t, db.Close())
}

// Build a 2-node cluster, kill one node and recover the other.
func TestIntegration_Recover(t *testing.T) {
	n := 2
	infos := make([]client.NodeInfo, n)
	for i := range infos {
		infos[i].ID = uint64(i + 1)
		infos[i].Address = fmt.Sprintf("@%d", infos[i].ID)
	}
	dirs := make([]string, n)
	nodes := make([]*dqlite.Node, 2)
	for i, info := range infos {
		dir, cleanup := newDir(t)
		defer cleanup()
		node, err := dqlite.New(info.ID, info.Address, dir, dqlite.WithBindAddress(info.Address))
		require.NoError(t, err)
		require.NoError(t, node.Start())
		nodes[i] = node
		dirs[i] = dir
	}

	store, err := client.DefaultNodeStore(":memory:")
	require.NoError(t, err)
	require.NoError(t, store.Set(context.Background(), infos))

	log := logging.Test(t)
	driver, err := driver.New(store, driver.WithLogFunc(log))
	require.NoError(t, err)

	driverName := registerDriver(driver)
	db, err := sql.Open(driverName, "test.db")
	require.NoError(t, err)
	defer db.Close()

	client, err := client.New(context.Background(), nodes[0].BindAddress())
	require.NoError(t, err)
	defer client.Close()

	require.NoError(t, client.Add(context.Background(), infos[1]))

	_, err = db.Exec("CREATE TABLE test (n INT)")
	require.NoError(t, err)

	nodes[0].Close()
	nodes[1].Close()

	node, err := dqlite.New(1, "@1", dirs[0], dqlite.WithBindAddress("@1"))
	require.NoError(t, err)

	require.NoError(t, node.Recover(infos[0:1]))

	require.NoError(t, node.Start())
	defer node.Close()

	// FIXME: this is necessary otherwise the INSERT below fails with "no
	// such table", because the replication hooks are not triggered and the
	// barrier is not applied.
	_, err = db.Exec("CREATE TABLE test2 (n INT)")
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO test(n) VALUES(1)")
	require.NoError(t, err)
}

func newDB(t *testing.T) (*sql.DB, []*dqlite.Node, func()) {
	n := 3

	infos := make([]client.NodeInfo, n)
	for i := range infos {
		infos[i].ID = uint64(i + 1)
		infos[i].Address = fmt.Sprintf("@%d", infos[i].ID)
	}

	servers, cleanup := newNodes(t, infos)

	store, err := client.DefaultNodeStore(":memory:")
	require.NoError(t, err)

	require.NoError(t, store.Set(context.Background(), infos))

	log := logging.Test(t)
	driver, err := driver.New(store, driver.WithLogFunc(log))
	require.NoError(t, err)

	driverName := fmt.Sprintf("dqlite-integration-test-%d", driversCount)
	sql.Register(driverName, driver)

	driversCount++

	db, err := sql.Open(driverName, "test.db")
	require.NoError(t, err)

	return db, servers, cleanup
}

func registerDriver(driver *driver.Driver) string {
	name := fmt.Sprintf("dqlite-integration-test-%d", driversCount)
	sql.Register(name, driver)
	driversCount++
	return name
}

func newNodes(t *testing.T, infos []client.NodeInfo) ([]*dqlite.Node, func()) {
	t.Helper()

	n := len(infos)
	servers := make([]*dqlite.Node, n)
	cleanups := make([]func(), 0)

	for i, info := range infos {
		dir, dirCleanup := newDir(t)
		server, err := dqlite.New(info.ID, info.Address, dir, dqlite.WithBindAddress(info.Address))
		require.NoError(t, err)

		cleanups = append(cleanups, func() {
			require.NoError(t, server.Close())
			dirCleanup()
		})

		err = server.Start()
		require.NoError(t, err)

		servers[i] = server

	}

	cleanup := func() {
		for _, f := range cleanups {
			f()
		}
	}

	return servers, cleanup
}

var driversCount = 0

func TestIntegration_ColumnTypeName(t *testing.T) {
	db, _, cleanup := newDB(t)
	defer cleanup()

	_, err := db.Exec("CREATE TABLE test (n INT, UNIQUE (n))")
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO test (n) VALUES (1)")
	require.NoError(t, err)

	rows, err := db.Query("SELECT n FROM test")
	require.NoError(t, err)
	defer rows.Close()

	types, err := rows.ColumnTypes()
	require.NoError(t, err)

	assert.Equal(t, "INTEGER", types[0].DatabaseTypeName())

	require.True(t, rows.Next())
	var n int64
	err = rows.Scan(&n)
	require.NoError(t, err)

	assert.Equal(t, int64(1), n)
}

