FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o myGolangApp

# ENV ENV_VAR_NAME=env_var_value

EXPOSE 5000

ENTRYPOINT ./myGolangApp

# CMD ["./app"]