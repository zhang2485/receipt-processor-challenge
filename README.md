To start the server, in the `receipt-processor-challenge` folder, 

First run these commands:
```
go mod init main
go mod tidy
```

Once the above commands have completed, run 
```
go run .
```

In a separate terminal, run this to perform a request with valid receipt metadata
```
curl -i -X POST -H "Content-Type: application/json" -d "{  \"retailer\": \"M&M Corner Market\",  \"purchaseDate\": \"2022-03-20\",  \"purchaseTime\": \"04:33\",  \"items\": [    {      \"shortDescription\": \"Gatorade\",      \"price\": \"2.25\"    },{      \"shortDescription\": \"Gatorade\",      \"price\": \"2.25\"    },{      \"shortDescription\": \"Gatorade\",      \"price\": \"2.25\"    },{      \"shortDescription\": \"Gatorade\",      \"price\": \"2.25\"    }  ],  \"total\": \"9.00\"}" http://localhost:8080/receipts/process
```
It should return a uuid.


Finally, to get the score of the previous request, run a command like this.
```
curl http://localhost:8080/receipts/{uuid}/points
```
The `uuid` used in this request is the one returned from the earlier request.
