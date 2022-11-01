package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	_ "github.com/lib/pq"
)

var srcDB, destDB *sql.DB

func init() {
	var err error
	srcDB, err = sql.Open("postgres", os.Getenv("SRC"))
	if err != nil {
		log.Fatal(err)
	}
	destDB, err = sql.Open("postgres", os.Getenv("DST"))
	if err != nil {
		log.Fatal(err)
	}
}

type Table struct {
	Name            string
	SourceSize      int64
	DestinationSize int64
	SizeDifference  float64
	SourceMaxId     int64
	DestMaxId       int64
	SourceRelTuples int64
	DestRelTuples   int64
	SourceCount     int64
	DestCount       int64
}

var tables []*Table

func CompareSize() {
	rows, err := srcDB.Query(`select table_name, pg_relation_size(quote_ident(table_name))
		from information_schema.tables
		where table_schema = 'public'
		order by 2;`)
	if err != nil {
		log.Fatal(err)
		return
	}
	for rows.Next() {
		var s Table
		err = rows.Scan(&s.Name, &s.SourceSize)
		if err != nil {
			log.Fatal(err)
		}
		if s.Name == "pg_stat_statements" {
			continue
		}

		tables = append(tables, &s)
	}

	rows, err = destDB.Query(`select table_name, pg_relation_size(quote_ident(table_name))
		from information_schema.tables
		where table_schema = 'public'
		order by 2;`)
	if err != nil {
		log.Fatal(err)
		return
	}

	var name string
	var size int64

	for rows.Next() {
		err = rows.Scan(&name, &size)
		if err != nil {
			log.Fatal(err)
		}

		if name == "pg_stat_statements" {
			continue
		}

		for _, s := range tables {
			if s.Name == name {
				s.DestinationSize = size
			}
		}
	}

	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(tabw, "%s\t%s\t%s\t%s\t\n", "Name", "Src Size", "Dest Size", "Difference")
	for _, s := range tables {
		if s.SourceSize == 0 && s.DestinationSize == 0 {
			continue
		}
		s.SizeDifference = float64(s.SourceSize-s.DestinationSize) / float64(s.SourceSize) * 100
		if s.SizeDifference < 0 {
			s.SizeDifference = -s.SizeDifference
		}
		if s.SizeDifference < 30 {
			continue
		}

		fmt.Fprintf(tabw, "%s\t%d\t%d\t%.2f\t\n", s.Name, s.SourceSize, s.DestinationSize, s.SizeDifference)
	}
	tabw.Flush()
}
