const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM4MTc2MDYsImlhdCI6MTYwMzc5OTYwNiwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.UArnSbCvUfsDbSO9dxUqtimLi-rHJFG9aI9veK3QHqyXfmmIUaoPnlOoXqCYtE4i0OIN9ttIK8hj8lXY5RUGKtZ2rBP-0v8SNV143f_Y2UZD8pFbA1uZPpGX5FuoaDauBJKWRB-LTL1_ASGaR7rWtSDv7MaFZSIxCiUH6C7cdqJ_TRdVcR-UdgvpoPF9Hio8-ZGQYbGtVjLY0QHpNxQGhB4TPnCONIDZREpahj2tDdkXxOPZtVhrQ8io_EMX5SsTx2TJ34nbNTnawb41qcj3x9J0kgTdF__SAiko0-056cplKnUHvsf_RokE1IvjbHKwo-IipWYE6abmThdQ6u9U9Q",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM4MTc2MTksImlhdCI6MTYwMzc5OTYxOSwibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.LAwRPk1co3c0JY0yrxD5qVabXs4bzMpUYq3V0eXfrJ2BSJGs00cKDmsIX8I67_Z6zw2FcGnkju7nZqJ9BHo_rhVW1AoZlPITSAmQYlCbj-UWOGfCTjIcB0BCRrGOH650dwgEBCBAj4BB7MVmPEhWOUBc4cJebDYI7Gn2TwOQv2PRbz4guvrhtFMcSSup1Cgj-Gj08C0rSzz_V6eo0lNHRRaPZiz-yTZ2uIN5_I21X9dk3ZWuSoL2uAxR9v4OsuP2DhUbwaFC3tTNtSoaO_pZob3tFEzRLU2zjvVSd2T3ndiYT0hWCZ_YcRjTGiOBMs_ewd3NdPfJ94eJHXx-bgnmcg",
    custoken3: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM4MTc2MzgsImlhdCI6MTYwMzc5OTYzOCwibmFtZSI6ImN1c3QzIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjZ9.r0ibDfco0vphM0WVv2X5yaJKAN3xP_e8jjULGsXmlrfi45KucNSf97NVvw0mSaqNtMpumnpc73pwFw1n78mXi4EnUnNRyBuH_uBxmja2zNP6b2WV60Mj7D6BHZQbWEsVaHIK1fqg7YkLXUvB_OWfjleSAUqzq8f-nIAryYuBgJiGWMvkEiiH5x8iqtD9jj3qlo9Qk3fpCKKHf80p3P2CMuxotgU59AHzTsXHsOI0TRAqEDy6t2W_5TY3A3WibKWC3F_29CLTQFL9LS6JfzoK3q8p5WUOB7kGY4sd97SJciiEY0gTETYu71a_ZKeyzCe18MEIpmmoA7dOJizzesWpVQ",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
    quotaLoad: function (requestList) {
        const url = 'http://localhost:9092/api/v1/reservation/reserve'
        autocannon({
            url: url,
            connections: 80,
            duration: 20,
            requests: requestList,
            excludeErrorStats: true
        }, (err, res) => {
            console.log('finished bench', err, res)
        })

    }
}