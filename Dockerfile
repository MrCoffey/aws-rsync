FROM golang:1.13.0-stretch

ENV GO111MODULE=on
ENV GOSUMDB=off

RUN mkdir -p $GOPATH/MrCoffey/s3-resync
WORKDIR $GOPATH/MrCoffey/s3-resync
COPY . .

RUN go build -o main
RUN mkdir /app && cp main /app

ENTRYPOINT [ "./entrypoint.sh" ]