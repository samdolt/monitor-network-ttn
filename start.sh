#!/usr/bin/bash
sudo systemctl start redis
sudo systemctl start influxdb
sudo systemctl start chronograf
TTN_APP_ID=bfhtest1 TTN_APP_ACCESS_KEY=ttn-account-v2.REPLACE_ME_WITH_A_KEY ./monitor-network-ttn
