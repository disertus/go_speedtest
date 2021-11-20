# go_speedtest

This set of tools includes two main components:
-  **go_speedtest_server** - server responsible for the scheduling of all speedtests (under construction)
- **go_speedtest_client** - client responsible for perfoming the speedtests and sending the results to a Postgres db

**go_speedtest_client** is a tiny client binary which: 
- performs scheduled speedtests relying on Ookla's speedtest CLI tool 
- ingests speedtest's json output into a Postgres db
- fetches a personal schedule of speedtests from a server every once in a while
