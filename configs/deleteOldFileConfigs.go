package configs

type DeleteOldFileConfig struct {
	DeleteAfter_Days     int
	DeleteAfter_Month    int
	DeleteAfter_Week     int
	DeleteEmptyFolder    bool
	KeepLastFileInFolder bool
}
