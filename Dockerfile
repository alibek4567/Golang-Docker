FROM golang

RUN mkdir /app
ADD . /app/
WORKDIR /app/cmd/web/
EXPOSE 4000
RUN go build -o main .
CMD ["/app/cmd/web/main"]