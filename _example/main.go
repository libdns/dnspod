package main

import (
	"context"
	"fmt"
	"os"
	"time"

	dnspod "github.com/libdns/dnspod"

	"github.com/libdns/libdns"
)

func main() {
	token := os.Getenv("DNSPOD_TOKEN")
	if token == "" {
		fmt.Printf("DNSPOD_TOKEN not set\n")
		return
	}
	zone := os.Getenv("ZONE")
	if zone == "" {
		fmt.Printf("ZONE not set\n")
		return
	}
	provider := dnspod.Provider{APIToken: token}

	records, err := provider.GetRecords(context.TODO(), zone)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	testName := "libdns-test"
	testID := ""
	for _, record := range records {
		fmt.Printf("%s (.%s): %s, %s\n", record.Name, zone, record.Value, record.Type)
		if record.Name == testName {
			testID = record.ID
		}

	}

	if testID != "" {
		// fmt.Printf("Delete entry for %s (id:%s)\n", testName, testID)
		// _, err = provider.DeleteRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
		// 	ID: testID,
		// }})
		// if err != nil {
		// 	fmt.Printf("ERROR: %s\n", err.Error())
		// }
		// Set only works if we have a record.ID
		fmt.Printf("Replacing entry for %s\n", testName)
		_, err = provider.SetRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
			Type:  "TXT",
			Name:  testName,
			Value: fmt.Sprintf("Replacement test entry created by libdns %s", time.Now()),
			TTL:   time.Duration(600) * time.Second,
			ID:    testID,
		}})
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	} else {
		fmt.Printf("Creating new entry for %s\n", testName)
		_, err = provider.AppendRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
			Type:  "TXT",
			Name:  testName,
			Value: fmt.Sprintf("This is a test entry created by libdns %s", time.Now()),
			TTL:   time.Duration(600) * time.Second,
		}})
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	}
}
