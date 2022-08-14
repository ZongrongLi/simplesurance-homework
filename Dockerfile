FROM golang:1.17-alpine as build

WORKDIR /go/src/
COPY . .

RUN go mod vendor &&  CGO_ENABLED=0 go build -mod=vendor -o /counter main.go routes.go schema.go utils.go

FROM xiaochengtech/alpine-timezone

COPY --from=build /counter /counter

RUN addgroup go \
    && adduser -D -G go go \counter \
    && chown -R go:go /counter \
    && mkdir data \
    && touch data/data.csv


ENV ENV dev
EXPOSE 8080

ENTRYPOINT ["/counter"]