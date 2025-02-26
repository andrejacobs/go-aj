package file_test

import (
	"context"
	"io/fs"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/andrejacobs/go-micropkg/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentWalker(t *testing.T) {
	// std lib walker
	expected := make([]string, 0, 10)
	err := filepath.Walk(tempDir, func(path string, info fs.FileInfo, err error) error {
		expected = append(expected, path)
		return nil
	})
	require.NoError(t, err)
	slices.Sort(expected)

	result := make([]string, 0, 10)
	walkFunc := func(path string, entry fs.DirEntry, fileInfo fs.FileInfo) error {
		result = append(result, path)
		return nil
	}

	walker := file.NewConcurrentWalker()

	// Test for re-entrancy issues
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()

		_, cancel, err := walker.StartWalking(ctx, tempDir, walkFunc)
		defer cancel()
		require.NoError(t, err)
		walker.Wait()
		require.False(t, walker.HadErrors())
		slices.Sort(result)

		assert.ElementsMatch(t, expected, result)
		result = result[:0]
	}
}

func TestConcurrentWalkerCancel(t *testing.T) {
	walker := file.NewConcurrentWalker()

	walkFunc := func(path string, entry fs.DirEntry, fileInfo fs.FileInfo) error {
		// Block a bit
		time.Sleep(time.Millisecond * 100)
		return nil
	}

	_, cancel, err := walker.StartWalking(context.TODO(), tempDir, walkFunc)
	defer cancel()
	require.NoError(t, err)

	// Cancel after a few milliseconds
	time.AfterFunc(time.Millisecond*100, func() {
		cancel()
	})

	walker.Wait()
}

func TestFileExcluder(t *testing.T) {
	excluder := fileExcluder{
		filenames: []string{"a", "b"},
	}
	walker := file.NewConcurrentWalker().SetFileExcluder(&excluder)

	walkFunc := func(path string, entry fs.DirEntry, fileInfo fs.FileInfo) error {
		filename := filepath.Base(path)
		if filename == "a" || filename == "b" {
			assert.Fail(t, "expected the excluder to be used")
		}
		return nil
	}

	_, _, err := walker.StartWalking(context.TODO(), tempDir, walkFunc)
	require.NoError(t, err)
	walker.Wait()
}

func TestDirExcluder(t *testing.T) {
	excluder := fileExcluder{
		filenames: []string{"d"},
	}
	walker := file.NewConcurrentWalker().SetDirExcluder(&excluder)

	walkFunc := func(path string, entry fs.DirEntry, fileInfo fs.FileInfo) error {
		filename := filepath.Base(path)
		if filename == "d" {
			assert.Fail(t, "expected the excluder to be used")
		}
		return nil
	}

	_, _, err := walker.StartWalking(context.TODO(), tempDir, walkFunc)
	require.NoError(t, err)
	walker.Wait()
}

func BenchmarkConcurrentWalker(b *testing.B) {
	// For small directories the concurrent walker is slower
	// point this to somewhere with a lot of files etc.
	//tempDir := "/Users/andre/temp"

	b.Run("filepath.Walk", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := filepath.Walk(tempDir, func(path string, info fs.FileInfo, err error) error {
				return nil
			})
			if err != nil {
				b.Error(err)
			}
		}
	})

	walkFunc := func(path string, entry fs.DirEntry, fileInfo fs.FileInfo) error {
		return nil
	}

	walker := file.NewConcurrentWalker()

	ctx := context.TODO()
	b.Run("ConcurrentWalker", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			_, cancel, err := walker.StartWalking(ctx, tempDir, walkFunc)
			if err != nil {
				b.Error(err)
				continue
			}
			defer cancel()
			walker.Wait()

			if walker.HadErrors() {
				b.Errorf("had errors")
			}
		}
	})
}

//-----------------------------------------------------------------------------

type fileExcluder struct {
	filenames []string
}

func (e *fileExcluder) Match(path string) (bool, error) {
	filename := filepath.Base(path)
	for _, match := range e.filenames {
		if filename == match {
			return true, nil
		}
	}
	return false, nil
}
