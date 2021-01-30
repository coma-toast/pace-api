# pace-api

This is the API for the PaCE app. This API will be a liaison between the iOS app and the Firestore DB. This a demo app and may not be fully functional. 


TODO: 

Demo:

* DB
* API
* Basic app

Pre-launch:
* Import take off sheet data as inventory
* DB of all shapes


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
