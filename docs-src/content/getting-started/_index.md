+++
title = "Server Setup"
description = ""
weight = 1
+++

This section is meant to get you started using Juzu Server via a [Docker](https://www.docker.com/products/docker-desktop) image.

## Installing the Juzu Server Image

The SDK communicates with a Juzu Server instance using [gRPC](https://grpc.io).  Cobalt distributes a docker image that contains the `juzusvr` binary and model files. The image can also include Cubic Server as well to get diarization and transcription results together.

<!--more-->

1. Contact Cobalt to get a link to the image file in AWS S3.  This link will expire in two weeks, so be sure to download the file to your own server.

2. Download with the AWS CLI if you have it, or with curl:

    ```bash
    URL="the url sent by Cobalt"
    IMAGE_NAME="name you want to give the file (should end with the same extension as the url, usually bz2)"
    curl $URL -L -o $IMAGE_NAME
    ```

3. Load the docker image

    ```bash
    docker load < $IMAGE_NAME
    ```

    This will output the name of the image (e.g. juzusvr-demo-en_us-16).

4. Start the juzu service

    ```bash
    docker run -p 2727:2727 -p 8080:8080 --name cobalt-juzu juzusvr-demo-en_us-16
    ```

    That will start a docker container with juzusvr listening for gRPC commands on port 2727 and http requests on 8080, and will stream the log to stdout. (You can replace `--name cobalt-juzu` with whatever name you want.  That just provides a way to refer back to the currently running container.)

    If the image provided also includes cubicsvr, then adding the arguments `-p CUBICSVR_GRPC_PORT:2728 -p CUBICSVR_HTTP_PORT:8081` to the `docker run` command above will enable you to access cubicsvr directly via `CUBICSVR_GRPC_PORT` and `CUBICSVR_HTTP_PORT` for gRPC and HTTP requests respectively. This is useful when you want to only obtain transcription in "true" streaming mode and diarization is not required.

5. Verify the service is running by calling

    ```bash
    curl http://localhost:8080/api/version
    ```

## Contents of the docker image

- **Base docker image** : debian-stretch-slim
- **Additional dependencies**
  - sox
  - a subset of pre-compiled kaldi binaries

### Cobalt-specific files

- **juzusvr** - binary for performing Speaker Diarization
- **diarization models** - models for extracting and scoring speaker embedding vectors

- If configured to include [Cubic](https://cobaltspeech.github.io/sdk-cubic/):
  - **cubicsvr** - binary for performing Automatic Speech Recognition
  - **cubic models** - speech recognition models and [associated files](https://cobaltspeech.github.io/sdk-cubic/getting-started/#cobalt-specific-files)
