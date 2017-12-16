# Codeship API Integration Tests

This directory contains additional test suites beyond the unit tests. Whereas the unit tests run very quickly (since they don't make any network calls), the tests in this directory are only run manually or in CI.

In short, these tests aim to be:

* **Non-Destructive:** Since these tests are setup to run against a live Codeship account, we have omitted any 'destructive' type tests such as those that create, update or delete data until we find a better way of testing these API calls.
* **Relatively Quick:** Even though these tests make real network calls, we still aim for them to be relatively quick (finish under a minute or so). Therefore these tests do not over exercise pagination or make multiple calls for the same data.
* **Happy Path Only:** Currently these tests only test the 'happy paths' meaning that we are testing for data that we know to exist and are not expecting any error conditions. This may change in future iterations of this test suite.

## Rate Limiting

Because these tests are making live network calls to a real Codeship account, they must follow the same rules regarding rate limiting as defined at: [https://apidocs.codeship.com/v2/introduction/rate-limiting](https://apidocs.codeship.com/v2/introduction/rate-limiting).

Be aware of the rate limit when adding new tests or when requesting data in loops. It may be necessary to use `time.Sleep` in some cases to avoid hitting the rate limit.

## Environment Variables

The following environment variables are **required** to be set before running this suite:

* `CODESHIP_USER` - the user the tests will use to authenticate
* `CODESHIP_PASSWORD` - the password of the user used to authenticate

## Running

Run these tests with exported environment variables:

`go test -v -tags=integration ./integration/...`

or without exporting environment variables:

`CODESHIP_USERNAME=XXX CODESHIP_PASSWORD=XXX go test -v -tags=integration ./integration/...`
