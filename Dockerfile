FROM alpine:latest

RUN apk add --no-cache --update ca-certificates
RUN apk add --no-cache python3 python3-dev && \
    python3 -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip3 install --upgrade pip setuptools && \
    if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
    if [[ ! -e /usr/bin/python ]]; then ln -sf /usr/bin/python3 /usr/bin/python; fi && \
    rm -r /root/.cache

RUN apk add --no-cache build-base git
RUN apk add --no-cache go

RUN go get gopkg.in/olahol/melody.v1
RUN pip3 install websocket-client
ADD main.go .
CMD go run main.go