package main

import (
	"fmt"
	"github.com/rob121/hubprom/hubitat"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/rob121/vhelp"
)

var api *hubitat.Api

func main() {

	vhelp.Load("config")

	conf,cerr := vhelp.Get("config")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if(cerr!=nil){

		log.Fatal(cerr)
	}



	devices := make(map[string][]string)

	api =  hubitat.NewApi(conf.GetString("access_token"),conf.GetString("hubitat_base_url"),false)

	go runPrometheus()

	hubitat.Config(conf.GetString("hubitat_ip"), true, devices)

	go pollData()

	go func() {

		for evt := range hubitat.Events {
			fmt.Printf("%#v\n", evt)

			switch evt.Name {
			case "illuminance":
				break
			case "humidity":
				break
			case "battery":
				batteryLevel.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
			case "level":
				levelOp.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "motion":
				motionOp.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "switch":
				switchOp.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "freeMemory":
				freeMemory.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "cpu5Min":
				cpuUsage.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "jvmFree":
				jvmFree.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "temperature":
				hubTemp.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			case "dbSize":
				dbSize.WithLabelValues(toNice(evt.DisplayName)).Set(convValue(evt.Value))
				break
			default:
				deviceStats.WithLabelValues(toNice(evt.DisplayName), evt.Name, evt.Unit).Set(convValue(evt.Value))

			}

		}

	}()

	select {}

}

func pollData(){

	go pollAttributes()

	t := time.NewTicker(5*time.Minute)

	for range t.C {

		pollAttributes()

	}


}

func pollAttributes(){

		 devs, _ := api.Attributes("battery")

		 for _, d := range devs {

			 hubitat.Events <- hubitat.Event{Name:"battery",DisplayName:d.Label,Value:d.Attributes.Battery}

		 }

}

func convValue(val string) float64 {

	if val == "on" || val == "active" {

		val = "1"
	}

	if val == "off" || val == "inactive" {

		val = "0"
	}

	val = strings.TrimSpace(val)

	f, err := strconv.ParseFloat(val, 64)

	if err != nil {

		log.Println("Value Conversion Error: %s", err)
		return 0

	}

	return f

}

func toNice(nm string) string {

	ret := strings.Replace(nm, "-", "_", -1)
	ret = strings.Replace(ret, " ", "_", -1)
	ret = strings.ToLower(ret)

	space := regexp.MustCompile(`(_)+`)
	ret = space.ReplaceAllString(ret, "_")

	return ret

}
