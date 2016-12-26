Andrew Yoon and Gabe Raybould
CS 417 Distributed Systems Assignment 1

This project was implemented in golang with use Google Protocol Buffers and Rest protocol.
The purpose was for us to familiarize ourselves with Go, serialization with Google Protocol Buffers, and RPC.

Code files:

	-main/client.go
	-main/placeServer.go
	-main/airportServer.go
	-airportdata/airports.go
	-placedata/places.go

Compilation and Information: 

	On three different terminals, you need to run airportServer.go, placeServer.go, and client.go.

	-go run airportServer.go
	-go run placeServer.go
	-go run client.go {state} {place}
		ex: -go run client.go nj princeton
		    -go run client.go ak anchorage

	placeServer.go listens to port 8080.
	airportServer.go listens to port 8081.

	Client.go takes in two parameters for state and place and checks if there are a indeed 2 inputs. 
		Please supply corrected formated state and place names
	Client.go then calls placeServer.go via an HTTP Get at the URL: "http://localhost:8080/getAllPlaces/" + state + "/" + town.
		ex:) "http://localhost:8080/getAllPlaces/nj/princeton"
	PlaceServer.go calls GetPlaceList() function at places.go. Places.go unmarshalls the list and send data back to PlaceServer.go. 
	PlaceServer.go find the given state and town by matching the state exactly (state is converted to all uppercases) and uses string.HasPrefix() to match the town.
	PlaceServer.go returns the data for the place that is found to Client.

	Client then calls airportServer.go via an HTTP Get at URL: "http://localhost:8081/getClosestAirports/{latitude}/{longitude}".
		example: "http://localhost:8081/getClosestAirports/40.361506/-74.651988"
	Using these coordinates, airportServer calls airport.go's GetClosestAirports method and recieves the five closest airports to the given place.

	Client then prints out the information to the client terminal.


Storage and Searching of Data:

	In both placeServer and airportServer, we use Golang's implementation of maps to store and search for data.

	The protocol buffers return lists for both the places.proto and airport.proto. We iterate the lists and save as a map and send the maps back to the servers.

	PLACE's are stored in maps of place object slices where the key is the state. So, if you search the map for the key "NJ" you will get a slice containing all the airports
	in New Jersey. You then iterate through that slice to get the place you were looking for. 

	AIRPORT's are also stored in maps where the key is the unique airport code and the value is a pointer to an airport struct. In airport server, in order to search the 5 closest
	airports, we created a new struct with the airport code and a distance field. We iterate through the airportMap and populate a slice of these structs that contain every single
	airports airport code and distances. We then use sort.Sort() on this slice and sort based on distance from smallest to largest. Then we take the first 5 elements and pass back
	to the client.

Bugs and Peculiarities

	If you enter a two word place such as "New Brunswick" or "Los Angeles" please surround with quotes.
  
Testcases can be viewed in testcases.txt file.
