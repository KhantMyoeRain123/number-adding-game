# REST API

### Errors
Response
```json
{
   "status":"error",
   "message":"This is an error message." 
}
```

### Make Room
#### /api/room/make
Response
```json
{
	"status": "success",
	"room_code": "77K43I",
	"room_state": 0, //0 means WAITING, 1 means STARTED, 2 means ENDED
    "host":true
}
```


### Get Room
#### /api/room/get/{roomCode}
Response
```json
{
	"status": "success",
	"room_code": "77K43I",
	"room_state": 0,
    "host":false 
}
```

# Websocket Events

### Error Event
```json
{
	"type": "error",
	"payload": {
        "errors":[...]
    }
}
```
This event is returned whenever there is an error.

### Initialize Event (Client->Server)
```json
{
	"type": "initialize",
	"payload": {
        "username":"James",
        "host":true,
        "room_code": "77K43I",
	    "room_state": 0
    }
}
```
This event asks the server to generate a ```"playerId"``` for the player and also put them in the room indicated by the ```"room_code"```.
### Initialized Event (Server->Client)
```json
{
	"type": "initialized",
	"payload": {
        "playerId":"7A84",
        "username":"James",
        "host":true,
        "room_code": "77K43I",
	    "room_state": 0
    }
}
```
This is the server's reply to the player who sent ```initialized``` event, indicating success.

### NewUser Event (Server -> Clients)
```json
{
	"type": "newUser",
	"payload": {
        "playerId":"7A84",
        "username":"Jeff",
    }
}
```
The server sends this to players other than the joining player after ```initialized``` is handled successfully.




# Initialization Flow
### Making a Room
The player presses the Make Room button which will send a POST request to the endpoint ```/api/room/make```. When the endpoint returns a success then the frontend will ask for a websocket connection to the server through the route ```ws://localhost:8080/ws```.
### Joining a Room
The player presses the Join Room button which will send a POST request to the endpoint ```/api/room/get/{roomCode}```. When the endpoint returns a success then the frontend will ask for a websocket connection to the server through the route ```ws://localhost:8080/ws```.
### After Joining or Making Room
After the websocket connection is established, the frontend will have to send an ```initialize``` event to the server. The server should respond with the ```initialized``` event to the sender of the ```initialize``` event and if there are already players in the room that was joined, send the ```newUser``` event to those players. Upon receiving ```initialized``` or ```newUser``` the frontend should update the UI appropriately.

# Game Flow
