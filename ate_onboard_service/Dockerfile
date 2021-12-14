FROM golang:1.16.5-alpine as builder

RUN mkdir /db-build
ADD . /db-build/
WORKDIR /db-build
RUN go build -o ate_onboard_service

FROM alpine
RUN mkdir /db
WORKDIR /db
COPY --from=builder /db-build/ /db/
EXPOSE 8080
CMD ["./ate_onboard_service"]