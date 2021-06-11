# Hubitat Prometheus

A prometheus collector for hubitat

## Configuration

```
{
  "access_token": "d7bcc9daa2f8e",
  "hubitat_base_url": "http://192.168.20.5/apps/api/12",
  "hubitat_ip": "192.168.20.5"
}
```


#### Sample Data:

```
# TYPE hubprom_battery_state gauge
hubprom_battery_state{device="master_closet_motion"} 37
# HELP hubprom_device_stats Device Stats
# TYPE hubprom_device_stats gauge
hubprom_device_stats{device="laundry_washer_power",name="power",unit="W"} 1.155
hubprom_device_stats{device="system_weather_driver",name="feelsLike",unit="°F"} 79
# HELP hubprom_motion_state Motion sensor states
# TYPE hubprom_motion_state gauge
hubprom_motion_state{device="office_ceiling_motion"} 0
# HELP hubprom_switch_state Switch States
# TYPE hubprom_switch_state gauge
hubprom_switch_state{device="kitchen_light"} 1
# HELP hubprom_temperature Hub Temp State
# TYPE hubprom_temperature gauge
hubprom_temperature{device="system_weather_driver"} 77
```
