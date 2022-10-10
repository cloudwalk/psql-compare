package main

import (
	"fmt"
	"log"
)

// SELECT c.relname FROM pg_class c WHERE c.relkind = 'S';

func UpdateSeqValue() {
	rows, err := srcDB.Query(`SELECT sequence_schema, sequence_name FROM information_schema.sequences;`)
	if err != nil {
		log.Fatal(err)
		return
	}
	for i := 0; rows.Next(); i++ {
		var schema, seqName string
		var value int64
		err = rows.Scan(&schema, &seqName)
		if err != nil {
			log.Fatal(err)
		}

		rowsValue, err := srcDB.Query(fmt.Sprintf(`SELECT last_value FROM %s.%s;`, schema, seqName))
		if err != nil {
			log.Fatal(err)
		}
		rowsValue.Next()
		err = rowsValue.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}

		_, err = destDB.Query(fmt.Sprintf(`ALTER SEQUENCE %s.%s RESTART WITH %d;`, schema, seqName, value))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(i, schema, seqName, value)
	}
}
