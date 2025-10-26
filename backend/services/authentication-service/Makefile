migrate-create:
	@migrate create -seq -ext sql -dir ./migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@migrate -path ./migrations -database "postgresql://User_Service:npg_iHUbkJ1sDB6L@ep-wild-river-a154rapw-pooler.ap-southeast-1.aws.neon.tech/User_Service_Database?sslmode=require&channel_binding=require" -verbose up

migrate-down:
	@migrate -path ./migrations -database "postgresql://User_Service:npg_iHUbkJ1sDB6L@ep-wild-river-a154rapw-pooler.ap-southeast-1.aws.neon.tech/User_Service_Database?sslmode=require&channel_binding=require" -verbose down

.PHONY: migrate-create migrate-up migrate-down

