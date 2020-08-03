FROM ubuntu:18.04

RUN true && \
    apt-get update && \
    # apt-get install -y apt-utils && \
    # DEBIAN_FRONTEND=noninteractive apt-get install -y sudo && \
    apt-get install -y sudo && \
true

# https://habr.com/ru/post/448480/
ARG USER=some
RUN true && \
    groupadd --gid 2000 $USER && \
    useradd --uid 2000 --gid $USER --shell /bin/bash --create-home $USER && \
    usermod -aG sudo $USER && \
    echo "$USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers && \
true

USER $USER
WORKDIR /home/$USER

RUN true && \
    sudo apt install -y curl less && \
true

