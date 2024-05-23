package golang

import (
	"GolangPractice/utils/logger"
	"fmt"
	"syscall"
	"testing"
	"time"
)

const (
	timeFormatDay = "2006-01-02"

	TimeFormatSecond = "2006-01-02 15:04:05"

	TimeFormatLog = "2006/01/02 15:04:05"

	TimeFormatMillisecond = "2006-01-02 15:04:05.000"

	TimeFormatISO = "2006-01-02T15:04:05Z"
)

func TestTimeFormat(t *testing.T) {
	// Format current time.
	nowStamp := time.Now()

	timeStr := nowStamp.Format(TimeFormatMillisecond)
	logger.Infoln("Current Time: %s", timeStr)
}

func TestTimeParse(t *testing.T) {
	// Assign location.
	loc, _ := time.LoadLocation("Local") // UTC， or other valid location name

	// another way to get local time zone is time.Local.
	// UTC is accessible if using time.UTC.
	timeStr, _ := time.ParseInLocation(TimeFormatLog, "2021/11/02 15:04:05", loc)
	logger.Infoln("Local time.Time: %s", timeStr)

	NewYorkLoc, _ := time.LoadLocation("America/New_York")

	// you can either define a format by yourself or use predefined patterns in `time` package.
	timeStr, _ = time.ParseInLocation(time.RFC3339, "2021/11/02 15:04:05", NewYorkLoc)
	logger.Infoln("New York time.Time: %s", timeStr)
}

func TestTimeCalc(t *testing.T) {
	logger.Infoln(time.Now().AddDate(1,1,1).String())
}

func TestGetFileCreationTime(t *testing.T) {
	var st syscall.Stat_t
	fileFullPath := "/root/demo.txt"
	err := syscall.Stat(fileFullPath, &st)
	if err != nil {
		fmt.Println(err)
		return
	}
	createTime := time.Unix(st.Ctim.Sec, 0)
	fmt.Println("file is created at ", createTime)
}
