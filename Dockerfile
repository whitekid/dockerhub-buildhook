ARG TARGET

FROM golang as build

WORKDIR /usr/app

COPY go.mod go.sum ./
RUN go get -v

COPY . ./
RUN make

FROM whitekid/debian:curl

WORKDIR /usr/app
COPY --from=builder dockerhub-buildhook /usr/app/dockerhub-buildhook
