test:
	go test ./utils/ -count=1
	go test ./api_client/ -count=1
	go test ./ws_client/ -count=1

test_v:
	go test ./utils/ -count=1 -v
	go test ./api_client/ -count=1 -v 
	go test ./ws_client/ -count=1 -v