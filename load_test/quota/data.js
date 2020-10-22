const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMzOTIyMTUsImlhdCI6MTYwMzM4ODYxNSwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.MWmndtQJOl49aGFlDf0ZnmU9FEoGCMUELnwfenYwyGC4fzOiM_apban9fBhTSHp-nGlZjvk1l8_bisktkFW__cmiAewqil4N9-BCn8Qwssq9iiQ9trozRDFdIybHb-FhBpyGAJt3q-q2Cwlzxc7TywOF9IDfNAtjQmRw73Iyxv-Cxx8zWSh1FCslfqvfcYWMM9Gm8azcPxXysoo9PfMdRx_ahRR5O8mVn12mMOkaI4cTJiVsm90W09vhm9YIS3FdxlMnlrmFMncH1Q-OJBxDY-wukkZ6IBe7vNkr5A4cgmSMuRRb415DKK1TFURWzAscVKkIAfbdDH8BFAtcP2Dh1Q",
    eventID1: "f659f0f0-188f-4e41-9de2-9e5a87594ade",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMzOTIyNzcsImlhdCI6MTYwMzM4ODY3NywibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.EL4wRtjvTnHHWRAJyktandZgMmgCFd7bk4WbXOdye1cxXwnNl5_8_u6vV-gDvyCXDTtuMNiFpE7up-C30fvksy-vT3bZhaQ4EwQH_pRX9iBrrwvo-9PMsQZcJeEtCtZS5KZQoLDDQLTeqLj9RGLm0Gwtf6-OMG9Zuw56IId7ocJI9ORcMc90VzHXFq3g_YjRxDbbkd9YW2ZKy_dn8mvu6-y_XthWaiOVeWO5l-G5l_gT7bVqtAkfn9It6ataERI54cc1kECRsq6eX76uHxSIioM3pUbrERH_-_RGFyex5eV-RQmuMk7Jnlaq0hd106Bfy0HyC3Du9MLQRXFWoEkDHQ",
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