redis-sv:
	@redis-server

load_blobs:
	@go run . -l

flush:
	@redis-cli flushall

.PHONY:
	redis-sv, flush
