FROM node:14-alpine
WORKDIR /app/

COPY package.json /app/
RUN npm install
COPY . /app/
RUN npm run-script build

FROM node:14-alpine
WORKDIR /app/
COPY --from=0 /app/dist ./dist
COPY package* ./
RUN npm install
EXPOSE 4200
CMD ["npm","run-script","start.prod"]