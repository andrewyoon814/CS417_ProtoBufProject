package placedata

import(
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"log"
)

func GetPlaceList() map[string][]*Place{

	//get name of places file and read into the placesIn variable
	var fname string = "placedata/places-proto.bin"
	placesIn, err := ioutil.ReadFile(fname)

	//file read error check
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("%s: File not found.", fname)
		} else {
			log.Fatalln("Error reading file:", err)
		}
	}

	// set struct to type PlaceList from definition in airportdata.pb.go file 
	placesList := &PlaceList{}
	
	//unmarshall the data from placesIn to placesList
	unmarshallResult := proto.Unmarshal(placesIn, placesList);

	if unmarshallResult != nil {
		log.Fatalln("Failed to parse placeData:", err)
	}

	//create a map of slices
	//each key(state) has a value(slice containing cities in that state)
	//this allows for less searching
	placeMap := make(map[string][]*Place)

	//iterate through list and populate the map
	for _,place := range placesList.Place {
		placeMap[place.GetState()] = append(placeMap[place.GetState()], place)
	}

	return placeMap
}



