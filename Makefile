create_migration:
	docker run -v $(CURDIR)/migrations:/migrations migrate/migrate create -seq -ext=.sql -dir=./migrations $(NAME)
