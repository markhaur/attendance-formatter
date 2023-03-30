package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type Attendance struct {
	Status       string
	DeviceName   string
	AuthDateTime string
	AuthDate     string
	AuthTime     string
	CardNO       string
	EmployeeID   string
}

func main() {
	dateString := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("./output%s.csv", dateString)
	fmt.Printf("filename: %s", filename)
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file %v\n", err)
		return
	}
	defer file.Close()

	// create a csv reader
	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	attendanceRecords := make([]Attendance, 0, 0)
	var checkoutTrack = make(map[string]int)
	counter := 0

	record, err := reader.Read()
	for true {
		record, err = reader.Read()
		if len(record) < 7 {
			break
		}
		attendanceRecord := mapSliceToAttendance(record)
		_, exists := checkoutTrack[attendanceRecord.EmployeeID]
		if exists {
			checkoutTrack[attendanceRecord.EmployeeID] = counter
		} else {
			if attendanceRecord.Status == "P20" {
				attendanceRecord.Status = "P10"
			}
			checkoutTrack[attendanceRecord.EmployeeID] = -1
		}
		attendanceRecords = append(attendanceRecords, attendanceRecord)
		counter += 1
	}

	// writing data to excel file
	outputFile, err := os.Create(fmt.Sprintf("./processed_output%s.csv", dateString))
	if err != nil {
		fmt.Printf("can't create file %v\n", err)
	}
	defer outputFile.Close()

	// create a new csv writer
	writer := csv.NewWriter(outputFile)
	writer.Write([]string{"Status", "deviceName", "authDateTime", "authDate", "authTime", "CardNo", "EmployeeID"})

	var savedIndex int
	for index, record := range attendanceRecords {
		savedIndex, _ = checkoutTrack[record.EmployeeID]
		if savedIndex == index {
			writer.Write([]string{"P20", record.DeviceName, record.AuthDateTime, record.AuthDate, record.AuthTime, record.CardNO, record.EmployeeID})
			continue
		}
		writer.Write([]string{record.Status, record.DeviceName, record.AuthDateTime, record.AuthDate, record.AuthTime, record.CardNO, record.EmployeeID})
	}

	writer.Flush()

}

func mapSliceToAttendance(record []string) Attendance {
	return Attendance{Status: record[0], DeviceName: record[1], AuthDateTime: record[2], AuthDate: record[3], AuthTime: record[4], CardNO: record[5], EmployeeID: record[6]}
}
