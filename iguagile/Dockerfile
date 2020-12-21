FROM golang:latest AS build

COPY . $GOPATH/src/engine

RUN cd $GOPATH/src/engine/ && \
    CGO_ENABLED=0 go build -o /app


FROM scratch

COPY --from=build /app /app

CMD ["/app"]
