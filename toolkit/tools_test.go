package toolkit

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Error("Expected string length of 10, but got", len(s))
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "allowed no rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    false,
		errorExpected: false,
	},
	{
		name:          "allowed rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    true,
		errorExpected: false,
	},
	{
		name:          "not allowed",
		allowedTypes:  []string{"image/png"},
		renameFile:    true,
		errorExpected: true,
	},
}

func TestTools_UploadFiles(t *testing.T) {
	// Очищаем директорию перед тестом
	if err := os.RemoveAll("./testdata/uploads"); err != nil {
		t.Fatal(err)
	}

	// Создаем директорию для загрузки
	if err := os.MkdirAll("./testdata/uploads", 0755); err != nil {
		t.Fatal(err)
	}

	for _, e := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer writer.Close()

			file, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
				return
			}
			defer file.Close()

			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
				return
			}

			_, err = io.Copy(part, file)
			if err != nil {
				t.Error(err)
				return
			}
		}()

		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadedFiles(request, "./testdata/uploads", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			t.Logf("Uploaded file: %+v", uploadedFiles[0])
			filePath := fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)
			t.Logf("Checking file at: %s", filePath)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("%s: Expected file to exist: %s", e.name, uploadedFiles[0].NewFileName)
			} else {
				t.Logf("File exists at: %s", filePath)
			}

			src := filePath
			dst := fmt.Sprintf("./testdata/uploads/keep_%s", uploadedFiles[0].NewFileName)
			t.Logf("Moving file from %s to %s", src, dst)
			if err := os.Rename(src, dst); err != nil {
				t.Errorf("Error moving file: %v", err)
			}
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: error expected but none received", e.name)
		}

		wg.Wait()
	}
}

func TestTools_UploadOneFile(t *testing.T) {

	// Очищаем директорию перед тестом

	if err := os.RemoveAll("./testdata/upload-one"); err != nil {
		t.Fatal(err)
	}

	// Создаем директорию для загрузки
	if err := os.MkdirAll("./testdata/upload-one", 0755); err != nil {
		t.Fatal(err)
	}

	for _, e := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		go func() {
			defer writer.Close()

			file, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
				return
			}
			defer file.Close()

			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
				return
			}

			_, err = io.Copy(part, file)
			if err != nil {
				t.Error(err)
				return
			}
		}()

		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools

		uploadedFile, err := testTools.UploadOneFile(request, "./testdata/upload-one", true)
		if err != nil {
			t.Error(err)
		}

		if !e.errorExpected {
			t.Logf("Uploaded file: %+v", uploadedFile)
			filePath := fmt.Sprintf("./testdata/upload-one/%s", uploadedFile.NewFileName)
			t.Logf("Checking file at: %s", filePath)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("%s: Expected file to exist: %s", e.name, uploadedFile.NewFileName)
			} else {
				t.Logf("File exists at: %s", filePath)
			}

			src := filePath
			dst := fmt.Sprintf("./testdata/upload-one/keep_%s", uploadedFile.NewFileName)
			t.Logf("Moving file from %s to %s", src, dst)
			if err := os.Rename(src, dst); err != nil {
				t.Errorf("Error moving file: %v", err)
			}
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: error expected but none received", e.name)
		}
	}
}
