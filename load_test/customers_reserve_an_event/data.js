const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM0NDg2MTEsImlhdCI6MTYwMzQ0NTAxMSwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.OCF4a6g65ljHSPGV3l2ARH0D7WrBWvZ6muFNFSzULOkcdAhQ2SDyLjcks1pUhC9Uf6rJyVU3IcC2J9qlJukpru8QgrqRr71b0wwNKP5CVUYeNHDmRpq8IZ_PD1LH2wVzmBwBYPwSiL2ZSESoNVsOMg-nyzW5_OE0Cu4Olux01PPn80tjY7HccgyBimAJbs3E4u-DFzxcXPhkjsXzoQ1Btc3Vawduzupm81vbIFW1QGSHsc8zrlz4Pozib08hBE_NA29DYpXw9dIXO7IbbQOVie-K3qSLRGo-jZg0-DqgY4L78RsPw8o3J0N8g8PoefqKdnfCqFBENYmqKhkPEqKeIA",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM0NDg2MTQsImlhdCI6MTYwMzQ0NTAxNCwibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.hogFsOfiu2Ja3Wzje8nnQmC1QVzTeKS3ElYIB4wu1gHrXdciXMVpkjkXgjj9oAKCSPUFnIucT-WcK08xWSn46QLZ1QEzhFlod2VgOw3FWtd2_mOkBjo6-96oMBN5hC4-uuyXQ-HtgUY2NxA95Nhbp-bX_ZKQMw64EOnAO1x-SGyAJ9uape91EjxdCQqh_qnoJou7pVF9xJGKbjaRT6bTM1AMgLRtSbgiRnkd6AtcdTNILaN9WH-74f2hVbokygR5XV6rJWO-vuYuBtQZrtl1Db2BlSkR-QjselZGWdFP9BABWoUClcZ4SV7LsqhTWYbbwO-N0_gjNE0Je46NwhGqRw",
    custoken3: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM0NDg2MTgsImlhdCI6MTYwMzQ0NTAxOCwibmFtZSI6ImN1c3QzIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjZ9.fOEE5Eu5WGk1qEDQ3j770XdKF7hL0OWa5RLI1XcHAvMSWM8M-2ckm0_n61S7Zu8RftS9k6cEBNoc3QUHu6UvFwyRd-S9MCTytYnwQ3qBjI6T5iVTlPyBHjPXLkYiFj_1-iZtgHV6w2Lj6SWT10tCMYwHFWQWqZETyznCsObNJOmnf5EEh7l-68P2N8jwV6bjSNu1GBTzis2XXYYojXG6NFBHuIGYrKIAi5yKNr9EXJrFrlU_2MB8Ref8mSVsaMsGDZcBGxP6NAFimjJTdsARnvu1P08-a2qZDX0PBDaPKx85zRIT4_szFSpfOU4_5Q5FB_FEJsc96wlieYKXvw05qw",
    eventID1: "d090f582-e72c-4710-9bd9-6276637e8dd9",
    quotaLoad: function (token, eventID, amount) {
        const url = 'http://localhost:9092/api/v1/reservation/reserve'
        autocannon({
            url: url,
            connections: 1000,
            duration: 5,
            requests: [
                {
                    method: 'POST',
                    body: JSON.stringify({
                        "authToken": token,
                        "eventID": eventID,
                        "amount": amount
                    }),
                },
            ]
        }, (err, res) => {
            console.log('finished bench', err, res)
        })

    }
}