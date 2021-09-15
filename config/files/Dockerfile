FROM python:3.8-slim-buster as build

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV http_proxy=http://proxy.ebiz.verizon.com:80
ENV https_proxy=http://proxy.ebiz.verizon.com:80
ENV HTTP_PROXY=http://proxy.ebiz.verizon.com:80
ENV HTTPS_PROXY=http://proxy.ebiz.verizon.com:80
ENV NO_PROXY="localhost.localdomain,$(printf '%s,' .{0..255}),cluster.local,verizon.com"
ENV no_proxy="localhost.localdomain,$(printf '%s,' .{0..255}),cluster.local,verizon.com"

COPY requirements.txt .
RUN pip install -r requirements.txt

WORKDIR /app
COPY . /app

ENTRYPOINT []