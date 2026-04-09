.PHONY: dev prod push deploy clean

dev:
	docker compose -f docker-compose.yml up -d

prod:
	docker compose -f docker-compose.prod.yml up -d

push:
	docker build --platform linux/amd64 -t tatoshko/bloodpressure-bot:latest .
	docker push tatoshko/bloodpressure-bot:latest
deploy:
	scp .env root@useful.team:/var/bot/blood-pressure-bot/.env
	scp docker-compose.prod.yml root@useful.team:/var/bot/blood-pressure-bot/docker-compose.yml
	ssh root@useful.team 'cd /var/bot/blood-pressure-bot && docker compose pull && docker compose up -d'
clean:
	ssh root@useful.team 'cd /var/bot/blood-pressure-bot && docker compose down -v && docker rmi -f tatoshko/bloodpressure-bot:latest'