# Openfish Webapp
Openfish Webapp is a web interface for interacting with the Openfish API.

## Deploying openfish
1) Build and run the web app using docker:
   ```bash
   docker build ./openfish-webapp -t openfish
   docker run -p 80:80 openfish 
   ```
2) Start go server:
   ```bash
   go run ./api/ 
   ```
3) Open the browser and visit http://localhost.

## Development guide
#### Getting started

Have the following prerequisites installed on your system:
- node v19.2.0 or later
- pnpm v8.5.1 or later
- go 1.20 or later

1) Start go server:
   ```bash
   go run ./api/ 
   ```

2) Install all dependencies in package.json:
   ```bash
   cd openfish-webapp
   pnpm i
   ```

3) Start the live development server:
   ```bash
   pnpm dev
   ```
   Visit http://localhost:5173/. Changes are updated live.

#### Linting / formatting
- `pnpm fmt` to format code.
- `pnpm check` to check for common issues.
