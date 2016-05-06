# IP Locator API Server

This is a API server that will return information about any IP address using the maxmind geolocation database in JSON format.

Getting it running is simple, simply build it using the build script or using `go build` and run the binary. Then you can access the API at `http://localhost:3000/v1/`. You can also append an ip address or hostname to that URL to get information on those. You can also add `ssl.cert` and `ssl.key` files on the main binary path to enable SSL on port 3001.
