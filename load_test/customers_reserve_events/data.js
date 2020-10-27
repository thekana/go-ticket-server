const autocannon = require('autocannon')

module.exports = {
    custoken1: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM3OTYwNDYsImlhdCI6MTYwMzc3ODA0NiwibmFtZSI6ImN1c3QxIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjR9.wCgoHLOlVNXxjDG5anmTMT_XxYgn3_PR3y9JYQS-mInZHBjfTRA8LUqNBLK5bRWtvNOmh-x0M65lr458DeAoIrn94B-qQ_8QuG8ci2UYYWM3YSaHTAsuGuS1Es55u9EsWBCqUf_sJ5bY7vinWRbxFaV3nAFHcqjJa36uN8xe25r-YEBQXJCwuuNac-Hf83mEWFIJkDmkPlXHBpNFXZk64zUSoNqTWVFCwIsiN2x4LRekHJYNQZ7eDEFOOhEP4d46_kWGQF0_vB5GpbvAZxBLa3hf99trbzVX0nOzFMVru53j4MyPLEx-r8JmIn1_7CPDbQALmhig3mSPbP_ieDyXQA",
    custoken2: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM3OTYwNDksImlhdCI6MTYwMzc3ODA0OSwibmFtZSI6ImN1c3QyIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjV9.oLwrz5lG3OxSkJoYQDsflEmBRl0nMjz20ZMRDKad8CvAZ04M0bd2S3B6x-buqB6X1yZkYQTyAceDrXKEHc7oia4eQiLmWJAg8KzVBCnkCUV53ez1n7O2yUUx62YDD8kw6xnJqKELYruFbKUbr0Uk6FA-nt7PmdEKG43o3p8zTJVNCQvoGmRL0_fP3Y-7we3y3_PsLPiVTN1w9rKdZVXOAaBCdnrfVo_TRyZG2atcEFtdnP05pB2aS3v8gVhblu1fSqZ6tq7Mvtuns_lRSFvqCMgxCZcfuowykJpJxiocKVLTRLZ_X7cZqThbJ77V-hLJ9QwA34tGhI7Ku9kpF7b3GA",
    custoken3: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM3OTYwNTEsImlhdCI6MTYwMzc3ODA1MSwibmFtZSI6ImN1c3QzIiwicm9sZSI6WyJjdXN0b21lciJdLCJ1aWQiOjZ9.XXCAscdKS0HH3bomMPRrw78VsHDvqsWRfUvp4jidRu0zC20TsynAQ_XoNlWeVWgdunsqhKbvT5TP_gDLP650t0lMZG8KQFpmE-Ms0LvMwWiiMbA6H2iOhOpr-krcWm6UIEBtdGudJPxViebggqOA9nGIfI3mn2_pehkEP9Llruv-HiOZlmEN46uxkkOxwkELQtD7LdELPwCe6opgoCzIiy-ZW9jJ5k90zpEd8uYDMTlweu54iFzqT7PYdfXBl_FeGN51AUASUDnTR-Gu1bTxx1YoFH3YHY7z0etZj_RISJnp57qT3UKHshXdW2eW6QfHaSKbqTVkUjjL-pc9SNYPWg",
    eventID1: 1,
    eventID2: 2,
    eventID3: 3,
    eventID4: 4,
    quotaLoad: function (requestList) {
        const url = 'http://localhost:9092/api/v1/reservation/reserve'
        autocannon({
            url: url,
            connections: 30,
            duration: 10,
            requests: requestList
        }, (err, res) => {
            console.log('finished bench', err, res)
        })

    }
}