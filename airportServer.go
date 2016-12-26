package main

import(
	"sort"
	"log"
	"net/http"
	"math"
	"strconv"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	airportdata "airportdata"
)

//this struct will be used to compose a slice that will be used to sort airports by distance
type distance struct{
	code string
	distance float64
}
type Distances []distance
var airportMap map[string]*airportdata.Airport

//sort interface function's Len, Less, and Swap
func (slice Distances) Len() int {
	return len(slice)
}

func (slice Distances) Less(i, j int) bool {
	return slice[i].distance < slice[j].distance;
}

func (slice Distances) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func main(){

	//call GetAirportList and store the airportList struct (defined in airportdata.pb.go)
	airportMap = airportdata.GetAirportList()

	//set up of route to handle get requests
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getClosestAirports/{latitude}/{longitude}", getClosestAirports).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func getClosestAirports(w http.ResponseWriter, r *http.Request) {

	//get latitude and longitude params that were passed and parse to float64
	params := mux.Vars(r);
	latitude, err := strconv.ParseFloat(params["latitude"],64)
	longitude, err := strconv.ParseFloat(params["longitude"],64)
	if(err != nil){
		log.Fatal("string to float conversion failed.")
	}

	//this slice of type Distances will hold slice of all airport codes and corrresponding distances
	var closest Distances

	//iterate through airportMap and populate list
	for _,airport := range airportMap {
     	closest = append(closest, distance{airport.GetCode(), getDistance(latitude,longitude,airport.GetLat(),airport.GetLon())})
	}

	//sort the closest array	
	sort.Sort(closest)

	//Airport type doesnt have append method. So populate top 5 airport results with a list...
	var returnList [5]*airportdata.Airport
	returnList[0] = airportMap[closest[0].code]
	returnList[1] = airportMap[closest[1].code]
	returnList[2] = airportMap[closest[2].code]
	returnList[3] = airportMap[closest[3].code]
	returnList[4] = airportMap[closest[4].code]

	//then turn into slices
	returnSlice := returnList[:]

	//create AirportList type and attach slice
	returnMsg := &airportdata.AirportList{
		Airport: returnSlice,
	}


	//matshall to protobuf and send
	var msg []byte

	msg, err = proto.Marshal(returnMsg)

	if(err != nil){
		 log.Fatal("problem occured while marshaling: ", err)
	}

	w.Write(msg)
}

//distance function 
func getDistance(placeLat float64, placeLon float64, airportLat float64, airportLon float64) float64 {

	//distance function used to calculate
	//d = 60 cos-1( sin(lat1) sin(lat2) + cos(lat1) cos(lat2) cos(lon2-lon1))

	//convert degree to radians
	placeLat = placeLat* math.Pi / 180
	placeLon = placeLon* math.Pi / 180
	airportLat = airportLat* math.Pi / 180
	airportLon = airportLon* math.Pi / 180

	//calculate distance
	s1, c1 := math.Sincos(placeLat)
	s2, c2 := math.Sincos(airportLat)
	clong := math.Cos(placeLon - airportLon)
	return math.Acos(s1*s2+c1*c2*clong)
}


