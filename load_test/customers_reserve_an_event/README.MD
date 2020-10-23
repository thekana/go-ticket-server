## Multiple customers trying to book one event
#### Purpose - check if lock is working correctly when different people are trying to access the same resource
Paste valid tokens in data.js <br>
Paste valid event id in data.js <br>
```shell script
node 1.js & node 2.js & node 3.js &
```
Then GET `localhost:9092/api/v1/print` and check console to see if the numbers add up to that event quota