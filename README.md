# image-promote

Image-promote plugin to pull Docker images and publish them to another container registry.

## Build

Build the binary with the following commands:

```sh
sh .build.sh
```

## Docker

Build the Docker image with the following commands:

```sh
docker build --rm=true -t plugins/image-promote .
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_TAG=TAG \
  -e PLUGIN_PUSH_REPO=IMAGE_TO_PUSH \
  -e PLUGIN_PUSH_REGISTRY=REGISTRY1 \
  -e PLUGIN_PUSH_USERNAME=USERNAME1 \
  -e PLUGIN_PUSH_PASSWORD=PASSWORD1 \
  -e PLUGIN_PULL_REPO=IMAGE_TO_PULL \
  -e PLUGIN_PULL_REGISTRY=REGISTRY2 \
  -e PLUGIN_PULL_USERNAME=USERNAME2 \
  -e PLUGIN_PULL_PASSWORD=PASSWORD2 \
  -e PLUGIN_INSECURE_REGISTRIES=REGISTRY1,REGISTRY2 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  --privileged \
  plugins/image-promote
```
