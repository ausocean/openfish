# Getting started

Openfish can be built from the docker images to run locally on your computer or to deploy to the cloud.

## Prerequisites
You will need installed:
- [Docker](https://www.docker.com) v24.0.5 or later
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- Linux (Windows is currently untested but may work)

## Deployment instructions
### Steps:
1) Setup google cloud credentials so your application can access the datastore:
   ```bash
   gcloud auth login
   gcloud config set project openfish-dev
   gcloud auth application-default login
   ```
   if the environment variable `$GOOGLE_APPLICATION_CREDENTIALS` is not set:
   ```bash
   GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/application_default_credentials.json
   ```

2) Build docker images:
   ```bash
   docker build . -t openfish-site -f ./docker/site.dockerfile
   docker build . -t openfish-api -f ./docker/api.dockerfile
   ```

3) Run both containers:
   ```bash      
   docker run -p 80:80 openfish-site
   docker run -p 8080:8080 -e GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcloud.json -v $GOOGLE_APPLICATION_CREDENTIALS:/tmp/gcloud.json:z openfish-api
   ```
4) Open the browser and visit http://localhost.
