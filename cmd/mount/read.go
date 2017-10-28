// +build linux darwin freebsd

package mount

import (
	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"github.com/ncw/rclone/fs"
	"github.com/ncw/rclone/vfs"
	"golang.org/x/net/context"
)

// ReadFileHandle is an open for read file handle on a File
type ReadFileHandle struct {
	*vfs.ReadFileHandle
}

// Check interface satisfied
var _ fusefs.Handle = (*ReadFileHandle)(nil)

// Check interface satisfied
var _ fusefs.HandleReader = (*ReadFileHandle)(nil)

// Read from the file handle
func (fh *ReadFileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) (err error) {
	var n int
	defer fs.Trace(fh, "len=%d, offset=%d", req.Size, req.Offset)("read=%d, err=%v", &n, &err)
	data := make([]byte, req.Size)
	n, err = fh.ReadFileHandle.ReadAt(data, req.Offset)
	if err != nil {
		return translateError(err)
	}
	resp.Data = data[:n]
	return nil
}

// Check interface satisfied
var _ fusefs.HandleFlusher = (*ReadFileHandle)(nil)

// Flush is called each time the file or directory is closed.
// Because there can be multiple file descriptors referring to a
// single opened file, Flush can be called multiple times.
func (fh *ReadFileHandle) Flush(ctx context.Context, req *fuse.FlushRequest) (err error) {
	defer fs.Trace(fh, "")("err=%v", &err)
	return translateError(fh.ReadFileHandle.Flush())
}

var _ fusefs.HandleReleaser = (*ReadFileHandle)(nil)

// Release is called when we are finished with the file handle
//
// It isn't called directly from userspace so the error is ignored by
// the kernel
func (fh *ReadFileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) (err error) {
	defer fs.Trace(fh, "")("err=%v", &err)
	return translateError(fh.ReadFileHandle.Release())
}
