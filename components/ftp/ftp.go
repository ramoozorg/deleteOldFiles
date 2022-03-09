package ftp

import (
	"bufio"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"time"
)

type FTP struct {
	Connection *ftp.ServerConn `json:"connection"`
}

func (f *FTP) Connect(config *FtpConfig) error {
	c, err := ftp.Dial(config.Address, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	if err = c.Login(config.Username, config.Password); err != nil {
		return err
	}
	f.Connection = c
	return nil
}

func (f *FTP) List(path string, recursive bool) []*Entry {
	entries, err := f.Connection.List(path)
	if err != nil {
		panic(err)
	}

	files := []*Entry{}

	for _, entry := range entries {
		if entry.Type != ftp.EntryTypeFolder || isValidDirectory(entry.Name) {
			files = append(files, &Entry{entry.Name, path + "/" + entry.Name, entry.Size, entry.Type, entry.Time})

		} else if recursive && entry.Type == ftp.EntryTypeFolder && isValidDirectory(entry.Name) {
			subFiles := f.List(path+"/"+entry.Name, recursive)
			files = append(files, subFiles...)
		}
	}
	return files
}

func (f *FTP) Store(path string, r io.Reader) error {
	return f.Connection.Stor(path, r)
}

func (f *FTP) Quit() error {
	return f.Connection.Quit()
}

func (f *FTP) Rename(from string, to string) error {
	return f.Connection.Rename(from, to)
}

func (f *FTP) Delete(path string) error {
	return f.Connection.Delete(path)
}

func (f *FTP) DeleteDirectory(path string, recursive bool) error {
	if recursive {
		return f.Connection.RemoveDirRecur(path)
	} else {
		return f.Connection.RemoveDir(path)
	}
}

func (f *FTP) DownloadFile(path string, localFilePath string) error {
	r, err := f.Connection.Retr(path)
	if err != nil {
		return err
	}
	defer r.Close()

	//Create and open local file
	outputFile, _ := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer outputFile.Close()

	// Create a buffered reader of remote file
	reader := bufio.NewReader(r)

	//p is the max buffer size
	p := make([]byte, 1024*4)

	//save bytes in file while reading buffer
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		outputFile.Write(p[:n])
	}
	return nil
}

func (f *FTP) MakeDir(path string) error {
	return f.Connection.MakeDir(path)
}

func isValidDirectory(directoryName string) bool {
	return directoryName != "." && directoryName != ".."
}
