const autocannon = require('autocannon')
data = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDU0Njg2ODQsImlhdCI6MTYwNTQ1MDY4NCwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.Yu1LHqj65uhcku6UUg6At3N3mb4KUGh1sSa2Qpli5ykv54aHXCm9f6kYollPUfHqpqXQHKTJluXgt0LClOXJPBA1vL1zjEC1GlDsef48WekukwaQhov-8yGoQokc5wUw_cWNWYDE04DWBgwlXc6G1dPiwfcxgIxRgS7vM_Z38hK-b3OSfSiJxovSJhSyqShjOxMQtcT_yKa9SmC3E3021TPBC8odCacFHcAPJZagOZCNLYjCjeeT3atibf3oXPWPsiUrSEjXSz4FbW7wByC07B_X-Q1cWyMMk-xsieRLO6SwfpQq7j24GvoAb-HnJ0ZShdhVoN3VgTb-PVMv0_G0WA",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    //const url = 'http://localhost:9092/api/v1/reservation/reserve'
    const url = 'https://www.ticketeer.ml/api/v1/reservation/reserve'
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