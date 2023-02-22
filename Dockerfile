FROM alpine:3.17.1

RUN apk add --no-cache caddy wireguard-tools

COPY ./cmd/waygate/waygate /waygate

EXPOSE 9001

ENTRYPOINT ["/waygate"]

#CMD ["client"]
