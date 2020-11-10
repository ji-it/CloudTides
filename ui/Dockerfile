FROM node:latest
WORKDIR /app/

COPY package.json /app/
RUN npm install

COPY . /app/
EXPOSE 4200

CMD ng serve