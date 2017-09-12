TODO

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
