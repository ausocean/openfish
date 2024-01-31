# Getting started

## Option 1: Using docker
Openfish can be built from the docker images to run locally on your computer or to deploy to the cloud.

If you are not planning on developing code this is the best choice for you.

### Prerequisites
You will need installed:
- [Docker](https://www.docker.com) v24.0.5 or later
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- Linux (Windows is currently untested but may work)

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
   docker build . -t openfish-webapp -f ./docker/webapp
   docker build . -t openfish-api -f ./docker/api
   ```

3) Run both containers:
   ```bash      
   docker run -p 80:80 openfish-webapp
   docker run -p 8080:8080 -e GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcloud.json -v $GOOGLE_APPLICATION_CREDENTIALS:/tmp/gcloud.json:z openfish-api
   ```
4) Open the browser and visit http://localhost.

## Option 2: For development
If you want to contribute code to the project, use these steps instead to get started:

### Prerequisites
You will need installed:
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- [node](https://nodejs.org/en) v19.2.0 or later
- [pnpm](https://pnpm.io/) v8.5.1 or later
- [go](https://go.dev/) 1.20 or later
- Linux (Windows is currently untested but may work)

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

2) Install npm dependencies:
   ```bash
   pnpm install
   ```

3) Start the API:
   ```bash
   go run ./api
   ```

4) Start the webapp using vite's development server:
   ```bash
   pnpm --filter ./openfish-webapp dev
   ```

5) Open the browser and visit http://localhost:5173/watch.html.

5) (Optional) Serve the documentation website using:
    ```bash
    pnpm --filter ./docs dev
    ```



