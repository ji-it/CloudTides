FROM node:14-alpine
WORKDIR /app/

COPY package.json /app/
RUN npm install

COPY . /app/

EXPOSE 4200
CMD ["npm","run-script","start.prod"]