const autocannon = require('autocannon')
data = {
    custoken1:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ1OTc5NjMsImlhdCI6MTYwNDU3OTk2MywibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.gM4nXrQER_okCCO53-3Xnis7RUnBKdDtnNmvbPZGt_g2fFQHcESKQ7FnpuTFCjiyaT8FVtl7FZkHCueN6S01RtpGtotBmQzvowniU3_TIrskgkLR1wENZSsh_GLGbot8hNAyUy50iGHb4aVsfP15Pqml_bVfa2WsxN6HdodaEQPgHQNYx21BKbVtM5ZwBcCZ3Ds5kl-znsTHlTK3pwrcRi9BNm6JR8O45gK9mrVbTxz7JP8YeswgUj1iLPbD41uqhjGPKktFb-0Vu4tCl2noV2_Ovpae6XV-2Fox3WSF0VMTUxVEhyxQ6rRtBMmU3_T7ErfGLh-Lw8sUNk0aQyZhhA",
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