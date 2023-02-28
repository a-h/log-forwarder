FROM golang:1.20.1
WORKDIR /entrypoint
COPY main.go .
COPY go.* .
RUN go build -ldflags "-s -w" -tags lambda.norpc -o entrypoint .

FROM logstash:8.6.2
COPY logstash.conf /usr/share/logstash/pipeline/logstash.conf
USER root
RUN chown -R logstash:root /usr/share/logstash/pipeline
USER logstash
COPY --from=0 /entrypoint/entrypoint .
CMD ["./entrypoint"]
