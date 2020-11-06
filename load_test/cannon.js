const autocannon = require('autocannon')
data = {
    custoken1:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ2OTExMjYsImlhdCI6MTYwNDY3MzEyNiwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.sddzRZCeg1o9RvDOwHBGmn5TpElfE5GFzM6u_yyhbfsE0nJHhkr4lmDiQ94mgpI_HdgEoEo6wCcxRWpyO12tOicrXQ0bENnivw6kOMGuOxGegOL5tZ4iFq3D7q5x949as2CfDMe1croH2uF2ACxRXFILRPXbal2BY_R1iB6beNibS-3Q98z7TrP4SKZsKg7s78zOkPeylp3kUf-9jlXzzfplHCvPKS3B5hL1RTr-mNa3lHwkCpEIMsAb83JtQ4ZYvW1sz5SelH6iiZQ0xJ_sBMqhkN0aynDDYxnBbjrAKyZ5-wBrD5QOs8SN3-ZV_hxV-RSUNljDi-JxXREdZWEwaQ",
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