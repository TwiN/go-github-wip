# Build the go application into a binary
FROM golang:alpine as builder
WORKDIR /go-github-wip
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o main .

# Run the binary on an empty container
FROM scratch
COPY --from=builder /go-github-wip/main .
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/main"]
