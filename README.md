# This handy script will DELETE ALL INFLUXDB SUBSCRIPTIONS!
Workaround for https://github.com/influxdata/kapacitor/issues/870
Kapacitor will re-create automatically it's currently used subscriptions.

## Build (deps are vendored using glide)
go get github.com/jeremyd/influxdb-subscription-cleaner

## Configure using Environment variables:
```
INFLUXDB_URL=http://myinflux:8086
INFLUXDB_DRYRUN=true (optional, output what we would have done)
```

## Run
```
INFLUXDB_URL=http://localhost:8086 cleaner
```