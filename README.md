fleetState Code Test
=================

## Summary

A simple web server, that can receive geolocation of vehicles and stream their speed and location

## Start service
  ### Docker
  In order to srart the service, a `Docker` should be installed. 

 How to install `Docker` 
  ``` 
  https://docs.docker.com/get-docker/
```
  If `Docker` is installed, please, run command,s from service root directory
  ```
  docker-compose -f docker/docker-compose.yaml up
``` 

  This command will run 2 docker containers on a local computer. 
  - container with `vehicle_simulator` (service that simulates vehicles movement)
  - container with `fleet_state` (expose port :8080)

## Example JSON Requests

  ### Vehicles
  List of all vehicles, that are sending information to service

  Request:
```GET: http://localhost:8080/v1/vehicle```

  Response:
```javascript
[
  {"vehicleID":"WAUMGAFL1DA105812"},
  {"vehicleID":"1G6KD54Y33U246700"},
  ...
  }
]
```

  ### POST vehicle
  Receives and store latest information about vehicle location. Calculates speed.

  Request:
```POST: http://localhost:8080/v1/vehicle/{vid}```

  Body:
```javascript
{
   "lat":15.322598,
   "lon":22.015082
}
```

Response: 
```
  StatusCode
```
  
  ### Get One vehicle stream
  Get One vehicle stream is streaming information about vehicle (approximately once every second, but can be done faster, if data will be send more frequently)

  Request:
```GET: http://localhost:8080/v1/vehicle/{vid}/stream```

Response: 
```javascript
   {
      "vehicleID":"2FMDK3JC9ABB70648",
      "latitude":15.322598,
      "longitude":22.015082,
      "speed":68,
      "measurement":"kmph"
   }
``` 

## Assumptions during development
 ## Database
 
 Based on project description and data, service will be working with, I decided not to use any database, instead all necessary data has been written into
 in memory storage. Since there was no requirement to keep any vehicle state, but latest, there's no overload on service with writing data to
 Database, nor read from it.
 
 As an improvement, storage like Reddis can be used in future.
 
 ## Service
 Service is calculating speed based on geolocation changes. Assumption was, that vehicle will send updates every second, but service still puts a time,
 when info was received for more precise speed calculation (in case of info deley). Since info is received in `float64` format, speed calculation is fast.
 Also service is working asynchronously, that will allow us to use server resources efficiently.
 
## Improvement on service
 Add unit tests. Add authentication and authorisation mechanisms. Add load balancing with `sticky sessions` before service. Add `swagger`.
 Add posibility to stream one vehicle to multiple endpoints.

