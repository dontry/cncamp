# Build stage
FROM golang:1.19-alpine AS build

# Set the working directory inside the container
WORKDIR /bin

# Copy the source from the current directory to the working Directory inside the container
COPY go.mod ./

# Download all the dependencies
RUN go mod download

# Copy the source code from the current directory to the working Directory inside the container
COPY metrics ./metrics
COPY *.go ./
COPY go.sum ./
COPY go.mod ./


# Build the Go app
# Turn off CGO since that can result in dynamic links to libc/libmusl.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o httpserver .

# Final stage
FROM scratch

# Copy the Pre-built binary file from the previous stage to /app
COPY --from=build /bin/httpserver /bin/httpserver


ENV VERSION=0.6.0
ENV PORT=8080

EXPOSE ${PORT}

ENTRYPOINT [ "/bin/httpserver" ]
CMD [ "-verbose"]
