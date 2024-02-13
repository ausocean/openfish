# Getting started

## Option 1: Using docker
Openfish can be built from the docker images to run locally on your computer or to deploy to the cloud.

If you are not planning on developing code this is the best choice for you.

### Prerequisites
You will need installed:
- Docker v24.0.5 or later
- Linux (Windows is currently untested but may work)

### Steps:
1) Build docker images:
   ```bash
   docker build . -t openfish-webapp -f ./docker/webapp
   docker build . -t openfish-api -f ./docker/api
   ```

2) Run both containers:
   ```bash
   docker run -p 8080:8080 openfish-api          
   docker run -p 80:80 openfish-webapp
   ```
3) Open the browser and visit http://localhost.

## Option 2: For development
If you want to contribute code to the project, use these steps instead to get started:

### Prerequisites
You will need installed:
- Linux (Windows is currently untested but may work)
- node v19.2.0 or later
- pnpm v8.5.1 or later
- go 1.20 or later

### Steps:
1) Install npm dependencies:
   ```bash
   pnpm install
   ```

2) Start the API:
   ```bash
   go run ./api
   ```

3) Start the webapp using vite's development server:
   ```bash
   pnpm --filter ./openfish-webapp dev
   ```

4) Open the browser and visit http://localhost:5173/watch.html.

5) (Optional) Serve the documentation website using:
    ```bash
    pnpm --filter ./docs dev
    ```



