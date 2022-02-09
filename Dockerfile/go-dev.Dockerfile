FROM golang:1.17.5-alpine3.15
RUN mkdir /app
ADD ../server /app/
WORKDIR /app
RUN go get -v github.com/cosmtrek/air
ENTRYPOINT ["air"]