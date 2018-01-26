# SQLImporter

SQLImporter is not a database migration tools. This is tools is for importing and executing sql files into test database.

This works is inspired By Hendra Huang database integration test: https://github.com/Hendra-Huang/databaseintegrationtest.

## What is this for?

To test integration with database. When we have a function that interract with database, the schema is imported and can be dropped afterwards.

## Sql file for SQLImporter

To differentiate each queries in sql files, `--end` need to be added.

```sql
CREATE TABLE IF NOT EXISTS `something1` (
    `something1_id` bigint(20) NOT NULL
)
--end

CREATE TABLE IF NOT EXISTS `something2` (
    `something2_id` bigint(20) NOT NULL,
    `field1` varchar(10) NOT NULL
)
--end
```