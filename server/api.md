# REST API

### Errors
Response
```json
{
   "status":"error"
   "message":"This is an error message." 
}
```

### Make Room
#### /api/room/make
Response
```json
{
    "status":"success"
    "roomCode":"R8PEP9" //this is an example room code
}
```

### Join Room
#### /api/room/join/{roomCode}
Response
```json
{
    "status":"success"
    "roomCode":"R8PEP9" //this is an example room code
}
```

# Websocket Events




# Initialization Flow
### Making a Room
The player presses the Make Room button which will send a POST request to the endpoint ```/api/room/make```. When the endpoint returns a success then the frontend will ask for a websocket connection to the server through the route ```ws://localhost:8080/ws```.
### Joining a Room
The player presses the Join Room button which will send a POST request to the endpoint ```/api/room/join```. When the endpoint returns a success then the frontend will ask for a websocket connection to the server through the route ```ws://localhost:8080/ws```.
### After Joining or Making Room
After the websocket connection is established, the frontend will have to send an ```initialize``` event to the server. The server should respond with the ```initialized``` event to the sender of the ```initialize``` event and if there are already players in the room that was joined, send the ```newUser``` event to those players. Upon receiving ```initialized``` or ```newUser``` the frontend should update the UI appropriately.

# Game Flow
