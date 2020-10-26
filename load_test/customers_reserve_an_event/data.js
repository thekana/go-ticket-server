const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM3NDM0NjMsImlhdCI6MTYwMzcyNTQ2MywibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.tGNVwvDzrIWvZG_KuQkDon4AUfJH5VFTOiUV8RDY3csjhEO1KpbYpE2kboGz7X_OhWesYxgih5TBCiMY7LToIusqCutwy_LfFRcRAF_MpVUnyVAr2ciZv6a8E9OCjnZERWZ5nbeJ1atRacosTyMuXNcgBXt8dTaXmOay9PG_PVMJ0GPUzwXGWsnp01Y5QBpQqM5pOMTYqZ8HJL91v2LmpSkvV5FgKhwTiZB2G1_aUXnPYORrRba_XZEoZSNiHWlu0a0xSy-HOMxrqY2mntKjHdH9NcsGOZK5Lm51kmhmzvxYXLyQR67-8CyFyXp9bXrefJ9bGrw49f3mcSvi18vlxQ",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM3NDM0NjUsImlhdCI6MTYwMzcyNTQ2NSwibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.KhCV-beU2mJM-zi36Qx06dKzwHK-nM-k7nle0ccT9NoouZkkQZckkxxy9_FlWBQeuc4pBUD8UIal36xdBMR8MW1m5zXZ6WD87q3k7qsriH_z2KW7SCaDTuh_ZMG6iGG01RMjgE5Rba56oCJnDGFuvx362A8oGPkRDOaw-fXd1VklE98i3ZTnLI_TL3OmZb6J7dY6dvOoDY-tsfX3DxZs5v30VEWqBmHGBNcKp7HSOrJ4HXiN_SsCKX2Ehmj4643jCNemb6TQMxr6A_DP--AgID1mBFUMI1WHzuxQeqkqYb6do3N4gyhNKLcGc7YKZ0ULgjiFDW3Skgt-_Xt8LhljHw",
    custoken3: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM3NDM0NjgsImlhdCI6MTYwMzcyNTQ2OCwibmFtZSI6ImN1c3QzIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjZ9.PoNeHQNVcG9xiO7ZmLceQUCOq8CiTyB769eP4bXX7YJaLqkdiFpKpQNEt39vNRSW54A99ZPccp9s-va_9K5WQWIC0zdALygYklh1VnJweuDcW_Wx9V5pexCgNRWtDuvaBAGDie5-EOHIiIDxNLKl500l3gNoEVdhr6Yjr8n0iUy9QLli7wb6asZtnYFxbGXBFISgVFC7_CzP1_hPj2d2J8SnStY5e-_Qc1wfUxwZrmt18SwLPWslvT4HdgE4o_8Ymry6-lgs5qK1FjLRU9az25ZiNzMwUJzlnK_uKLkE8nEA2rGaZMTgz-CAEQRL4N99M4pczRgGWZkecqpSJwCy0A",
    eventID1: 1,
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