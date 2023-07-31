# Stage 1: Build --------------------------

FROM golang:latest AS build

ARG GIT_COMMIT
ENV GIT_COMMIT ${GIT_COMMIT}

ARG LOG_LEVEL
ENV LOG_LEVEL ${LOG_LEVEL}

WORKDIR /app
ADD . /app

RUN make build-for-docker

# Stage 2: Application --------------------

FROM alpine

ARG GIT_COMMIT
ENV GIT_COMMIT ${GIT_COMMIT}

ARG LOG_LEVEL
ENV LOG_LEVEL ${LOG_LEVEL}

RUN apk update \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build /app/proxy /app

EXPOSE 80 443

CMD [ "./proxy", "run" ]
