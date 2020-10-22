'use strict'

const eventID = "ccccc97f-b357-4815-935d-94ff5b2680da"
const token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMzODEwMzEsImlhdCI6MTYwMzM3NzQzMSwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.HlqGQZeuElslL4isQjSj4qkU1JS0PLDQsdFvprnQ5ftn6GWjpPakSYO-McojINWbqDzpwRDJ9ZlC1l8ND1_r3ZPxLR3uZX4kOtmDrioje5vcIUp4oiLo8jNjauDG0KVML4l0_bhf9vf72yjRvhEAHfp0Wo7bYp064jaj0Li_I-ywGLhdJUuggaYHbvvh0MDk8GaJlh6ISJHp5Rp5QfSfYM-tiQXqjn6xswMCC9Q5H73BLIYSg2y23O4bUvwaZy_GTpfxwc06zAqn6TICTvt7W1C80l3vKXnNHingGNoDAtPNmCE3t4BwbmHC5Df1BLs1i2fDLp9XG-NqySD5DfsfHw"

const autocannon = require('autocannon')

function startBench() {
    const url = 'http://localhost:9092/api/v1/reservation/reserve'

    autocannon({
        url: url,
        connections: 1000,
        duration: 10,
        requests: [
            {
                method: 'POST', // this should be a post for logging in
                body: JSON.stringify({
                    "authToken": token,
                    "eventID": eventID,
                    "amount": 4
                }),
            },
        ]
    }, finishedBench)

    function finishedBench(err, res) {
        console.log('finished bench', err, res)
    }
}
startBench()