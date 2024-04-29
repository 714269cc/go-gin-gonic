package gin

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

type mockFileSystem struct {
	open func(name string) (http.File, error)
}

func (m *mockFileSystem) Open(name string) (http.File, error) {
	return m.open(name)
}

func TestOnlyFilesFS_Open(t *testing.T) {
	var testFile *os.File
	mockFS := &mockFileSystem{
		open: func(name string) (http.File, error) {
			return testFile, nil
		},
	}
	fs := &OnlyFilesFS{fs: mockFS}

	file, err := fs.Open("foo")

	assert.NoError(t, err)
	assert.Equal(t, testFile, file.(neutralizedReaddirFile).File)
}

func TestOnlyFilesFS_Open_err(t *testing.T) {
	testError := errors.New("mock")
	mockFS := &mockFileSystem{
		open: func(_ string) (http.File, error) {
			return nil, testError
		},
	}
	fs := &OnlyFilesFS{fs: mockFS}

	file, err := fs.Open("foo")

	assert.ErrorIs(t, err, testError)
	assert.Nil(t, file)
}

func Test_neuteredReaddirFile_Readdir(t *testing.T) {
	n := neutralizedReaddirFile{}

	res, err := n.Readdir(0)

	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestDir_listDirectory(t *testing.T) {
	testRoot := "foo"
	fs := Dir(testRoot, true)

	assert.Equal(t, http.Dir(testRoot), fs)
}

func TestDir(t *testing.T) {
	testRoot := "foo"
	fs := Dir(testRoot, false)

	assert.Equal(t, &OnlyFilesFS{fs: http.Dir(testRoot)}, fs)
}
