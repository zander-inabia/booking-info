FROM golang:1.12-alpine as builder
RUN apk add --no-cache ca-certificates git && \
      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
      chmod +x /go/bin/dep

ENV PROJECT purchases
WORKDIR /go/src/$PROJECT

# restore dependencies
COPY Gopkg.* ./
RUN dep ensure --vendor-only -v
COPY . .
RUN go build  -o /go/bin/purchases .

FROM alpine as release
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/purchases /purchases
ENV APP_PORT=9081
EXPOSE 9081
ENTRYPOINT ["/purchases"]
