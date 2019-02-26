lab2.1

Description: 
My own implementation of a wall of clocks. This works using to go programs, clock2.go and clockWall.go. clock2.go is the server for the time and clockWall.go displays the city and the time at the city. clockWall.go works with goroutines inside a for loop, to start a goroutine for each parameter recieved. 

Compilation:

run 'go build clock2.go' and start the servers you want. 3 examples: 

TZ=US/Eastern ./clock2 -port 8010
TZ=Asia/Tokyo ./clock2 -port 8020
TZ=Europe/London ./clock2 -port 8030

Once the servers you want are running, you can run 'go run clockWall.go' passing as parameters the information of the servers. An example using the servers started above:

go run clockWall.go NewYork=localhost:8010 Tokyo=localhost:8020 London=localhost:8030