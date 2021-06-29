FROM node:14-alpine
WORKDIR /app/
ARG MOD

COPY package.json /app/
RUN npm install
COPY . .
RUN npm run-script build.${MOD}

FROM nginx:1.19-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=0 /app/dist/tides-ui /usr/share/nginx/html
EXPOSE 80