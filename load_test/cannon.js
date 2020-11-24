const autocannon = require('autocannon')
data = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDYyMTk4NjcsImlhdCI6MTYwNjIwMTg2NywibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.dmBzjrRPVAZ6QEMJc8gia9-SGMgeJm9-mkjpbjxy481Ch8xtIeaQQfKPkTYjzuWq5kbDPL2tTQ2Vn-JTdeR9IG9u8jxnkJGcLaaD5zOqgVkdmDIPSLpG7pmbzzUubQO3rj6lfqz9yLate5EhL8dXR7YMKc4esWXwW-AEfVZOZ37m9-jic_V-yQGVj1954w0jty83qRdRGGdhJeJq6NYZVXqoEEOHgFr4QCJCy21TCl0ZbeaZLMpuywHtFxVaEs0kUTIHpts-4yNEFEdKpXLnzIqB7qoI-o74rKIOhkweQaw3uxnjSf1ZDnnh_4c7-v4WrKUu6uWWqWYkWRhX2VYb6w",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
}

function start() {
    //const url = 'http://localhost:9092/api/v1/reservation/reserve'
    //const url = 'https://www.ticketeer.ml/api/v1/reservation/reserve'
    const url = 'https://www.goticket.tk/api/v1/reservation/reserve'
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