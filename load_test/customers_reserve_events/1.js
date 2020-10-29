const data = require('./data.js')

data.quotaLoad([
    {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": data.eventID4,
            "amount": 3
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": data.eventID2,
            "amount": 4
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": data.eventID3,
            "amount": 5
        }),
    },
])
