const data = require('./data.js')

data.quotaLoad([
    {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": 7,
            "amount": 3
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": 2,
            "amount": 4
        }),
    }, {
        method: 'POST',
        body: JSON.stringify({
            "authToken": data.custoken1,
            "eventID": 4,
            "amount": 5
        }),
    },
])
