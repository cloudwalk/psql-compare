package main

import (
	"database/sql"
	"fmt"
	"log"
)

func CalculateReplicationLag(slotName, subName string) {
	rows, err := srcDB.Query(`SELECT stat.pid, stat.client_addr, stat.state, stat.sync_state, stat.write_lsn,
		pg_wal_lsn_diff(pg_current_wal_lsn(), stat.write_lsn)::BIGINT AS replication_lag,
		pg_wal_lsn_diff(sent_lsn, stat.write_lsn)::BIGINT AS write_lag,
		pg_wal_lsn_diff(sent_lsn, stat.flush_lsn)::BIGINT AS flush_lag,
		pg_wal_lsn_diff(sent_lsn, stat.replay_lsn)::BIGINT AS replay_lag
		FROM pg_catalog.pg_stat_replication stat
		JOIN pg_catalog.pg_replication_slots slot ON (stat.pid = slot.active_pid)
		WHERE slot.slot_name = $1;`, slotName)
	if err != nil {
		log.Fatal(err)
	}
	ok := rows.Next()
	if !ok {
		log.Fatal("No replication slot found")
	}

	var pid, client_addr, state, sync_state, write_lsn, replication_lag, write_lag, flush_lag, replay_lag string
	err = rows.Scan(&pid, &client_addr, &state, &sync_state, &write_lsn, &replication_lag, &write_lag, &flush_lag, &replay_lag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pid:%s client_addr:%s state:%s sync_state:%s write_lsn:%s replication_lag:%s write_lag:%s flush_lag:%s replay_lag:%s \n", pid, client_addr, state, sync_state, write_lsn, replication_lag, write_lag, flush_lag, replay_lag)

	rows, err = destDB.Query(`SELECT stat.*,
		pg_wal_lsn_diff(stat.received_lsn, $1)::BIGINT AS replication_lag
		FROM pg_catalog.pg_stat_subscription stat
		WHERE subname = $2`, write_lsn, subName)
	if err != nil {
		log.Fatal(err)
	}
	ok = rows.Next()
	if !ok {
		log.Fatal("No subscription found")
	}

	var subid, subname, pidDest, relid, received_lsn, last_msg_send_time, last_msg_receipt_time, latest_end_lsn, latest_end_time, replication_lagDest sql.NullString
	err = rows.Scan(&subid, &subname, &pidDest, &relid, &received_lsn, &last_msg_send_time, &last_msg_receipt_time, &latest_end_lsn, &latest_end_time, &replication_lagDest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("subid:%s subname:%s pid:%s relid:%s received_lsn:%s last_msg_send_time:%s last_msg_receipt_time:%s latest_end_lsn:%s latest_end_time:%s replication_lag:%s \n", subid.String, subname.String, pidDest.String, relid.String, received_lsn.String, last_msg_send_time.String, last_msg_receipt_time.String, latest_end_lsn.String, latest_end_time.String, replication_lagDest.String)

	fmt.Println("Replication lag:", replication_lagDest.String)
}
