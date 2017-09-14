GO code to request spark dashboard page, parse it and extract stats of 
running apps in a prometheus format and CSV too : the main objective is 
to compensate the lack of completude of spark 1.5.1 API about running 
apps stats.

# Usage

Spark dashboard url must be defined in a env var named `SPARK_DASHBOARD_URL`.
Example :
````
export SPARK_DASHBOARD_URL="http://<spark_url>/"
export SPARK_LOGIN=login
export SPARK_PASSWORD=pass
make run
````

By default, create a HTTP server that listen on port 8080. To change it,
simply define an env var `PORT` with the given value.

# Dev mode

To test locally without have a running spark available in network, you can 
run a mock page by running a second server described in mock folder

````
cd mock
go run main.go
export SPARK_DASHBOARD_URL=http://localhost:8088/myApp/appStreamingStatistics.html
export SPARK_DASHBOARD_URL=http://localhost:8088/mainPage/mainpage.html
cd ..
make run
````
