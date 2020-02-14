FROM golang:1.13-alpine
EXPOSE 34320
EXPOSE 34320/udp
EXPOSE 34321
EXPOSE 34321/udp
RUN apk add --no-cache bash

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY server-go /bin/server-go
