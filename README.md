# Sendify Code Test

## How to run

This is a simple Go http server so you only have to run this from the command line: 

`go run *.go`

To create a shipment run this CURL command:

`curl --location --request POST 'http://localhost:8080/shipments' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sender": "DE",
    "receiver": "DE",
    "weight": 28.8,
    "customerId": "hans"
}'`

To list shipments you have to first create a shipment and use the customerId in the query params:

`curl --location --request GET 'http://localhost:8080/shipments?customerId=hans'`

## Discussion: Deployment

Currently the app is in a very simple prototype state with an in memory database. Several steps are needed to make it production ready:

1. Add a real database to the service, for example Postgres or MySQL. 
    a. **Local**: I would use a Postgres docker image as well as docker compose to start the db locally and init it with a database schema.
    b. **Cloud**: Google Cloud has a Cloud SQL service where you simply need the credentials and some configuration to get it to connect to your service.
    c. Using GORM we can define tables as Go structs and have the database schema be auto migrated to the db if there are any changes. That way we can version control the database schema.
2. Dockerize the application so that it is more "microservice-friendly".
    a. Build the Go binary.
    b. Build the Docker file and upload it to a registry.
    c. Deploy the Docker image to Cloud Run which is a serverless service on GCP.
    d. In Cloud Run you can add a custom domain and serve that as the external API.
3. Other points to consider is API docs and automating it as much as possible and keeping the docs as close to the code as possible. A good thing would be to create an API docs page in Swagger or other frameworks.
4. Unit testing the business logic, integration tests for the database and acceptance testing for the API itself.
5. Some form of authentication/authorization middleware to ensure that customer A can't access the records of customer B for example. Firebase Auth is a GCP service that can provide and handle JWT tokens for this purpose.

## Scaling up globally

This is complex and requires a long checklist of many things to consider but I will just mention a few here:

1. Build pipelines and automated tests. Managed by Terraform or some other IaaC tool.
2. Monitoring, logging and alerts for resource usage and rate of HTTP 4xx errors.
3. Docs need to be kept up to date. Possibly add tutorials and quick guides for easy integration towards the APIs. 
4. Daily backups of databases to ensure disaster recovery.
5. Global load balancer to reduce latency for multi-regional requests.
6. API tier pricing models for consumers that use API more than others/have a heavy load.
7. Traceability using unique request IDs within the system for fault tracing.
8. Check for security vulnerabilities on a regular basis using linters, static code analysis, depedency checkers etc.
9. Using caching like Redis to reduce network load and improve speed for common type of operations like looking up regions/country codes. 
