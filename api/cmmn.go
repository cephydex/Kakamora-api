package api

import (
	"fmt"
	"log"
	"os"
	"svclookup/xutil"
	"time"
)

func Accra() *time.Location {
	time.Local = time.UTC
	accra, err := time.LoadLocation("Africa/Accra")

	if err != nil {
		// 	panic(err)
		fmt.Println("Error loading location:", err)
	}

	return accra
}

const fPath = "monitoring.txt"

func CheckNCreate() {
	if _, err := os.Stat(fPath); os.IsNotExist(err) {
		fmt.Printf("File: %s does not exist!", fPath)
		_, fe := os.Create("monitoring.txt")
		if fe != nil {
			panic(fe)
		}
	}
}

func AppendToFile(msg_str string) {
	data, fErr := os.ReadFile(fPath)
    if fErr != nil {
        log.Fatal(fErr)
    }
	
    msg := []byte( msg_str)
    data = append(data, msg...)

    err := os.WriteFile(fPath, data, 0777)
    if err != nil {
        log.Fatal(err)
    }
}

func CleanRespItems(arr []xutil.RespItem) []xutil.RespItem {
	mySlice := []xutil.RespItem{}

	for _, item := range arr {
		if item.Site != "" {
			// fmt.Println("VAL CHK", item, i)				
			mySlice = append(mySlice, item)
		}
	}
	
	return mySlice
}