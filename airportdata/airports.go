package airportdata

import(
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"log"	
)

//reads the airportdata.pb.go function and returns map of airports
func GetAirportList() map[string]*Airport{

	//get name of airport file and read in.
	var fname string = "airportdata/airports-proto.bin"
	airportIn, err := ioutil.ReadFile(fname)

	//file read error check
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found.", fname)
		} else {
			log.Fatalln("Error reading file:", err)
		}
	}

	
	// set struct to type AirportList from definition in airportdata.pb.go file 
	airportList := &AirportList{}
	
	//unmarshall the data
	unmarshallResult := proto.Unmarshal(airportIn, airportList);

	if unmarshallResult != nil {
		log.Fatalln("Failed to parse airportData:", err)
	}

	//create a map where key = airportcode and value = airportStruct
	airportMap := make(map[string]*Airport)

	//iterate through list and populate the map
	for _,airport := range airportList.Airport {
		airportMap[airport.GetCode()] = airport
	}

	return airportMap
}



