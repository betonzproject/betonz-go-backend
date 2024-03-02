FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY main.go ./
COPY internal/ ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /betonz-go

FROM alpine:3.19

# https://stackoverflow.com/questions/60103251/time-loadlocation-works-regularly-but-throws-an-error-on-my-docker-instance-how
RUN apk --no-cache add tzdata

COPY template.html .
COPY --from=builder betonz-go .

EXPOSE 8080

CMD [ "/betonz-go" ]
