
for i in {1..100}
do
curl -X PUT \
  https://lk8z3zvs06.execute-api.us-west-2.amazonaws.com/dev/signal \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8f23b85a-59c6-4262-86f9-80549c4064a9' \
  -d '{
	"user_id": "test user",
	"doc_id": "12345",
	"cat_id": "catabc",
	"query": "ipad 64 gb"
}'


curl -X PUT \
  https://lk8z3zvs06.execute-api.us-west-2.amazonaws.com/dev/signal \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8f23b85a-59c6-4262-86f9-80549c4064a9' \
  -d '{
	"user_id": "test user",
	"doc_id": "23456",
	"cat_id": "catabc",
	"query": "ipad 64 gb"
}'

done

