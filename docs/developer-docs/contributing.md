# Openfish Contributing Guide
**Authors:** Scott Barnard

Hi there! 👋 We’re thrilled that you’re interested in contributing to OpenFish, our open-source project dedicated to classifying marine species.

Before diving in, please take a moment to familiarize yourself with our guidelines. These will help streamline the contribution process and ensure a smooth collaboration. Here’s what you need to know:

## Raising issues.
If you find a bug or want to request a feature, please raise an issue https://github.com/ausocean/openfish/issues.

For bugs, try to include any information that you think is helpful in describing the problem and ideally the steps you need to reproduce the issue. The better the information you include, the quicker it is likely to be fixed.

Use the labels `document`, `bug` and `enhancement` to categorize your issue.

## Choosing an issue to work on.
Issues with the `good first issue` label are a good choice for your first contribution.

## Repo Setup
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
   pnpm site install
   ```

3) Start the OpenFish API:
   ```bash
   go run ./cmd/openfish
   ```

4) Start the site using vite's development server:
   ```bash
   pnpm site dev
   ```

5) Open the browser and visit http://localhost:5173/watch.

6) (Optional) Serve the documentation website using:
    ```bash
    pnpm docsite install
    pnpm docsite dev
    ```

## Common tasks
### Formatting code
::: code-group
```bash [Go packages]
go fmt ./... && swag fmt
```
```bash [TypeScript packages]
pnpm fmt
```
:::

### Running unit tests
::: code-group
```bash [Go packages]
go test -v ./... -short
```
```bash [TypeScript packages]
pnpm test
```
:::

### Linting (TypeScript packages only)
::: code-group
```bash [Checking code]
pnpm check
```
```bash [Applying fixes automatically]
pnpm check --fix
```
:::

## Submitting a pull request (PR)
- If the PR is related to an issue, link to it in the description.
- Check that your code passes unit tests.
- Check that your code is formatted properly.
- Check that your code does not have any linting issues.

## Summary
Thank you for joining the OpenFish community! We can’t wait to see your contributions in action 🌊🐠🦑!
