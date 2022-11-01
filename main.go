package main

import "os"

// Set SRC and DST
// SRC="postgres://USER:PASSWORD@SRC_HOST:SRC_PORT/DB?PARAMS"
// DST="postgres://USER:PASSWORD@DST_HOST:DST_PORT/DB?PARAMS"
// SLOT is the name of the replication slot from SRC_CONN
// SUB is the name of the subscription from DST_CONN
func main() {
	SlotName := os.Getenv("SLOT")
	SubName := os.Getenv("SUB")
	CalculateReplicationLag(SlotName, SubName)

	CompareSize()
	CompareCount()
	//CompareRelTuples()
	//CompareSequenceID()
}
