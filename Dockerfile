### Build stage for the website frontend
FROM node:15 as website
RUN apt-get update
RUN apt-get install -y protobuf-compiler libprotobuf-dev
WORKDIR /code
COPY ./website/package.json ./
COPY ./website/package-lock.json ./
RUN npm ci --no-audit --prefer-offline
COPY ./protos/ ../protos/
COPY ./website/ ./
RUN npm run codegen
RUN npm run build

### Build stage for the website backend server
FROM golang:1.15.6-alpine as server
RUN apk add gcc musl-dev
RUN apk add protobuf
RUN apk add protobuf-dev
WORKDIR /code
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=1
ENV GO111MODULE=on
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./protos/ ./protos/
COPY ./codegen.sh ./
RUN ./codegen.sh
COPY ./main.go ./main.go
COPY ./internal/ ./internal/
RUN go build -o falcry

### Server
FROM alpine:3.10
COPY --from=server /code/falcry /usr/local/bin/falcry
COPY --from=website /code/build /website/build
CMD ["falcry"]
