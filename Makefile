build:
	docker-compose build messaggio

up:
	docker-compose up -d

up_log:
	docker-compose up

stop:
	docker-compose stop

rm:
	docker-compose stop \
	&& docker-compose rm \
	&& sudo rm -rf pgdata/