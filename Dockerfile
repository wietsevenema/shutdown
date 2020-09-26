# A named stage: "build"
FROM golang:1.15 AS build
 
WORKDIR /go/src/app
COPY go.* /go/src/app/
RUN go mod download
 
COPY . /go/src/app
RUN go build -o /go/bin/app

# A new stage: "run"
FROM gcr.io/distroless/base-debian10:nonroot AS run

# Copy the binary from stage build
COPY --from=build /go/bin/app /
COPY web /go/bin/web

WORKDIR /go/bin/
ENTRYPOINT ["/app"]
