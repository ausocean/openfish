FROM node:20.5-alpine as build-stage
WORKDIR /app
RUN npm install -g pnpm
COPY package*.json ./
RUN pnpm i
COPY ./ .
RUN pnpm build

FROM nginx as production-stage
RUN mkdir /app
COPY --from=build-stage /app/dist /usr/share/nginx/html
