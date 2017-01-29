curl -X POST -d '{ "name":"mr. Johnson"}' http://localhost:8080/v1/example/politeGreeting

echo ""

curl -X POST -d '{ "name":"Johnnie"}' http://localhost:8080/v1/example/coolGreeting

echo ""


curl -X POST -d '{ "a":"2", "b":"3", "c":"4"}' http://localhost:8080/v1/example/compute

echo ""

