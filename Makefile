setup:
	go get -u github.com/pressly/goose/cmd/goose

migrate:
	@goose -dir ./migrations mysql "root:root@/project?parseTime=true" status