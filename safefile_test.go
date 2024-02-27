// safefile_test.go -- tests for safefile impl

package utils

import (
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"flag"
	"fmt"
	mrand "math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/opencoff/go-mmap"
)

var testDir = flag.String("testdir", "", "Use 'T' as the testdir for COW tests")

func TestSimpleFile(t *testing.T) {
	assert := newAsserter(t)
	tmpdir := t.TempDir()

	fn := filepath.Join(tmpdir, "file-1")

	ck1, err := createFile(fn, 0)
	assert(err == nil, "can't create tmpfile: %s", err)

	sf, err := NewSafeFile(fn, 0, 0, 0600)
	assert(err != nil, "%s: bypassed overwrite protection", fn)

	buf := make([]byte, 128*1024)
	randbuf(buf)

	sf, err = NewSafeFile(fn, OPT_OVERWRITE, 0, 0600)
	assert(err == nil, "%s: can't create safefile: %s", fn, err)
	assert(sf != nil, "%s: nil ptr", fn)

	n, err := sf.Write(buf)
	assert(err == nil, "%s: write error: %s", sf.Name(), err)
	assert(n == len(buf), "%s: partial write: exp %d, saw %d", sf.Name(), len(buf), n)

	err = sf.Close()
	assert(err == nil, "%s: close: %s", sf.Name(), err)

	ck2 := sha256.Sum256(buf)
	assert(0 == subtle.ConstantTimeCompare(ck1, ck2[:]), "%s: file checksums match after change!", fn)

	ck3, err := fileCksum(fn)
	assert(err == nil, "%s: cksum error: %s", fn, err)
	assert(1 == subtle.ConstantTimeCompare(ck2[:], ck3), "%s: file checksums mismatch", fn)
}

func TestAbort(t *testing.T) {
	assert := newAsserter(t)
	tmpdir := t.TempDir()

	fn := filepath.Join(tmpdir, "file-1")

	ck1, err := createFile(fn, 0)
	assert(err == nil, "can't create tmpfile: %s", err)

	buf := make([]byte, 128*1024)
	randbuf(buf)

	sf, err := NewSafeFile(fn, OPT_OVERWRITE, 0, 0600)
	assert(err == nil, "%s: can't create safefile: %s", fn, err)
	assert(sf != nil, "%s: nil ptr", fn)

	n, err := sf.Write(buf)
	assert(err == nil, "%s: write error: %s", sf.Name(), err)
	assert(n == len(buf), "%s: partial write: exp %d, saw %d", sf.Name(), len(buf), n)

	sf.Abort()

	// File original contents shouldn't change
	ck3, err := fileCksum(fn)
	assert(err == nil, "%s: cksum error: %s", fn, err)
	assert(1 == subtle.ConstantTimeCompare(ck1, ck3), "%s: file checksums mismatch", fn)
}

func TestCow(t *testing.T) {
	assert := newAsserter(t)

	tmpdir := t.TempDir()
	if len(*testDir) > 0 {
		tmpdir = filepath.Join(*testDir, t.Name())
		err := os.MkdirAll(tmpdir, 0750)
		assert(err == nil, "can't make COW tmpdir %s: %s", tmpdir, err)
		t.Logf("%s: Using %s as the COW tempdir\n", t.Name(), tmpdir)
		t.Cleanup(func() {
			os.RemoveAll(tmpdir)
		})
	}

	fn := filepath.Join(tmpdir, "file-1")

	const (
		_ChunkSize      int = 8192
		_ModifiedChunks     = 8
		_MaxChunks          = 64
		_FileSize           = _MaxChunks * _ChunkSize
	)

	ck1, err := createFile2(fn, _MaxChunks, _ChunkSize)
	assert(err == nil, "can't create tmpfile: %s", err)

	buf := make([]byte, _ModifiedChunks*_ChunkSize)
	randbuf(buf)

	// calculate checksums of new chunks
	ck2 := make([][]byte, 0, len(ck1))
	for i := 0; i < _ModifiedChunks; i++ {
		off := i * _ChunkSize
		end := off + _ChunkSize
		sum := sha256.Sum256(buf[off:end])
		ck2 = append(ck2, sum[:])
	}

	sf, err := NewSafeFile(fn, OPT_OVERWRITE|OPT_COW, 0, 0600)
	assert(err == nil, "%s: can't create safefile: %s", fn, err)
	assert(sf != nil, "%s: nil ptr", fn)

	n, err := sf.Write(buf)
	assert(err == nil, "%s: write error: %s", sf.Name(), err)
	assert(n == len(buf), "%s: partial write: exp %d, saw %d", sf.Name(), len(buf), n)

	err = sf.Close()
	assert(err == nil, "%s: close: %s", sf.Name(), err)

	// Only 8 chunks of total have changed. So the rest should be fine.
	fd, err := os.Open(fn)
	assert(err == nil, "%s: open %s", fn, err)

	st, err := fd.Stat()
	assert(err == nil, "%s: stat %s", fn, err)
	assert(st.Size() == int64(_FileSize), "%s: file size: exp %d, saw %d", fn, _FileSize, st.Size())

	buf = buf[:_ChunkSize]
	ck3 := make([][]byte, 0, _MaxChunks)
	for i := 0; i < _MaxChunks; i++ {
		n, err := fd.Read(buf)
		assert(err == nil, "%s: read %s", fn, err)
		assert(n == len(buf), "%s: partial read", fn)
		sum := sha256.Sum256(buf)
		ck3 = append(ck3, sum[:])
	}
	fd.Close()

	// The file must now have combination of old and new checksums

	// the first N chunks are the new ones
	for i := 0; i < _ModifiedChunks; i++ {
		x := ck1[i]
		a := ck2[i]
		b := ck3[i]
		assert(0 == subtle.ConstantTimeCompare(x, b), "%s: Chunk %d retained orig contents", fn, i)
		assert(1 == subtle.ConstantTimeCompare(a, b), "%s: Chunk %d mismatch", fn, i)
	}

	// and the rest must be same as the original file
	for i := _ModifiedChunks; i < _MaxChunks; i++ {
		a := ck1[i]
		b := ck3[i]
		assert(1 == subtle.ConstantTimeCompare(a, b), "%s: Chunk %d mismatch", fn, i)
	}
}

func fileCksum(nm string) ([]byte, error) {
	fd, err := os.Open(nm)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	h := sha256.New()
	_, err = mmap.Reader(fd, func(b []byte) error {
		h.Write(b)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return h.Sum(nil)[:], nil
}

// create a file and return cryptographic checksum
func createFile(nm string, sizehint int64) ([]byte, error) {
	fd, err := os.OpenFile(nm, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	if sizehint <= 0 {
		for sizehint < 1048576 {
			sizehint = mrand.Int63n(20 * 1048576)
		}
	}

	buf := make([]byte, 65536)
	h := sha256.New()

	// fill it with random data
	for sizehint > 0 {
		randbuf(buf)
		h.Write(buf)
		n, err := fd.Write(buf)
		if err != nil {
			return nil, err
		}
		if n != len(buf) {
			return nil, fmt.Errorf("%s: partial write (exp %d, saw %d)", nm, len(buf), n)
		}
		sizehint -= int64(n)
	}

	if err = fd.Sync(); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func createFile2(nm string, chunks, chunksize int) ([][]byte, error) {
	fd, err := os.OpenFile(nm, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	buf := make([]byte, chunksize)
	ck := make([][]byte, 0, chunks)

	for i := 0; i < chunks; i++ {
		randbuf(buf)
		sum := sha256.Sum256(buf)
		ck = append(ck, sum[:])

		n, err := fd.Write(buf)
		if err != nil {
			return nil, err
		}
		if n != len(buf) {
			return nil, fmt.Errorf("%s: partial write (exp %d, saw %d)", nm, len(buf), n)
		}
	}
	return ck, nil
}

func randbuf(b []byte) []byte {
	n, err := crand.Read(b)
	if err != nil || n != len(b) {
		panic(fmt.Sprintf("can't read %d bytes of crypto/rand: %s", len(b), err))
	}
	return b
}
