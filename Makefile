postgres:
	docker run --name postgres -e POSTGRES_USER=root -d postgres

build:
	docker build --no-cache=true -t izqui/botfounder .

run:
	docker run -e FOUNDERBOT_TOKEN=$(FOUNDERBOT_TOKEN) -P --link postgres:postgres  izqui/botfounder &