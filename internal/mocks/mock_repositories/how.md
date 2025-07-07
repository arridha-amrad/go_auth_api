mockgen -source=internal/repositories/redis_repository.go -destination=internal/mocks/mock_repositories/mock_redis_repositories.go -package=
mockrepositories

mockgen -source=internal/repositories/user_repository.go -destination=internal/mocks/mock_repositories/mock_user_repositories.go -package=mo
ckrepositories