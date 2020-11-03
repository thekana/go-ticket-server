const autocannon = require('autocannon')
data = {
    custoken1:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ0MjY5OTIsImlhdCI6MTYwNDQwODk5MiwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.toTXgHyns0nxnHHVNw2Et5-0Wgup7ZAMeWAoof_P4PdPcgPCiiQk4q_Eyjk-2xEmXpX_38EMB9SBls1DEqPlSKpueFOWWAVgdoqRMYtwE95YUCd-ab4CsShdAhhBKO_yVFtZM7mYyA_QbPKnYQUYv7-VfyE9ZhKBn0iqAMmSuiUdupbSsDv-u5dEa3YKPYD2Vx92gumcezOug8gWSwQ7-bVvNIQu2VJYnOvf2GHLiHMM3zAU9URDh_mCYo5JabcDupfZYvffEkiR_uaerNcs9EuBwdRFnWmPst_s5izJ7TNCFBLcFVX_PMx1V9y0io8Q2NfbOEYijn_USncfMA4SHA",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    const url = '54.251.152.89/api/v1/reservation/reserve'
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