# go_speedtest

go_speedtest is a tiny binary that ingests the json output of the speedtest CLI tool into a struct.

Its full list of features should include:
- fetch a personal schedule of speedtests from a server via a GET request with one's mac-address passed as parameter
- insert the results of performed speedtests into a Postgres database
