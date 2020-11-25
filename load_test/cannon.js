const autocannon = require('autocannon')
data = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDYzMDY4NzgsImlhdCI6MTYwNjI4ODg3OCwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.kQUWO3WrAuUuZCwROyaPzetXytlIC1AjxIUrhbxLpLmIXdMTxGP0bU5TtzHKRm-7xJ1LX-rrFuLjvTDhaWpLs2y8STnpMS4Uu7hPY_gkjS_y5UWmsB_ab9xpXau2Gr1jHyXqV-Yb0I71FbbQ_kASi1GQPcGFkExGn_7aT-4eu5aZE-xY7tx45Bx27A_JhCs_olaWr4nwY-XptgoDl2s3NFccFgPj5itvbOK-XE8Q4TG2DPm5YYxZlYWaN7-T2Xs3Wp97PMJL_BhxSFGeHf12LkQtj9pJFFXF_PxwggsHaF4WSrRD7bbQn0908D5ImSMQCIgrZu02HirVTbANKJ9KnQ",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    const url = 'http://localhost:9092/api/v1/reservation/reserve'
    //const url = 'https://www.ticketeer.ml/api/v1/reservation/reserve'
    //const url = 'https://www.goticket.tk/api/v1/reservation/reserve'
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