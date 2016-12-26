package main

import(
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"github.com/golang/protobuf/proto"
	placedata "placedata"
)

//made global so that reading the protobuf file and putting data into map is only done once
var placesMap map[string][]*placedata.Place

func main(){

	//call GetPlaceList and store the Place struct (defined in placedata.pb.go)
	placesMap = placedata.GetPlaceList()

	//set up of route to handle get requests
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getAllPlaces/{state}/{town}", getAllPlaces).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getAllPlaces(w http.ResponseWriter, r *http.Request) {

	//get state and town params that were passed
	params := mux.Vars(r);
	var state string = strings.ToUpper(params["state"])
	var town string = strings.ToLower(params["town"])
	
	//find list in map that contains our key(state) and iterate through all the values
	for _ , value := range placesMap[state] {

		//lowercase prefix search 
		if(strings.HasPrefix(strings.ToLower(value.GetName()), town)){
			
			//set up the protobuf struct
        	place := &placedata.Place{
        		State:  proto.String(value.GetState()),
        		Name:   proto.String(value.GetName()),
        		Lat:    proto.Float64(value.GetLat()),
        		Lon:    proto.Float64(value.GetLon()),
        	}

        	//marshall to []byte and write it
        	var msg []byte
        	var err error

        	msg, err = proto.Marshal(place)

        	if(err != nil){
        		 log.Fatal("problem occured while marshaling: ", err)
        	}

        	w.Write(msg)
        }
    }

}