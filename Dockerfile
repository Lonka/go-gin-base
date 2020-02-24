FROM golang:1.12.6 as go_base
ENV APP_PATH=/app
RUN mkdir -p ${APP_PATH}
WORKDIR ${APP_PATH}
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag

FROM go_base as builder
ENV APP_PATH=/app
ADD . ${APP_PATH}
WORKDIR ${APP_PATH}
RUN sed -i 's/^RunMode = debug/RunMode = release/g' conf/app.ini
RUN swag init
RUN CGO_ENABLED=0 go build .

FROM alpine
LABEL maintainer="Lonka Liu"
WORKDIR /app
COPY --from=builder /app/go_gin_base /app/go_gin_base
COPY --from=builder /app/conf /app/conf
EXPOSE 8080
ENTRYPOINT [ "./go_gin_base" ]
