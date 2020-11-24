const autocannon = require('autocannon')
data = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDYyNTA1NjAsImlhdCI6MTYwNjIzMjU2MCwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.q0S2E1ufygHeolwetgVyyHzC82z6L8c5cSJZUspscf5LC3CxKwYYXNz6MIFCasExXmV5qTB3xY6CuUu7n1c5jeIYbnI9EZJzjFvdENSLdvmKdfLE0AoNLyiQd_qBsDJ_jYK_Sce858M6tGg0Bu6EFPlHmugPyf1AQKJt56ELXe6gjx8tpDkz8edANc-OT4t3THXL6z-QsilSyu5HTO6TYcI9-xX9bwX1FpY17FbNpwn-014rLmnlMaOf8NqPE3SVmipgflaQxLIrGnqsRW65xH_ntJpVZQrXBf1etYMJfFNogrVuEvHCTiVL41s4SZ679SkB8GOcec0w024sVJcV8w",
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