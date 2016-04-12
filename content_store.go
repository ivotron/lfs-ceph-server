package main

import (
	"bytes"
	"errors"
	"io"

	"github.com/ceph/go-ceph/rados"
)

var (
	errHashMismatch = errors.New("Content hash does not match OID")
	errSizeMismatch = errors.New("Content size does not match")
)

// ContentStore provides a simple file system based storage.
type ContentStore struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
}

// NewContentStore creates a ContentStore at the base directory.
func NewContentStore(base string) (*ContentStore, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}
	conn.ReadConfigFile("/etc/ceph/ceph.conf")
	conn.Connect()
	ioctx, err := conn.OpenIOContext("data")
	if err != nil {
		return nil, err
	}

	return &ContentStore{conn, ioctx}, nil
}

// Get takes a Meta object and retreives the content from the store, returning
// it as an io.Reader.
func (s *ContentStore) Get(meta *MetaObject) (io.Reader, error) {
	stat, err := s.ioctx.Stat(meta.Oid)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, stat.Size)
	bytesRead, err := s.ioctx.Read(meta.Oid, buf, 0)
	if err != nil {
		return nil, err
	}
	if uint64(bytesRead) != stat.Size {
		return nil, errSizeMismatch
	}
	return bytes.NewReader(buf), nil
}

// Put takes a Meta object and an io.Reader and writes the content to the store.
func (s *ContentStore) Put(meta *MetaObject, r io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	buf.Bytes()
	if err := s.ioctx.Write(meta.Oid, buf.Bytes(), 0); err != nil {
		return err
	}
	return nil
}

// Exists returns true if the object exists in the content store.
func (s *ContentStore) Exists(meta *MetaObject) bool {
	if _, err := s.ioctx.Stat(meta.Oid); err != nil {
		if err == rados.RadosErrorNotFound {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}
