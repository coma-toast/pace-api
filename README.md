# pace-api

This is the API endpoint for the PaCE app. This API will be a liaison between the iOS app and the DB data


TODO: 

Demo:

* DB
* API
* Basic app

Pre-launch:
* Import take off sheet data as inventory
* DB of all shapes

TODO:
* When no items found, return an empty json instead of `null`
* Return a `status` as well as the json (maybe?)


## Service file:
```
[Unit]
Description=API server for PaCE
After=network.target

[Service]
Type=simple
ExecStart=/bin/bash /home/jason/www-data/pace-api/pace-api.sh
TimeoutStartSec=0

[Install]
WantedBy=default.target
```
## Script to start:

`/home/jason/www-data/pace-api/pace-api.sh`
```
#!/bin/bash
/home/jason/www-data/pace-api/pace-api -conf=/home/jason/www-data/pace-api/
```
