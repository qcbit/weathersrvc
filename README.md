# Weather Service

A service that serves forecasted weather for a latitude and longitude pair responding with a short forecast for the area for today, whether it be "Sunny", "Partly Cloudy", etc., and a characterization "hot", "cold", or "moderate" for the temperature.

It uses the [National Weather Service API Web Service](https://www.weather.gov/documentation/services-web-api) as the data source.

## Sample Data

Example latitude and longitude coordinates (39.7456,-97.0892).

## How to use

From the terminal, run the program:

```make install```

```make run```

From the terminal, stop the program:

```control+c```

### Get the weather

Add coordinates (latitude, longitude) to the endpoint:

```http://localhost:8080/location/latitude,longitude```

Use curl:

```curl http://localhost:8080/location/39.7456,-97.0892```

\- or -

Use a browser and go to <http://localhost:8080/location/39.7456,-97.0892>.

\- or -

Use [Postman](https://www.postman.com/downloads/).

## Get Your Coordinates

<https://www.gps-coordinates.net/my-location>
