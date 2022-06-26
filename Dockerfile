# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env ./
COPY cert/*.* ./cert/
RUN go mod download

COPY *.go ./
COPY db/*go ./db/
COPY db/memo/*.go ./db/memo/
COPY db/tp/*.go ./db/tp/

COPY middleware/*.go ./middleware/
COPY routes/*.go ./routes/
COPY model/user/*.go ./model/user/
COPY src/*.go ./src/

RUN go build -o /mend

EXPOSE 8010
CMD [ "/mend" ]