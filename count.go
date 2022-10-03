package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"text/tabwriter"
)

var wg sync.WaitGroup
var guard chan int

func CompareCount() {
	guard = make(chan int, 50)

	fmt.Println("Table count:", len(tables))

	for i, t := range tables {
		fmt.Println("START--->", i, t.Name)
		guard <- 1
		go CalculateSrcCount(t)
		guard <- 1
		go CalculateDestCount(t)
	}

	fmt.Println("Waiting")
	wg.Wait()
	fmt.Printf("\n\n\n\nDONE---->\n")

	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(tabw, "%s\t%s\t%s\t%s\t\n", "Name", "SRC Count", "Dest Count", "Difference")
	for _, s := range tables {
		if s.SourceCount == s.DestCount {
			continue
		}
		fmt.Fprintf(tabw, "%s\t%d\t%d\t%d\t\n", s.Name, s.SourceCount, s.DestCount, s.SourceCount-s.DestCount)
	}
	tabw.Flush()
}

func CalculateSrcCount(t *Table) {
	wg.Add(1)
	defer wg.Done()

	rows, err := srcDB.Query(`SELECT count(*) FROM ` + t.Name)
	if err != nil {
		log.Printf("Table:%s Error:%v", t.Name, err)
		<-guard
		return
	}
	ok := rows.Next()
	if !ok {
		log.Fatal("Empty response from SQL query")
	}
	err = rows.Scan(&t.SourceCount)
	if err != nil {
		log.Printf("Table:%s Error:%v", t.Name, err)
		<-guard
		return
	}
	//fmt.Println("SRC:", t.Name, t.SourceCount, t.DestCount)
	<-guard
}

func CalculateDestCount(t *Table) {
	wg.Add(1)
	defer wg.Done()

	rows, err := destDB.Query(`SELECT count(*) FROM ` + t.Name)
	if err != nil {
		log.Printf("Table:%s Error:%v", t.Name, err)
		<-guard
		return
	}
	ok := rows.Next()
	if !ok {
		log.Fatal("Empty response from SQL query")
	}
	err = rows.Scan(&t.DestCount)
	if err != nil {
		log.Printf("Table:%s Error:%v", t.Name, err)
		<-guard
		return
	}
	//fmt.Println("DEST:", t.Name, t.SourceCount, t.DestCount)
	<-guard
}
