TARGET_DIR=target

COMMIT_ID:=$(shell git rev-parse HEAD)
BRANCH=$(shell git branch | grep \* | sed 's/ /\n/g' | head -2 | tail -1)

# collect packages and dependencies for later usage
PACKAGES=$(shell go list ./... | grep -v /vendor/)

# choose the environment, if BUILD_URL environment variable is available then we are on ci (jenkins)
ifdef BUILD_URL
ENVIRONMENT=ci
else
ENVIRONMENT=local
endif