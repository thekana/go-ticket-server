const data = require('./data.js')

data.quotaLoad([
    {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken3,
            "eventID": data.eventID1,
            "amount": 3
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken3,
            "eventID": data.eventID2,
            "amount": 2
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken3,
            "eventID": data.eventID3,
            "amount": 7
        }),
    },
])
