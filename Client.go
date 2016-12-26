package main

import(
	"log"
	"math"
	"os"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	placedata "placedata"
	airportdata "airportdata"
)

func main(){

	var state string 
	var town string

	//make sure 2 arguments are passed and store into above variables
	if(len(os.Args) != 3){
		fmt.Println("Incorrect number of arguments.")
		return
	}else{
		
		state = os.Args[1]
		town = os.Args[2]
	}
	
	//make get call to the placeServer with args
	res, err := http.Get("http://localhost:8080/getAllPlaces/" + state + "/" + town)

    if(err != nil){
       log.Fatalln("Connection to placeServer could not be made")
    }

    //read body of byte array
    body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	//create a struct of type Place() and unmarshall the body msg
	place := &placedata.Place{}
	unmarshalErr := proto.Unmarshal(body, place)

	if (unmarshalErr != nil ) {
		log.Fatal("problem unmarshaling: ", unmarshalErr)
	}

	//save the longitude and latitude recieved and call airportServer
	var longitude float64 = place.GetLon()
	var latitude float64 = place.GetLat()

	fmt.Println("Place that was inputed was : ")
	fmt.Println(place)
	fmt.Println("");


	//make get call to the airportServer with args
	airportRes, err := http.Get("http://localhost:8081/getClosestAirports/" + strconv.FormatFloat(latitude, 'f', -1, 64) + "/" + strconv.FormatFloat(longitude, 'f', -1, 64))

    if(err != nil){
       log.Fatalln("Connection to airportServer could not be made")
    }

    //read body of marshalled byte array from airportServer
    airportBody, err := ioutil.ReadAll(airportRes.Body)
	res.Body.Close()


	if err != nil {
		log.Fatal(err)
	}


	//create a struct of type AirportList() and unmarshall the body msg
	airports := &airportdata.AirportList{}
	airportUnmarshalErr := proto.Unmarshal(airportBody, airports)

	if (airportUnmarshalErr != nil ) {
		log.Fatal("problem unmarshaling: ", airportUnmarshalErr)
	}

	fmt.Println("5 Closest Airports are : ")
	for _,airport := range airports.Airport {
		fmt.Println(airport, "DISTANCE: ", getDistance(latitude,longitude,airport.GetLat(), airport.GetLon()), " miles")
	}
	fmt.Println("")
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
	
	return (math.Acos(s1*s2+c1*c2*clong) * 180 * 60 * 1.1507794 ) / math.Pi
}

