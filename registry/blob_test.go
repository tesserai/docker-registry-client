package registry_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"

	digest "github.com/opencontainers/go-digest"
)

func TestRegistry_UploadBlob(t *testing.T) {
	blobData := []byte("This is a test blob.")
	digest := digest.FromBytes(blobData)

	foreachWritableTestcase(t, func(t *testing.T, tc *TestCase) {
		content := bytes.NewBuffer(blobData)
		ctx := context.Background()
		err := tc.Registry(t).UploadBlob(ctx, tc.Repository, digest, content, nil)
		if err != nil {
			t.Error("UploadBlob() failed:", err)
		}
	})
}

func TestRegistry_UploadBlobFromFile(t *testing.T) {
	const filename = "testdata/blob"

	// prepare UploadBlob() parameters
	blobData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	digest := digest.FromBytes(blobData)
	body := func() (io.ReadCloser, error) {
		// NOTE: the file will be closed by UploadBlob() (more precisely the http.Client)
		return os.Open(filename)
	}

	// run tests
	foreachWritableTestcase(t, func(t *testing.T, tc *TestCase) {
		blobReader, err := body()
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()
		err = tc.Registry(t).UploadBlob(ctx, tc.Repository, digest, blobReader, body)
		if err != nil {
			t.Error("UploadBlob() failed:", err)
		}
	})
}
