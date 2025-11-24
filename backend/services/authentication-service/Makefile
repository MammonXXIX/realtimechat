migrate-create:
	@migrate create -seq -ext sql -dir ./migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@migrate -path ./migrations -database "postgresql://neondb_owner:npg_JHT3cu5AGdCP@ep-wild-flower-a1nvwm29-pooler.ap-southeast-1.aws.neon.tech/authentication-service?sslmode=require&channel_binding=require" -verbose up

migrate-down:
	@migrate -path ./migrations -database "postgresql://neondb_owner:npg_JHT3cu5AGdCP@ep-wild-flower-a1nvwm29-pooler.ap-southeast-1.aws.neon.tech/authentication-service?sslmode=require&channel_binding=require" -verbose down

.PHONY: migrate-create migrate-up migrate-down

