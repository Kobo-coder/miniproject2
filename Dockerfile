FROM golang:1.17

COPY go.mod go.sum /
WORKDIR /
RUN [ "go", "mod", "download" ]
COPY . /src/miniproject2
WORKDIR /src/miniproject2
RUN [ "go", "build", "-o", "/build/miniproject2", "github.com/Kobo-coder/miniproject2"]
EXPOSE 50000
ENTRYPOINT [ "/build/miniproject2" ]