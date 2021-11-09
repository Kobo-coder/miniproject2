FROM golang:1.17

COPY . /src/miniproject2
WORKDIR /src/miniproject2
RUN [ "make", "compile-proto" ]
RUN [ "go", "build", "-o", "/build/miniproject2", "github.com/Kobo-coder/miniproject2"]
ENTRYPOINT [ "/build/miniproject2" ]