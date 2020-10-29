const autocannon = require('autocannon')
data = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQwMDE5OTMsImlhdCI6MTYwMzk4Mzk5MywibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.EkTgUPSRR3vl0nxhXqTRmO3aTLktbf2WMdu2CzO6nlPa62D8oqRrRzoQ_1sWdMSzb3WTWQlm2FAoVqv3LADmPMCZAiLo31b3z9T1vgNkgM1GAcoacctDXFNn4Q_Q-IqJvBmeghNuhjvXM_NHdLsS3QH5KS65UW9glo4OhIBeeXaMMxVK4XrvdGgcYHig_CxpgkzCXUo-1tIbG_eqRU1j9inF0VjaSLPor8XQjszMlYgHA25oUhj3ByYYD02qkDf7pPtcVnTLx9-rbARUcJsdi7tiJ3IIHWfbwxB1LQh2rliaAiU0hycjnjL61IYfcUTp4w-WAIulJPHwJb-suRZqzA",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    const url = 'http://localhost:9092/api/v1/reservation/reserve'
    autocannon({
        url: url,
        connections: 80, // matters a lot can't be more than a hundred and less than cut off amount
        duration: 10,
        requests: [
            {
                method: 'POST',
                body: JSON.stringify({
                    "authToken": data.custoken1,
                    "eventID": data.eventID1,
                    "amount": 1
                }),
            }, {
                method: 'POST',
                body: JSON.stringify({
                    "authToken": data.custoken1,
                    "eventID": data.eventID2,
                    "amount": 1
                }),
            }, {
                method: 'POST',
                body: JSON.stringify({
                    "authToken": data.custoken1,
                    "eventID": data.eventID3,
                    "amount": 1
                }),
            }, {
                method: 'POST',
                body: JSON.stringify({
                    "authToken": data.custoken1,
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