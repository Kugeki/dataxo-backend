.PHONY: run, run-and-attach, stop

run:
	docker compose up -d --build

run-and-attach:
	docker compose up --build

stop:
	docker compose down
