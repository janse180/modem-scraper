[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

What does this do?
==========
This application polls these pages on the Arris SB8200 cable modem:
* http://192.168.100.1/cmconnectionstatus.html
* http://192.168.100.1/cmswinfo.html
* http://192.168.100.1/cmeventlog.html

The data from those pages is populated into structs, and is
then published out to MQTT as well as InfluxDB. The event log is
also published and a local BoltDB database is leveraged to prevent
recording duplicates.

My intent is to use this data in Home Assistant (MQTT) as well
as InfluxDB/Grafana (for graphing metrics over longer periods
of time).

This is currently a work in progress... and was built for my own
use. That said, if it's useful for someone else, then cool beans.

This has recently been updated to incorporate the work from https://github.com/mdonoughe/modem_status
in order to work around the lateset Comcast-pushed firmware update for these
modems which adds in an ancient SSL certificate plus a username and password
login page.

TODO:
* Add unit tests.
* Build and publish docker container automatically.
