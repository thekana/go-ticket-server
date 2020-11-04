const autocannon = require('autocannon')
data = {
    custoken1:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ0OTM4OTAsImlhdCI6MTYwNDQ3NTg5MCwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.Flsq7eBxdjfi6pZMBi2rSkTXSAfAolo-48Lg9tG3o2eDTxNY9sARziUMtz1kUgsqzeLIaorcRtorKxPhhdXeam6TovGo9Q_CfkbrWGdGpdg6QJHpz0XFiHxdhf0bD4vKmH3mpc110iGxlbMjGVDiW7oGMEXSGEneicGRR61HbVLultLYkzHrh2H3QDXP311RgPLCJOkOJpSXOwJ7i3oiGpbe9UetfYSoNMIOQkA-M5apDo4qfKq7fysK5TM1TrVkcWTO-v9w-LYhu2V5Od8dkpWZpe43AkdZXreSllWQpqQI04hmsI596800G6P_aS4qv8puX_xs9Sa1SBgVrZ6_xg",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    const url = 'https://www.ticketeer.ml/api/v1/reservation/reserve'
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