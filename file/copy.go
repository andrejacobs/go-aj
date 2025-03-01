package file

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/andrejacobs/go-aj/file/contextio"
)

// Copy the source file to the destination and return the number of bytes that were copied
func CopyFile(ctx context.Context, source string, destination string) (int64, error) {
	src, dest, srcInfo, err := openFilesForCopying(source, destination)
	if err != nil {
		return 0, fmt.Errorf("failed to copy the file %q to %q. %v", source, destination, err)
	}
	defer src.Close()
	defer dest.Close()

	wc, err := copyN(ctx, src, dest, srcInfo.Size())
	if err != nil {
		return wc, fmt.Errorf("failed to copy the file %q to %q. %v", source, destination, err)
	}

	return wc, nil
}

// Copy N bytes from the source file to the destination and return the number of bytes that were copied
func CopyFileN(ctx context.Context, source string, destination string, count int64) (int64, error) {
	src, dest, _, err := openFilesForCopying(source, destination)
	if err != nil {
		return 0, fmt.Errorf("failed to copy the file %q to %q. %v", source, destination, err)
	}
	defer src.Close()
	defer dest.Close()

	wc, err := copyN(ctx, src, dest, count)
	if err != nil {
		return wc, fmt.Errorf("failed to copy the file %q to %q. %v", source, destination, err)
	}

	return wc, nil
}

func openFilesForCopying(source string, destination string) (*os.File, *os.File, fs.FileInfo, error) {
	src, err := os.Open(source)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open the source file %q. %v", source, err)
	}

	srcStat, err := src.Stat()
	if err != nil {
		src.Close()
		return nil, nil, nil, fmt.Errorf("failed to do Stat() on the source file %q. %v", source, err)
	}

	dest, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcStat.Mode().Perm())
	if err != nil {
		src.Close()
		return nil, nil, nil, fmt.Errorf("failed to create the destination file %q. %v", destination, err)
	}

	return src, dest, srcStat, nil
}

func copyN(ctx context.Context, src io.Reader, dest io.Writer, count int64) (int64, error) {
	in := contextio.NewReader(ctx, src)
	out := contextio.NewWriter(ctx, dest)

	wc, err := io.CopyN(out, in, count)
	return wc, err
}
