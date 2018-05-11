

startcb:
	docker start /db
stopcb:
	docker stop /db

test:
	go test -v -run TestGetProduct
	go test -v -run TestBoostSearch
	go test -v -run TestFilteredSearch
	go test -v -run TestCounters

