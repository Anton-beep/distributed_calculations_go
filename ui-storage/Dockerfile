FROM node:21-alpine3.18

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH
ENV REACT_APP_STORAGE_API_URL=http://localhost:8080

COPY package.json ./
COPY package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts@3.4.1 -g --silent
RUN npm install -g --silent serve

COPY . ./

EXPOSE 3000

RUN npm run build
CMD ["serve", "-s", "build"]