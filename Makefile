tests:
	go test -v ./... -coverprofile coverage.out

coverage-html: tests ;
	go tool cover -html=coverage.out

coverage-console: tests ;
	go tool cover -func=coverage.out

generate-mocks:
	mockgen -package=mock -source=chat/domain/repository.go > pkg/chat/mock/domain_repository_mock.go
	mockgen -package=mock -source=chat/domain/service.go > pkg/chat/mock/domain_service_mock.go
	mockgen -package=mock -source=chat/domain/event.go > pkg/chat/mock/domain_event_mock.go
	mockgen -package=mock -source=bot/domain/service.go > pkg/bot/mock/domain_service_mock.go
	mockgen -package=mock -source=bot/domain/event.go > pkg/bot/mock/domain_event_mock.go
