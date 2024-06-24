ARG GO_VERSION=1.22.3
ARG GO_PORT=7041

FROM golang:${GO_VERSION}-alpine AS builder
# FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./svclookup ./main.go 

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/svclookup .
# EXPOSE 8080
EXPOSE ${GO_PORT}
# ENTRYPOINT ["./svclookup"]
# RUN PWD
ENTRYPOINT ["./app/svclookup"]