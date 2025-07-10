example command:
mockgen -source=internal/services/password_service.go -destination=internal/mocks/mock_services/mock_password_service.go -package=mockservices

mockgen -source=internal/services/email_service.go -destination=internal/mocks/mock_services/mock_email_service.go -package=mockservices

mockgen -source=internal/services/token_service.go -destination=internal/mocks/mock_services/mock_token_service.go -package=mockservices

mockgen -source=internal/services/user_service.go -destination=internal/mocks/mock_services/mock_user_service.go -package=mockservices
