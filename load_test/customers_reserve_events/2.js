const data = require('./data.js')

data.quotaLoad([
    {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken2,
            "eventID": data.eventID2,
            "amount": 3
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken2,
            "eventID": data.eventID3,
            "amount": 4
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken2,
            "eventID": data.eventID4,
            "amount": 5
        }),
    },
])
