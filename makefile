generate_doc:
	swag init --parseDependency --parseInternal
start:
	docker-compose up -d  