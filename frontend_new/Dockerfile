# build stage
FROM node:lts-alpine as altklausur_ausleihe-frontend-build-stage
WORKDIR /app
COPY package*.json *.config.js index.html ./
RUN npm install
COPY ./src ./src
COPY ./public ./public

ARG SERVER_WS
ENV SERVER_WS "$SERVER_WS"
ARG SERVER_HTTP
ENV SERVER_HTTP "$SERVER_HTTP"
ENV NODE_OPTIONS "--openssl-legacy-provider"

RUN sed -i "s,SERVER\_WS,\"$SERVER_WS\",g" ./src/vue-apollo.js
RUN sed -i "s,SERVER\_HTTP,\"$SERVER_HTTP\",g"  ./src/vue-apollo.js
RUN npm run build

# production stage
FROM nginx:stable-alpine as altklausur_ausleihe-frontend-production-stage
ARG NGINX_CONFIG_NAME
COPY --from=altklausur_ausleihe-frontend-build-stage /app/dist /usr/share/nginx/html
COPY $NGINX_CONFIG_NAME /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
