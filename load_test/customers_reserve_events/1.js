const data = require('./data.js')

data.quotaLoad([
    {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": data.eventID4,
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
    },
])
