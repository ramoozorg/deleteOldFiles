package managers

import (
	"fmt"
	ftp2 "github.com/jlaffaye/ftp"
	"os"
	"ramooz.org/deleteOldBackup/components/ftp"
	"ramooz.org/deleteOldBackup/configs"
	"sort"
	"time"
)

func DeleteOldFiles(path string, deleteConfig *configs.DeleteOldFileConfig, ftpConfig *ftp.FtpConfig) {
	ftpAccount := &ftp.FTP{}
	if err := ftpAccount.Connect(ftpConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deleteOldFiles(path, deleteConfig, ftpAccount)
}
func deleteOldFiles(path string, deleteConfig *configs.DeleteOldFileConfig, ftpAccount *ftp.FTP) {
	files := ftpAccount.List(path, false)
	if len(files) == 0 && deleteConfig.DeleteEmptyFolder {
		ftpAccount.Delete(path)
		return
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModifyTime.UnixMilli() < files[j].ModifyTime.UnixMilli()
	})
	remainingFilesCount := 0
	for _, file := range files {
		if file.Type == ftp2.EntryTypeFile {
			remainingFilesCount++
		}
	}
	for _, file := range files {
		switch file.Type {
		case ftp2.EntryTypeFile:
			if remainingFilesCount == 1 && deleteConfig.KeepLastFileInFolder {
				fmt.Println("file kept:", file.Path)
				continue
			}
			if isDeletableFile(file, deleteConfig) {
				fmt.Println("file deleting:", file.Path)
				if err := ftpAccount.Delete(file.Path); err != nil {
					fmt.Println("error occurred on deleting file:", err)
				} else {
					fmt.Println("file deleted:", file.Path)
					remainingFilesCount--
				}
			}
		case ftp2.EntryTypeFolder:
			deleteOldFiles(file.Path, deleteConfig, ftpAccount)
		}
	}
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
func isDeletableFile(file *ftp.Entry, deleteConfig *configs.DeleteOldFileConfig) bool {
	now := time.Now()
	now = date(now.Year(), int(now.Month()), now.Day())
	fileTime := file.ModifyTime
	fileTime = date(fileTime.Year(), int(fileTime.Month()), fileTime.Day())
	diffInDays := int(now.Sub(fileTime).Hours() / 24)
	diffInMonths := now.Year()*24 + int(now.Month()) - (fileTime.Year()*24 + int(fileTime.Month()))
	diffInWeeks := int(now.Sub(fileTime).Hours() / 24 / 7)

	if fileTime.Day() == 1 && diffInMonths < deleteConfig.DeleteAfter_Month {
		return false
	}
	if fileTime.Weekday() == time.Saturday && diffInWeeks < deleteConfig.DeleteAfter_Week {
		return false
	}

	if diffInDays < deleteConfig.DeleteAfter_Days {
		return false
	}
	return true
}
