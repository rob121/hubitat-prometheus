package main

import (
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (


	levelOp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "hubprom_level_state",
			Help:      "Switch States",
		},
		[]string{
			// Which user has requested the operation?
			"device",
			// Of what type is the operation?
		},
	)

	batteryLevel = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "hubprom_battery_state",
			Help:      "Battery States",
		},
		[]string{
			// Which user has requested the operation?
			"device",
			// Of what type is the operation?
		},
	)

	switchOp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "hubprom_switch_state",
			Help:      "Switch States",
		},
		[]string{
			// Which user has requested the operation?
			"device",
			// Of what type is the operation?
		},
	)
	motionOp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "hubprom_motion_state",
			Help:      "Motion sensor states",
		},
		[]string{
			// Which user has requested the operation?
			"device",
			// Of what type is the operation?
		},
	)

	cpuUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hubprom_cpu_load",
		Help: "CPU Load",
	},[]string{
		"device",
    })

	deviceStats = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hubprom_device_stats",
		Help: "Device Stats",
	},[]string{
		"device",
		"name",
		"unit",
	})

	dbSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "hubprom_db_size",
		Help:      "Size of DB on disk",
	},[]string{
		"device",
	})

	jvmFree = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "hubprom_jvm_free",
		Help:      "JVM Free",
	},[]string{
		"device",
	})

	freeMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "hubprom_mem_free",
		Help:      "Hub Memory Free",
	},[]string{
		"device",
	})

	hubTemp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "hubprom_temperature",
		Help:      "Hub Temp State",
	},[]string{
		"device",
	})

)

func runPrometheus() {

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9911", nil)
}

