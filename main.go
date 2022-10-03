package main

// Set SRC_CONN and DST_CONN
// SRC_CONN="postgres://USER:PASSWORD@SRC_HOST:SRC_PORT/DB?PARAMS"
// DST_CONN="postgres://USER:PASSWORD@DST_HOST:DST_PORT/DB?PARAMS"
func main() {
	CompareSize()
	CompareCount()
	CompareRelTuples()
	CompareSequenceID()
}
