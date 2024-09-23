package sync15

import (
	"fmt"
	"io"
	"net/http"

	"github.com/juruen/rmapi/config"
	"github.com/juruen/rmapi/log"
	"github.com/juruen/rmapi/model"
	"github.com/juruen/rmapi/transport"
)

type BlobStorage struct {
	http        *transport.HttpClientCtx
	concurrency int
}

func NewBlobStorage(http *transport.HttpClientCtx) *BlobStorage {
	return &BlobStorage{
		http: http,
	}
}

const ROOT_NAME = "root"

func (b *BlobStorage) PutRootUrl(hash string, gen int64) (string, int64, error) {
	log.Trace.Println("fetching  ROOT url for: " + hash)
	req := model.BlobRootStorageRequest{
		Broadcast:  true, //TODO
		Hash:       hash,
		Generation: gen,
	}
	var res model.BlobStorageResponse

	if err := b.http.Post(transport.UserBearer, config.RootPut, req, &res); err != nil {
		return "", 0, err
	}
	return res.Url, res.MaxUploadSizeBytes, nil
}
func (b *BlobStorage) PutUrl(hash string) (string, int64, error) {
	log.Trace.Println("fetching PUT blob url for: " + hash)
	var req model.BlobStorageRequest
	var res model.BlobStorageResponse
	req.Method = http.MethodPut
	req.RelativePath = hash
	if err := b.http.Post(transport.UserBearer, config.UploadBlob, req, &res); err != nil {
		return "", 0, err
	}
	return res.Url, res.MaxUploadSizeBytes, nil
}

func (b *BlobStorage) GetUrl(hash string) (string, error) {
	log.Trace.Println("fetching GET blob url for: " + hash)
	var req model.BlobStorageRequest
	var res model.BlobStorageResponse
	req.Method = http.MethodGet
	req.RelativePath = hash
	if err := b.http.Post(transport.UserBearer, config.DownloadBlob, req, &res); err != nil {
		return "", err
	}
	return res.Url, nil
}

func (b *BlobStorage) GetReader(hash, filename string) (io.ReadCloser, error) {
	return b.http.GetStream(transport.UserBearer, config.BlobUrl+hash, filename)
}

func (b *BlobStorage) UploadBlob(hash, filename string, reader io.Reader) error {
	log.Trace.Println("uploading blob ", filename)

	return b.http.PutStream(transport.UserBearer, config.BlobUrl+hash, reader, filename)
}

// SyncComplete notifies that the sync is done
func (b *BlobStorage) SyncComplete(gen int64) error {
	log.Info.Println("TODO: sync in root")
	return nil
	// req := model.SyncCompletedRequest{
	// 	Generation: gen,
	// }
	// return b.http.Post(transport.UserBearer, config.SyncComplete, req, nil)
}

func (b *BlobStorage) WriteRootIndex(roothash string, gen int64) (int64, error) {
	log.Info.Println("writing root with gen: ", gen)
	req := model.BlobRootStorageRequest{
		Broadcast:  true, //TODO
		Hash:       roothash,
		Generation: gen,
	}
	var res model.BlobRootStorageResponse

	err := b.http.Put(transport.UserBearer, config.RootPut, req, &res)
	if err != nil {
		return 0, err
	}
	if res.Hash != roothash {
		return 0, fmt.Errorf("hash mismatch")
	}

	return res.Generation, nil
}
func (b *BlobStorage) GetRootIndex() (string, int64, error) {
	var res model.BlobRootStorageResponse
	err := b.http.Get(transport.UserBearer, config.RootGet, nil, &res)
	if err != nil {
		return "", 0, err
	}

	log.Info.Println("got root gen:", res.Generation)
	return res.Hash, res.Generation, nil

}
