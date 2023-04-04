FROM golang:alpine

WORKDIR /app
ADD . .

RUN go get ./...

RUN GOOS=linux GOARCH=amd64 go build -o looker .

ENV IFUNNY_USER_AGENT "IFUNNY_USER_AGENT='iFunny/8.19.11(22222) iphone/16.2 (Apple; iPhone14,5)'"
ENV IFUNNY_ADMIN_ID ""
ENV IFUNNY_BEARER ""
ENV IFUNNY_STATS_CONNECTION ""
ENTRYPOINT [ "./looker" ]
