# Build stage.
FROM node:20.5-alpine as build-stage
WORKDIR /src
RUN npm install -g pnpm

# Install dependencies - copy done seperately to help caching.
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm i

# Build webapp.
COPY openfish-webapp ./openfish-webapp
RUN pnpm --filter ./openfish-webapp build

# Production container.
FROM nginx as production-stage
COPY --from=build-stage /src/openfish-webapp/dist /usr/share/nginx/html
