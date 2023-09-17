# Use an official Go runtime as a parent image
FROM golang

# Set the working directory inside the container
WORKDIR /app

# Copy the local hello.go file to the container
COPY main.go .

# Build the Go application inside the container
RUN go mod init go-module
RUN go get github.com/gin-gonic/gin
RUN go mod tidy

EXPOSE 8010

# Run the executable generated in the previous step
#CMD ["./main"]
CMD ["go", "run", "main.go"]




