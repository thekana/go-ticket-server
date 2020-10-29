const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM5NjU2NDIsImlhdCI6MTYwMzk0NzY0MiwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.LhP7iJbtqUGh9CL6vvJkX49hdNUvkELTCfUfvRlPJgN9NhXmx4Q262XZ_WQeZBwxdpLBFnGvWYQi3wZ_cK5FBFMN6s2dgJq8dbZtEOH7ATS62gUp-65c75onzP5PZas1FUz6XDwqHwHgcPyfnsDVrFgnPXdOyv-it-6LfsPbM-KykiUqdw7GlLzXTi51sjTZtslgGsfNBCMhdPCaNe5EfZ2t6H_Mn5uh98BS2M2Nap_u6rePli3EOodJr4E_shsIQ9hXgmdlDHsUQvrlkdZ54Nxy3TOKvx3CH0If88O5rV0Gaps3bSqkxtQFQVUO0W_HTvNamwEWL72l7QMZKrVJ9g",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM5NjU2NDQsImlhdCI6MTYwMzk0NzY0NCwibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.Gn6tEbS_x5WOOIg2uhGf2BHntWIScmO3J15kvDHbpuqgcdFWoP6Ft2xpY4IaMfAivDUFc6mMT_KVcbbBraT-yyEP8tiIjifczU2ovkzhsXqZbgZwtWFFD0kIXTZdZjlbri1aJodkoY_N1RxQ7MENqsdFILsBujnrKWURLoOw54JZSMW7JS4-N3rXW8Qji-G399kbiD8s0QlAIxuHPjVYGI-OiHo4vCz-TZOQ5cidPKysWMmcGjs_ZFMPpG_LF5vfqc0Tz02sap4CE72jKcSGl9PyOBn03-_xbSbyWJUscqNuU7oatVW8fgWFT0GllbTSpxvYwwecr02OVNfbYx5MjQ",
    custoken3: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM5NjU2NDUsImlhdCI6MTYwMzk0NzY0NSwibmFtZSI6ImN1c3QzIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjZ9.NpXyodF3xqY80cA8hXKZeGDQ9EDkeFXkgJY3EYGgd_HnFUFv9oy0xQgGL6k9Ac5HatzwzUYQ-ktg4SftojUyGZjGNjL2UIY1eo5-_wVSgSUxEcfozulib4jrWXZdW-JlzaGHXEwYmGmdqb7otSRmJVlIFZtTgcZ60d90KTSnI8OQSkqfR29me-CyqKBSe2kINOtyRBTFtQbx-XQzgAfJsdC6dj83d62s9gBdfah_di3UPkWH5dY3Ys1y3-G07vIFvpcGUw_Egfejx4Mm6dxytt_WLjyFq5a0SYeqq5g2HXr_Xu-pEEf4Kb6hqkiDoSA48CN71Yv2JT-lVjs3TFwf5A",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
    quotaLoad: function (requestList) {
        const url = 'http://localhost:9092/api/v1/reservation/reserve'
        autocannon({
            url: url,
            connections: 30,
            duration: 20,
            requests: requestList,
            excludeErrorStats: true
        }, (err, res) => {
            console.log('finished bench', err, res)
        })

    }
}