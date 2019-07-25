# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/yamamushi/durouter

# Create our shared volume
RUN mkdir /durouter

# Get the du-discordbot dependencies inside the container.
RUN go get gopkg.in/mgo.v2
#RUN cd /go/src/github.com/anaskhan96/soup && git checkout ad448eafe
RUN go get github.com/BurntSushi/toml
RUN go get github.com/go-chi/chi
RUN go get github.com/go-chi/docgen
RUN go get github.com/go-chi/render

# Install and run du-discordbot
RUN go install github.com/yamamushi/durouter

# Run the outyet command by default when the container starts.
WORKDIR /durouter
ENTRYPOINT /go/bin/durouter

VOLUME /durouter