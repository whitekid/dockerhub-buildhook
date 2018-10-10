ARG TARGET

FROM golang as build

WORKDIR /usr/app

COPY go.mod go.sum ./
RUN go get -v

COPY . ./
RUN make

FROM whitekid/debian:runtime

WORKDIR /usr/app
COPY --from=build /usr/app/bin/dockerhub-buildhook /usr/app/dockerhub-buildhook
