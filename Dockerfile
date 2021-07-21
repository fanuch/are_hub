# alpine distro
FROM alpine

# install go to compile against libraries
RUN apk add go

# create and change to application directory
WORKDIR /app/

# copy go.mod first and download dependencies
COPY go.mod ./
RUN go mod download

# copy source files
COPY . ./

# move deployment configuration to build/execution directory
RUN mv config.json.docker cmd/are_hub/config.json

# change to build directory and compile
WORKDIR ./cmd/are_hub/
RUN go build

# start the application
# remember to expose the port set in config.json!
CMD ["are_hub", "--config", "./config.json"]
