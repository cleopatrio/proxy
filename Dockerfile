# Stage 1: Build --------------------------

FROM golang:alpine AS build

WORKDIR /app
ADD . /app

RUN cd /app && go build -o ./server

# Stage 2: Application --------------------

FROM alpine

RUN apk update \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build /app/server /app

EXPOSE 80 443

CMD [ "./server", "run" ]
