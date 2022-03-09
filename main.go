package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"ramooz.org/deleteOldBackup/components/ftp"
	"ramooz.org/deleteOldBackup/configs"
	"ramooz.org/deleteOldBackup/managers"
)

func main() {
	jsonData, err := ioutil.ReadFile("configs/configs.json")
	if err != nil {
		fmt.Println("config file not found", err)
		os.Exit(1)
	}

	serviceConfig := map[string]interface{}{}

	if err := json.Unmarshal(jsonData, &serviceConfig); err != nil {
		fmt.Println("config file didn't parsing to json", err)
		os.Exit(1)
	}
	ftpConfigMap := serviceConfig["ftp_account"].(map[string]interface{})
	ftpConfig := &ftp.FtpConfig{
		Address:  ftpConfigMap["address"].(string),
		Username: ftpConfigMap["username"].(string),
		Password: ftpConfigMap["password"].(string),
	}
	deleteConfig := &configs.DeleteOldFileConfig{
		DeleteAfter_Month:    3,
		DeleteAfter_Week:     3,
		DeleteAfter_Days:     3,
		KeepLastFileInFolder: true,
		DeleteEmptyFolder:    true,
	}
	managers.DeleteOldFiles(ftpConfigMap["deletable_files_path"].(string), deleteConfig, ftpConfig)

}
