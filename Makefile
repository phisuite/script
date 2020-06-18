EDITOR?=localhost:8010
INSPECTOR?=localhost:8020
FILE?=main

inspect:
	go run ./src inspect -inspector=${INSPECTOR} ${FILE}

update:
	go run ./src update -editor=${EDITOR} -inspector=${INSPECTOR} ${FILE}
