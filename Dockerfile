FROM golang:1.18-bullseye

ENV TERRAFORM_VERSION=0.14.0

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list \
    && apt update \
    && apt install goreleaser -y

RUN apt install unzip -y

RUN wget -nv https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/local/bin

WORKDIR /root/terraform-provider-pfptmeta
