# Build the go application into a binary
FROM golang:alpine as builder
WORKDIR /go-github-wip
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o main .
RUN apk --update add ca-certificates

# Run the binary on an empty container
FROM scratch
COPY --from=builder /go-github-wip/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/main"]
