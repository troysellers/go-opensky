package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

// Response is a struct
type Response struct {
	Time   int64
	States []State
}

// State is a struct .. defined https://opensky-network.org/apidoc/rest.html#operation
type State struct {
	Icao24        string  // [0] Unique ICAO 24-bit address of the transponder in hex string representation.
	Callsign      string  // [1] Callsign of the vehicle (8 chars). Can be null if no callsign has been received.
	OriginCountry string  // [2] Country name inferred from the ICAO 24-bit address.
	TimePosition  float64 // [3] Unix timestamp (seconds) for the last position update. Can be null if no position report was received by OpenSky within the past 15s.
	TimeVelocity  float64 // [4]Unix timestamp (seconds) for the last velocity update. Can be null if no velocity report was received by OpenSky within the past 15s.
	Longitude     float64 // [5] WGS-84 longitude in decimal degrees. Can be null.
	Latitude      float64 // [6] WGS-84 latitude in decimal degrees. Can be null.
	Altitude      float64 // [7] Barometric or geometric altitude in meters. Can be null.
	OnGround      bool    // [8] boolean	Boolean value which indicates if the position was retrieved from a surface position report.
	Velocity      float64 // [9]Velocity over ground in m/s. Can be null.
	Heading       float64 // [10] Heading in decimal degrees clockwise from north (i.e. north=0°). Can be null.
	VerticalRate  float64 //[11] Vertical rate in m/s. A positive value indicates that the airplane is climbing, a negative value indicates that it descends. Can be null.
	Sensors       []int   // [12]IDs of the receivers which contributed to this state vector. Is null if no filtering for sensor was used in the request.
}

func (st *State) String() string {
	return fmt.Sprintf(`Icao24 [%v] Callsign [%v] OriginCountry[%v] TimePosition[%v] `+
		`TimeVelocity[%v] Longitude[%v] Latitude[%v] Altitude[%v] OnGround[%v] `+
		`Velocity[%v] Heading[%v] VerticalRate[%v] Sensors[%v]`,
		st.Icao24, st.Callsign, st.OriginCountry, st.TimePosition,
		st.TimeVelocity, st.Longitude, st.Latitude, st.Altitude,
		st.OnGround, st.Velocity, st.Heading, st.VerticalRate, st.Sensors)
}

func main() {
	var url string
	if len(os.Args) != 3 {
		url = "https://opensky-network.org/api/states/all"
	} else {
		url = "https://" + os.Args[1] + ":" + os.Args[2] + "@opensky-network.org/api/states/all"
	}
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	fmt.Printf("Status : %v\n", resp.Status)

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		check(err)
		f, err := os.Create("temp.txt")
		check(err)
		defer f.Close()
		f.Write(body)

		var str interface{}
		err = json.Unmarshal(body, &str)
		check(err)

		m := str.(map[string]interface{})
		fmt.Printf("Lets get the time [%d]\n", m["time"])
		fmt.Printf("Lets get the states [%v]\n", reflect.TypeOf(m["states"]))

		st := m["states"].([]interface{})

		states := make([]*State, 0)

		for _, st := range st {
			stArray := st.([]interface{})
			if len(stArray) != 13 {
				fmt.Println("We don't have 13 state variables. We have ", len(stArray))
				panic("ARGH!!")
			}

			state := &State{}
			if stArray[0] != nil {
				state.Icao24 = stArray[0].(string)
			}
			if stArray[1] != nil {
				state.Callsign = stArray[1].(string)
			}
			if stArray[2] != nil {
				state.OriginCountry = stArray[2].(string)
			}
			if stArray[3] != nil {
				state.TimePosition = stArray[3].(float64)
			}
			if stArray[4] != nil {
				state.TimeVelocity = stArray[4].(float64)
			}
			if stArray[5] != nil {
				state.Longitude = stArray[5].(float64)
			}
			if stArray[6] != nil {
				state.Latitude = stArray[6].(float64)
			}
			if stArray[7] != nil {
				state.Altitude = stArray[7].(float64)
			}
			state.OnGround = stArray[8].(bool)
			if stArray[9] != nil {
				state.Velocity = stArray[9].(float64)
			}
			if stArray[10] != nil {
				state.Heading = stArray[10].(float64)
			}
			if stArray[11] != nil {
				state.VerticalRate = stArray[11].(float64)
			}
			if stArray[12] != nil {
				state.Sensors = stArray[12].([]int)
			}

			states = append(states, state)

		}
		process(states)
		fmt.Printf("We have [%d] states\n", len(states))

		if len(os.Args) != 3 {
			fmt.Println("You can pass your own user credentials to this function, if you want to.")
			fmt.Println("e.g 'go run opensky.go username supersecretpassword' (without the single quotes...)")
		}
	}
}
func check(err error) {
	if err != nil {
		fmt.Printf("We have an error : %v\n", err)
		panic(err)
	}
}

func process(states []*State) {

	for _, st := range states {
		fmt.Println(st.String())
	}

}
