# deleteOldFiles
Delete old files based on a specific month, week, and day from FTP 

## config
for run deleteOldFiles put your config in path `configs/configs.json`
```json
{
  "version": 1,
  "ftp_account": {
    "address": "127.0.0.1:21",
    "username": "",
    "password": "",
    "deletable_files_path":""
  }
}
```
## config in main
in main you can change some config fields

```go
func main() {
	jsonData, err := ioutil.ReadFile("configs/configs.json")
  ...
	deleteConfig := &configs.DeleteOldFileConfig{
		DeleteAfter_Month:    3,
		DeleteAfter_Week:     3,
		DeleteAfter_Days:     3,
		KeepLastFileInFolder: true,
		DeleteEmptyFolder:    true,
	}

}
```
- change `configs.json` path
- **DeleteAfter_Month** delete files after x month and save lastest file
- **DeleteAfter_Week** delete files after x week and save latest file
- **DeleteAfter_Days** delete filest afte x days and save latest file
- **KeepLastFileInFolder** keep last file in folder ( mean 1 file in folder )
- **DeleteEmptyFolder** remove empty folders
