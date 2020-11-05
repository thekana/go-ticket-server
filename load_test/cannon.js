const autocannon = require('autocannon')
data = {
    custoken1:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ1NzA5NTksImlhdCI6MTYwNDU1Mjk1OSwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.16JVkqWNG-w4hOkArUP7MYYx9Vydnijr-m3HrxObHmFPn7YbWuNdFyYM4Vn9BTkgwz1It1YLBhnq85m7NVxnMI9l4Qb7_XWzJGySWWjCHHI1VmlSLuG8DWrHl4J5slypawzfW_qGcdsPY3AFn-JhcqK1G3fD824UdiPwbqmv_YniK9zR9TCqFepst5PwKZi--VIUaakPXqAwpKL0fg8v7hws20wm8hvy9J9S8CvVHf_a-_jxkuDOSgYWFOjAYIW76cmO5QzgPBYvDqLtdfWyFrfYGso9R9IIHE-PIBHWmB8f5WkLJ2j96XJSm_fLm0qvO_TmWjRu7r-pu-P9isYKcg",
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