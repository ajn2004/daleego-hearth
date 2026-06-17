.PHONY: migrate-up migrate-down migrate-all-down migrate-reset migrate-status

migrate-up:
	goose -dir migrations postgres "$$DATABASE_URL" up
	$(MAKE) migrate-status

migrate-down:
	goose -dir migrations postgres "$$DATABASE_URL" down
	$(MAKE) migrate-status

migrate-all-down:
	goose -dir migrations postgres "$$DATABASE_URL" down-to 0
	$(MAKE) migrate-status

migrate-reset:
	$(MAKE) migrate-all-down
	$(MAKE) migrate-up

migrate-status:
	goose -dir migrations postgres "$$DATABASE_URL" status
