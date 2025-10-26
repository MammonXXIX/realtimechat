migrate-create:
	@migrate create -seq -ext sql -dir ./migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@migrate -path ./migrations -database "postgresql://Contact_Service:npg_02wnulKRiULf@ep-wild-river-a154rapw-pooler.ap-southeast-1.aws.neon.tech/Contact_Service_Database?sslmode=require&channel_binding=require" -verbose up

migrate-down:
	@migrate -path ./migrations -database "postgresql://Contact_Service:npg_02wnulKRiULf@ep-wild-river-a154rapw-pooler.ap-southeast-1.aws.neon.tech/Contact_Service_Database?sslmode=require&channel_binding=require" -verbose down

.PHONY: migrate-create migrate-up migrate-down
