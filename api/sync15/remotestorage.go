package sync15

import "io"

type RemoteStorage interface {
	GetRootIndex() (hash string, generation int64, err error)
	GetReader(hash, name string) (io.ReadCloser, error)
}

type RemoteStorageWriter interface {
	UpdateRootIndex(hash string, generation int64) (gen int64, err error)
	GetWriter(hash, name string, writer io.WriteCloser) error
}
