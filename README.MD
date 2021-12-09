Run db container first, so that app does not connect to db while db has not started yet
1. docker-compose up -d db 
2. docker-compose up -d app

- CRUD endpoint at /user, and palindrome checker at /palindrome

get(read), delete and update needs userId sent in as header, 
update, and create need request body as {"firstName": "Vinh", "lastName": "Tran"}

/palindrome userId sent in as header


