# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM golang:1.13

# Install protoc
RUN apt-get update && apt-get upgrade -y
RUN apt-get install protobuf-compiler -y

# Copy the source from the current directory the working directory, excluding
# the deployment directory.
WORKDIR /mixer
COPY . .
RUN rm -r deployment

# Install protobuf go plugin
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.23.0
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v0.0.0-20200824180931-410880dd7d91

# Only download the two files. Can `git clone` entire library if needed.
RUN mkdir -p /mixer/proto/google/api/
RUN curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto \
         --output /mixer/proto/google/api/annotations.proto
RUN curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto \
         --output /mixer/proto/google/api/http.proto
RUN protoc \
    --proto_path=proto \
    --go_out=. \
    --go-grpc_out=. \
    --go-grpc_opt=requireUnimplementedServers=false \
    proto/*.proto

# Install the Go app.
RUN go install .

# Adding the grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.5 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

ENTRYPOINT ["/go/bin/mixer"]