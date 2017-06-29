#!/usr/bin/bash
sudo systemctl start redis
sudo systemctl start influxdb
sudo systemctl start chronograf
TTN_APP_ID=bfhtest1 TTN_APP_ACCESS_KEY=ttn-account-v2.IN7kWMC4CxjsPW7NO6gcIYz8by7ai38pyIDtgXzYeBA ./monitor-network-ttn
