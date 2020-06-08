FROM golang:1.14 as builder

WORKDIR /workspace
COPY . /workspace

RUN make build-linux

FROM alpine:3

RUN apk update && apk add curl

COPY --from=builder /workspace/build/linux-amd64/gosea /usr/local/bin/gosea

RUN addgroup -S gosea && adduser -S gosea -G gosea
USER gosea

HEALTHCHECK CMD curl --insecure -f https://localhost:8000/health || exit 1;

EXPOSE 8000

ENTRYPOINT ["ls"]

CMD [ "-h" ]