# Copyright 2020
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: push

#
# Docker tag with v prefix to differentiate the official release build, triggered by git tagging
# this is pushed to datastax Dockerhub repo
#

IMAGE_REPO_BASE ?= itestmycode
GIT_COMMIT = $(shell git rev-list -1 HEAD)
TAG ?= $(GIT_COMMIT)

prom2dd:
    go build -o prom2dd -tags musl src/main.go

build: prom2dd

test:
	go test ./...

container:
	docker build -t $(IMAGE_REPO_BASE)/prom2dd:$(TAG) .
	docker tag $(IMAGE_REPO_BASE)/prom2dd:$(TAG) $(IMAGE_REPO_BASE)/prom2dd:latest 

push: container
	docker push $(IMAGE_REPO_BASE)/prom2dd:$(TAG)
	docker push $(IMAGE_REPO_BASE)/prom2dd:latest

clean:
	go clean --cache
	docker rmi $(IMAGE_REPO_BASE)/prom2dd:$(TAG)
	docker rmi $(IMAGE_REPO_BASE)/prom2dd:latest

static-check: lint

lint:
	golangci-lint run

