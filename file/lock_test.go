package file_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/andrejacobs/go-aj/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcquireLockfile(t *testing.T) {
	lockPath := filepath.Join(os.TempDir(), "unit-test.lock")
	os.Remove(lockPath)
	defer os.Remove(lockPath)

	// Lock
	lock, err := file.AcquireLockfile(lockPath)
	require.NoError(t, err)
	require.NotNil(t, lock)
	assert.Equal(t, lockPath, lock.Path())
	assert.Equal(t, os.Getpid(), lock.Pid())

	// Can't lock (even though it is the same PID)
	fail, err := file.AcquireLockfile(lockPath)
	require.NotNil(t, fail)
	assert.ErrorIs(t, err, file.ErrLockfileAcquired)
	assert.Equal(t, lockPath, fail.Path())
	assert.Equal(t, os.Getpid(), fail.Pid())

	// Can release what you own as many times as you want
	for i := 0; i < 5; i++ {
		err = lock.Release()
		assert.NoError(t, err)
	}

	// Lock
	lock, err = file.AcquireLockfile(lockPath)
	defer lock.Release()
	require.NoError(t, err)
	require.NotNil(t, lock)
	assert.Equal(t, lockPath, lock.Path())
	assert.Equal(t, os.Getpid(), lock.Pid())
}

func TestAcquireLockfileReEntrant(t *testing.T) {
	lockPath := filepath.Join(os.TempDir(), "unit-test.lock")
	os.Remove(lockPath)
	defer os.Remove(lockPath)

	lock, err := file.AcquireLockfile(lockPath)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		// Lock again
		lock, err := file.AcquireLockfileReEntrant(lockPath)
		require.NoError(t, err)
		require.NotNil(t, lock)
		assert.Equal(t, lockPath, lock.Path())
		assert.Equal(t, os.Getpid(), lock.Pid())
	}

	require.NoError(t, lock.Release())
}

func TestReleaseNotOwnedLockfile(t *testing.T) {
	lockPath := filepath.Join(os.TempDir(), "unit-test.lock")
	os.Remove(lockPath)
	defer os.Remove(lockPath)

	// Simulate that some other process owns the lock file
	f, err := os.Create(lockPath)
	require.NoError(t, err)
	_, err = f.WriteString(strconv.Itoa(os.Getpid() + 100))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	fail, err := file.AcquireLockfile(lockPath)
	require.NotNil(t, fail)
	assert.ErrorIs(t, err, file.ErrLockfileAcquired)

	// Can't release what you don't own
	err = fail.Release()
	assert.ErrorIs(t, err, file.ErrLockfileNotOwned)
}

func TestInvalidLockfile(t *testing.T) {
	lockPath := filepath.Join(os.TempDir(), "unit-test.lock")
	os.Remove(lockPath)
	defer os.Remove(lockPath)

	f, err := os.Create(lockPath)
	require.NoError(t, err)
	_, err = f.WriteString("lol-nan")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	fail, err := file.AcquireLockfile(lockPath)
	require.NotNil(t, fail)
	assert.ErrorIs(t, err, file.ErrLockfileAcquired)
	var numErr *strconv.NumError
	assert.ErrorAs(t, err, &numErr)
}
