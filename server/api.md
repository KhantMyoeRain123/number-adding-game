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
#### ```/api/room/make?playerName={playerName}```
Response
```json
{
	"status": "success",
    "playerId":"AD34",
	"room_code": "77K43I",
	"room_state": 0
}
```
### Get Room
#### ```/api/room/join/{roomCode}?playerName={playerName}```
Path Parameters

| Parameter  | Definition | Type
| -------- | ------- | ------- | 
|roomCode | The code of the room the player wants to know exists or not | string|
Response
```json
{
	"status": "success",
    "playerId":"AD34",
	"room_code": "77K43I",
	"room_state": 0
}
```
### Upgrade to Websocket
#### ```ws://localhost:8080/ws?playerId={playerId}```
Query Parameters

| Parameter  | Definition | Type
| -------- | ------- | ------- | 
| playerName  | The username entered by the player    | string |
| playerId  | The id of the player    | string |
| roomCode | The code of the room the player is joining     | string |
| host    | Whether or not the player is the host of the room| bool|



# Websocket Events
### NewUser Event (Server -> Clients)
```json
{
	"type": "newUser",
	"payload": {
        "playerId":"55C4",
        "username":"Jeff",
    }
}
```
The server sends this to players other than the joining player after ```initialized``` is handled successfully.

# Initialization Flow
### Making a Room
The player presses the Make Room button which will send a POST request to the endpoint ```/api/room/make```. When the endpoint returns a success then the frontend will ask for a websocket connection to the server through the route ```ws://localhost:8080/ws``` with the required query parameters.
### Joining a Room
The player presses the Join Room button which will send a POST request to the endpoint ```/api/room/join/{roomCode}```. When the endpoint returns a success then the frontend will ask for a websocket connection to the server through the route ```ws://localhost:8080/ws``` with the required query parameters.

# Game Flow
