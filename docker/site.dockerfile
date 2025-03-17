# Build stage.
FROM node:23.10-alpine as build-stage
WORKDIR /src
RUN npm install -g pnpm

# Install dependencies - copy done seperately to help caching.
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm i

# Build site.
COPY site ./site
RUN pnpm site build

# Production container.
FROM nginx as production-stage
COPY --from=build-stage /src/site/dist /usr/share/nginx/html
