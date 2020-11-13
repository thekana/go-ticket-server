const autocannon = require('autocannon')
data = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDUyNzQ4NjUsImlhdCI6MTYwNTI1Njg2NSwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.rvolSiccDfC4zZgFYZcrNy-iM1lR8bO_KDw4EwsadONd1lwk6aclMele7ibPDq7Z49agKeVlJiH8QJkSMWtfewLGba_YQ2JuUKp4S9rqQdufp8irEbm-Xo49qztkQnGYptQuAzKouKgRGWyxQvmrwGd8KQuFDPaoJvA6rDoc1pUZsU8YjKiK-DPvT9kgeXCbrcJWLnbNHw74xGB1GUX3uJ1y7bNA8kH6ExMFVctd9owHfEHtdVU2zxNKvrubzl_OvCozwAwq2wV1GCVDJB2I1Pzqo54MPZjX_-iJ_THhf_sAw8DfNc1IYIaFS3-Hi7AlVrdMfhQ5UOg3xaaeuEHSQA",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    const url = 'http://localhost:9092/api/v1/reservation/reserve'
    //const url = 'https://www.ticketeer.ml/api/v1/reservation/reserve'
    autocannon({
        url: url,
        connections: 80, // matters a lot can't be more than a hundred and less than cut off amount
        duration: 10,
        headers: {
            // by default we add an auth token to all requests
            'Authorization': "Bearer " + data.custoken1
        },
        requests: [
            {
                method: 'POST',
                body: JSON.stringify({
                    "eventID": data.eventID1,
                    "amount": 1
                }),
            }, {
                method: 'POST',
                body: JSON.stringify({
                    "eventID": data.eventID2,
                    "amount": 1
                }),
            }, {
                method: 'POST',
                body: JSON.stringify({
                    "eventID": data.eventID3,
                    "amount": 1
                }),
            }, {
                method: 'POST',
                body: JSON.stringify({
                    "eventID": data.eventID4,
                    "amount": 1
                }),
            },
        ],
        excludeErrorStats: true
    }, (err, res) => {
        console.log('finished bench', err, res)
    })
}

start()