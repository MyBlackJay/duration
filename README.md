# ISODURATION
It's multifunctional and fast a Go module for parsing [ISO 8601 duration format](https://www.digi.com/resources/documentation/digidocs/90001488-13/reference/r_iso_8601_duration_format.htm)


## Why use?
ISO 8601 duration format is a convenient format for recording the duration of something, for example, you can use this format in your configurations to indicate the duration of the wait or for time-based events.

## Features
- fast parsing of raw strings in ISO 8601 duration format
- convenient tools for obtaining and reverse conversion of time.Duration
- possibility to get each period and time element in float64 format
- human-readable errors open for import and comparison
- yaml serialization and deserialization
- json serialization and deserialization 

## Installation
```
go get github.com/MyBlackJay/isoduration
```

## Example

#### [Code:](https://go.dev/play/p/SPc-oa4lcNi)
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/MyBlackJay/isoduration"
	"gopkg.in/yaml.v3"
	"time"
)

type Config struct {
	FlushThresholdDuration *isoduration.Duration `yaml:"flush_threshold_duration" json:"flush_threshold_duration"`
}

type Main struct {
	Config Config `yaml:"config" json:"config"`
}

func yamlToStruct() {
	yml := `
config:
  flush_threshold_duration: PT30M30S
`
	m := &Main{}

	if err := yaml.Unmarshal([]byte(yml), m); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(m.Config.FlushThresholdDuration.ToTimeDuration())
}

func jsonToStruct() {
	jsn := `{"config": {"flush_threshold_duration": "PT30M30S"}}`
	m := &Main{}

	if err := json.Unmarshal([]byte(jsn), m); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(m.Config.FlushThresholdDuration)
}

func structToYaml() {
	tp := &Main{
		Config: Config{
			FlushThresholdDuration: isoduration.NewDuration(1, 1, 5, 0, 1, 30, 12, true),
		},
	}

	if out, err := yaml.Marshal(tp); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(string(out))
	}

}

func structToJson() {
	tp := &Main{
		Config: Config{
			FlushThresholdDuration: isoduration.NewDuration(1, 1, 5, 0, 1, 30, 12, true),
		},
	}

	if out, err := json.Marshal(tp); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(string(out))
	}
}

func timeDelta() {
	str := "PT10H5M"
	strNeg := "-PT10H5M"
	now := time.Now()

	td := time.Hour * 20

	iso, _ := isoduration.ParseDuration(str)
	isoNeg, _ := isoduration.ParseDuration(strNeg)

	fmt.Println(td - iso.ToTimeDuration())
	fmt.Println(td - isoNeg.ToTimeDuration())
	fmt.Println(now.Add(iso.ToTimeDuration()).UTC())
	fmt.Println(now.Add(isoNeg.ToTimeDuration()).UTC())

}

func fromTimeDuration() {
	dur := (time.Hour * isoduration.DayHours * isoduration.YearDays * 2) + (60 * time.Hour) + (30 * time.Second)
	v := isoduration.NewFromTimeDuration(dur)
	c := isoduration.FormatTimeDuration(dur)

	fmt.Println(v, c)
}

func main() {
	yamlToStruct() // 30m30s
	jsonToStruct() // 30m35s
	structToYaml() // config:
	//                  flush_threshold_duration: -P1Y1M5DT1H30M12S
	structToJson()     // {"config":{"flush_threshold_duration":"-P1Y1M5DT1H30M12S"}}
	timeDelta()        // 9h55m0s | 30h5m0s | 2024-01-19 21:25:38.660823592 +0000 UTC | 2024-01-19 01:15:38.660823592 +0000 UTC
	fromTimeDuration() // P2Y2DT12H30S P2Y2DT12H30S
}

```

