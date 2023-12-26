## App Description:
* My iteration of this take home challenge reads all provided files within the [data package](data/)
* A CitiesMap is created and stored in the client with the cities located inside [cites.json](data/cities.json)
* Tmp cities located [here](data/tmp), are compared and validated against the map
* Files for both valid elements, invalid elements and unprocessable files are then generated inside of the [results package](results)
* This app can be run using either `go run` or a docker container.
* Due to time constraints, there is only partial testing of some of the functions

## App Directions:
* Install Go
  - For Linux:
    - Run `sudo apt-get update
sudo apt-get install golang
`
  - For Mac (Homebrew)
    - Run `brew install go
`
* From the CITY_VALIDATION_s3 directory, run `go mod tidy` to install external dependencies
* From the root of the directory run `go run city_validation_s3.go`
* Check the results package for the generated files

## Docker:
* Download Docker Desktop and make sure it is running
* From the root of the app run `docker-compose up --build`
* Check the results package for the generated files


