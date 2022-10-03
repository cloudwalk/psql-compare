package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func CompareSequenceID() {
	var seqID sql.NullInt64
	for _, t := range tables {
		rows, err := srcDB.Query(`SELECT MAX(id) FROM ` + t.Name)
		if err != nil {
			log.Printf("Table:%s Error:%v", t.Name, err)
			continue
		}
		ok := rows.Next()
		if !ok {
			log.Fatal("Empty response from SQL query")
		}
		err = rows.Scan(&seqID)
		if err != nil {
			log.Printf("Table:%s Error:%v", t.Name, err)
			continue
		}
		t.SourceMaxId = seqID.Int64

		rows, err = destDB.Query(`SELECT MAX(id) FROM ` + t.Name)
		if err != nil {
			//log.Fatal(err)
			log.Printf("Table:%s Error:%v", t.Name, err)
			continue
		}
		ok = rows.Next()
		if !ok {
			log.Fatal("Empty response from SQL query")
		}
		err = rows.Scan(&seqID)
		if err != nil {
			//log.Fatal(err)
			log.Printf("Table:%s Error:%v", t.Name, err)
			continue
		}
		t.DestMaxId = seqID.Int64
	}

	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(tabw, "%s\t%s\t%s\t%s\t\n", "Name", "Src Seq", "Dest Seq", "Difference")
	for _, s := range tables {
		if s.SourceMaxId == s.DestMaxId {
			continue
		}
		fmt.Fprintf(tabw, "%s\t%d\t%d\t%d\t\n", s.Name, s.SourceMaxId, s.DestMaxId, s.SourceMaxId-s.DestMaxId)
	}
	tabw.Flush()
}
