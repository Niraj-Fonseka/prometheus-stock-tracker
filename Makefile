

all: build push

build:
	docker build -t stock-expoter . 

run:
	docker run -d -p 2112:2112 stock-expoter

push:
	docker tag stock-expoter dockerfonseka/stock-expoter:latest