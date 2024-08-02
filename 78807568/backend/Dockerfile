FROM golang
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o bin
CMD [ "/app/bin" ]