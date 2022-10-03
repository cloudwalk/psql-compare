# psql-compare
This is used to compare data between two databases

## Functions:
**CompareSize**: This just compares the storage size of tables, but storage size can be very different between two databases. So its not a very reliable metric, but this function is important to get the list of all the tables. As its just a single query to databases, its very quick.

**CompareCount**: This function runs `count(*)` query on all the tables in the databases. That is why this function is very very slow, but its quite a reliable way to compare the two tables.

**CompareRelTuples**: This function runs `SELECT reltuples::bigint FROM pg_catalog.pg_class` query on all the tables in the databases. This query is quicker than CompareCount function but less reliable.

**CompareSequenceID**: This function runs `MAX(id)` query on all the tables in the databases. This is also quite a reliable metric, but it will ignore any table which does not have an id as integer.

**CalculateReplicationLag**: This function calculates the replication lag for the replication. You need to set `SLOT_NAME` and `SUB_NAME` environment variable for this to work.

## ENV variables:

Set the following Environment variables:
- SRC_CONN="postgres://USER:PASSWORD@SRC_HOST:SRC_PORT/DB?PARAMS"
- DST_CONN="postgres://USER:PASSWORD@DST_HOST:DST_PORT/DB?PARAMS"
- SLOT_NAME is the name of the replication slot from SRC_CONN
- SUB_NAME is the name of the subscription from DST_CONN
