const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM0NjU2MzYsImlhdCI6MTYwMzQ0NzYzNiwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.zlEwL99gfP6UZZDf1wZ7deu7fmKCK9CRC0grtOGFWLASojyjKr9fRrDAsWiayCaNjfYBpRY1eoGL53cttMQpB8P5c8OPaORdbwhi4TEo1vuv3jKKCutuxV5W5_enngVzYrQQKIxBBEqjEVSy_sjjJlA8-dJp_TdLBJSP9ZhdCD3FtrUvh3eK8BdeuGv0XbGCdfowaI8AgT5btQuqCDq5mOvR1lwpRhH2qRjVkHfCLDAz5_lS795NN-s643ffUOCwHyQTXlWzlMR2LEJMm4ijHcUSVGoGSrRHqc_Gu1CkrawqR9R5TF-EXK3op7lW9PR4qZacry2TVH6kuh4_Hi2yMg",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM0NjU2NDYsImlhdCI6MTYwMzQ0NzY0NiwibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.OLjbsu83oRJYIyXFB7KuYKDwRxDRP5xAbC-qBsU6X9Ft1tZD5e7blmiVwSSfZCol0pGk-fkilZx-XGyUM2cdLecuyP_bdycTW3OavG-ZfoNfGLN8vxPaPT_aE8gBd3afx_-EypHwvaffQG2S8N81BF3opxkHZVtX3yW-dseRgStM_ysFc1JviIgcI39mtqZqweoR_ft3K-HzmqeRiRJvHe_tiFkANAY6Tsl020j583yFC1WCfSa1vIr9ykXTPWEoFYlzItsXOnCz_f-N9gQMY8kblFOhDAUD294WeLH7_VNBOueR22pDLgnwdJWn2iZ5vvCJFC4IeRhsE_rQN6Y0vw",
    custoken3: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM0NjU2NjQsImlhdCI6MTYwMzQ0NzY2NCwibmFtZSI6ImN1c3QzIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjZ9.y9KvK11qxZXm34XBqn0fiuMcqgTkdX3tyEr-kKfTAoY-CrfKmVfHznhuGPzeFJKSXmX2TM98GxWRh7U1UgHF2YrfkhntlIoalrc-oJ9JvLmOFhLsCmU62sEyh4wkYcG8zQ8vR7_mpoVEZeBsJgYUNj_TsHLBqmL8phNCovZ4UUNwknwTtyzk8BP6q6JfNmlGYEV0ZKUy4hxcIOXj1g2DFPXtltdL-YidmJZ03QcxeKgTYHtDjfOv44StlhJDrh37XFXXb-A09eHm-nURqvkZS1ZO3UAtYW7HVhR9_OtqquTzNJi8CNVy0KfrqvBq1SNyN3zv7WVp3zzsykua0QaJAQ",
    eventID1: "8370dc97-e6bc-4089-802d-861a02b7effc",
    eventID2: "c99f9be0-f217-4a44-b548-44b27f67c08c",
    eventID3: "9acec916-cfd2-4d02-8867-4aa38b402ce8",
    eventID4: "ca3abd76-ebd5-4cf1-8e89-c771fc357470",
    quotaLoad: function (requestList) {
        const url = 'http://localhost:9092/api/v1/reservation/reserve'
        autocannon({
            url: url,
            connections: 1000,
            duration: 5,
            requests: requestList
        }, (err, res) => {
            console.log('finished bench', err, res)
        })

    }
}