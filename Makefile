default:
	@echo "=============building Local API============="
	docker build .

up: default
	@echo "=============starting api locally============="
	docker-compose up -d

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	docker-compose run web go test -v -cover ./...
	
clean: down
	@echo "=============cleaning up============="
	docker system prune -f
	docker volume prune -f