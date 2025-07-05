run:
	$(MAKE) back & $(MAKE) front
	wait

back:
	cd backend && go run ./src/*.go

front:
	cd frontend && yarn dev

prepare_frontend:
	cd frontend && yarn install