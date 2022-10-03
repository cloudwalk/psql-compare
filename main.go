package main

import "os"

// Set SRC_CONN and DST_CONN
// SRC_CONN="postgres://USER:PASSWORD@SRC_HOST:SRC_PORT/DB?PARAMS"
// DST_CONN="postgres://USER:PASSWORD@DST_HOST:DST_PORT/DB?PARAMS"
// SLOT_NAME is the name of the replication slot from SRC_CONN
// SUB_NAME is the name of the subscription from DST_CONN
func main() {
	SlotName := os.Getenv("SLOT_NAME")
	SubName := os.Getenv("SUB_NAME")

	//CompareSize()
	//CompareCount()
	// CompareRelTuples()
	CalculateReplicationLag(SlotName, SubName)
	//CompareSequenceID()
}
