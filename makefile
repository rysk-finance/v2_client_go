test:
	go test ./utils/ -count=1
	go test ./api_client/ -count=1 
	go test ./ws_client/ -count=1 

test_utils:
	go test ./utils/ -count=1 -cover

test_api_client:
	go test ./api_client/ -count=1 -cover

test_ws_client:
	go test ./ws_client/ -count=1 -cover

test_unit: 
	go test --tags=unit ./utils/ -count=1 -cover
	go test --tags=unit ./api_client/ -count=1  -cover
	go test --tags=unit ./ws_client/ -count=1  -cover

test_integration: 
	go test --tags=integration ./utils/ -count=1 -cover
	go test --tags=integration ./api_client/ -count=1  -cover
	go test --tags=integration ./ws_client/ -count=1  -cover

coverage:
	go test ./utils/ -count=1 -coverprofile=utils_coverage.out
	go tool cover -func=utils_coverage.out
	go test ./api_client/ -count=1 -coverprofile=api_client_coverage.out
	go tool cover -func=api_client_coverage.out
	go test ./ws_client/ -count=1 -coverprofile=ws_client_coverage.out
	go tool cover -func=ws_client_coverage.out
