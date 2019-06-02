NAME=web_driver
CSVPATH=/go/src/web_driver/walmart_reviews.csv

build:
	docker build -t $(NAME) .

run: build
	docker run --mount src="/Users/jrieck/go/src/adrian/walmart_reviews.csv",target=$(CSVPATH),type=bind --rm --name $(NAME) -d $(NAME) ./driver && docker logs -f $(NAME)

debug-run: build
	docker run --rm  -d $(NAME) sleep 300