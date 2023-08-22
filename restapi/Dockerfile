# use latest 1.x golang debian buster image for the building stage
FROM golang:1.21 AS builder

WORKDIR /usr/local/src/restapi

# disable CGO
ENV CGO_ENABLED=0

# copy only module information to take advantage of cache and layers
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy source files
COPY . ./

# build executable
RUN go build -o restapi

# make a fresh start to final image
FROM alpine:3.15.4

# copy executable
COPY --from=builder /usr/local/src/restapi/restapi /usr/local/bin/restapi

# use the executable as the main program for image
ENTRYPOINT [ "/usr/local/bin/restapi" ]
