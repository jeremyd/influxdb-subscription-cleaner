# This handy script will DELETE ALL INFLUXDB SUBSCRIPTIONS!
Workaround for https://github.com/influxdata/kapacitor/issues/870
Kapacitor will re-create automatically it's currently used subscriptions.

## Build
go build cleaner.go

## Configure using Environment variables:
INFLUXDB_URL=http://myinflux:8086
INFLUXDB_DRYRUN=true (optional, output what we would have done)

## Run
```
INFLUXDB_URL=http://localhost:8086 ./cleaner.go
```
