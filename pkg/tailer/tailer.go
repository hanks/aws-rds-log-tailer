package tailer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

const (
	TIME_FORMAT   = "2006-01-02-15"
	LOG_NAME_BASE = "error/postgresql.log."
)

func Execute(out string, dbID string) {
	marker := ""
	cur := time.Now().UTC()
	prev := cur

	for {
		cur = time.Now().UTC()
		// handle hour change and did not finish to
		// download previous hour's log
		if prev.Hour() != cur.Hour() && marker != "" {
			cur = prev
		}

		logFileName := buildLogFileName(cur)
		marker, err := fetchLog(out, dbID, logFileName, marker)

		if err != nil {
			log.Fatalf("error, %s", err)
			os.Exit(1)
		}
		prev = cur
		log.Printf("%s: %s", logFileName, marker)

		// sleep for 30 seconds to avoid to make intensive api calls to rds server
		time.Sleep(30 * time.Second)
	}
}

func buildLogFileName(t time.Time) string {
	suffix := t.Format(TIME_FORMAT)
	logFileName := fmt.Sprintf("%s%s", LOG_NAME_BASE, suffix)

	log.Printf("LogFileName is %s", logFileName)
	return logFileName
}

func fetchLog(path string, dbID string, logFileName string, marker string) (string, error) {
	output, err := downloadLogWithToken(dbID, logFileName, marker)
	if err != nil {
		msg := fmt.Sprintf("download db log error, %s", err)
		log.Fatal(msg)
		return "", errors.New(msg)
	}

	logMsg, marker := *output.LogFileData, *output.Marker
	log.Printf("marker is %s", marker)

	appendLog(path, logMsg)

	return marker, nil
}

func downloadLogWithToken(dbID string, logFileName string, marker string) (*rds.DownloadDBLogFilePortionOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	svc := rds.New(sess)
	var input *rds.DownloadDBLogFilePortionInput
	if marker == "" {
		input = &rds.DownloadDBLogFilePortionInput{
			DBInstanceIdentifier: aws.String(dbID),
			LogFileName:          aws.String(logFileName),
		}
	} else {
		input = &rds.DownloadDBLogFilePortionInput{
			DBInstanceIdentifier: aws.String(dbID),
			LogFileName:          aws.String(logFileName),
			Marker:               aws.String(marker),
		}
	}

	result, err := svc.DownloadDBLogFilePortion(input)

	if err != nil {
		msg := fmt.Sprintf("download db log error, %s", err)
		log.Fatal(msg)
		return nil, errors.New(msg)
	}

	return result, nil
}

func appendLog(path string, s string) error {
	dir, err := filepath.Abs(filepath.Dir(path))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("create new dir <%s>", dir)
		err = os.MkdirAll(dir, 0777)

		if err != nil {
			msg := fmt.Sprintf("open file %s error, %s", path, err)
			log.Fatal(msg)
			return errors.New(msg)
		}

		log.Printf("create new dir <%s> successfully", dir)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		msg := fmt.Sprintf("open file %s error, %s", path, err)
		log.Fatal(msg)
		return errors.New(msg)
	}

	s += "\n"

	w := bufio.NewWriter(f)
	w.WriteString(s)
	w.Flush()

	return nil
}
