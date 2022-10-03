package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func CompareRelTuples() {
	var relTuples int64
	for _, t := range tables {
		rows, err := srcDB.Query(`SELECT reltuples::bigint FROM pg_catalog.pg_class WHERE relname = '` + t.Name + `'`)
		if err != nil {
			log.Printf("Table:%s Error:%v", t.Name, err)
			continue
		}
		ok := rows.Next()
		if !ok {
			log.Fatal("Empty response from SQL query")
		}
		err = rows.Scan(&relTuples)
		if err != nil {
			log.Printf("Table:%s Error:%v", t.Name, err)
			continue
		}
		t.SourceRelTuples = relTuples

		rows, err = destDB.Query(`SELECT reltuples::bigint FROM pg_catalog.pg_class WHERE relname = '` + t.Name + `'`)
		if err != nil {
			log.Fatal(err)
		}
		ok = rows.Next()
		if !ok {
			log.Fatal("Empty response from SQL query")
		}
		err = rows.Scan(&relTuples)
		if err != nil {
			log.Fatal(err)
		}
		t.DestRelTuples = relTuples

		fmt.Println(t.Name, t.SourceRelTuples, t.DestRelTuples)
	}

	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(tabw, "%s\t%s\t%s\t%s\t\n", "Name", "Src RelTuples", "Dest RelTuples", "Difference")
	for _, s := range tables {
		if s.SourceRelTuples == 0 {
			continue
		}
		if s.DestRelTuples == 0 {
			continue
		}

		diff := s.SourceRelTuples - s.DestRelTuples
		if diff < 0 {
			diff = -diff
		}
		diffPer := (float64(diff) / float64(s.SourceRelTuples)) * 100.0

		fmt.Fprintf(tabw, "%s\t%d\t%d\t%.2f\t\n", s.Name, s.SourceRelTuples, s.DestRelTuples, diffPer)
	}
	tabw.Flush()
}
