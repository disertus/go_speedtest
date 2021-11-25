# go_speedtest

This set of tools includes two main components:
- **speedtest scheduler** - server responsible for the scheduling of all speedtests and returning them to clients (under construction)
- **speedtest runner** - client responsible for fetching of the schedule from server, perfoming scheduled speedtests and sending the results to a Postgres db


**speedtest runer** is a tiny client binary which:
- fetches a personal schedule of speedtests from a server every once in a while
- performs scheduled speedtests based on Ookla's speedtest CLI tool 
- ingests speedtest's json output into a Postgres db

