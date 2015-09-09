postgres:
	docker run --restart=on-failure:10 --name postgres -e POSTGRES_USER=root -d postgres

build:
	docker build --no-cache=true -t izqui/botfounder .

run:
	docker run --detach=true --restart=on-failure:10 -e FOUNDERBOT_TOKEN=$(FOUNDERBOT_TOKEN) -e FOUNDERBOT_URL=$(FOUNDERBOT_URL) -p 3000:3000 --name botcast --link postgres:postgres izqui/botfounder