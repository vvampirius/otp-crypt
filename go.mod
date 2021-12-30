module github.com/vvampirius/otp-crypt

go 1.17

replace github.com/vvampirius/otp-crypt/coder => ./coder

require (
	github.com/pquerna/otp v1.3.0
	github.com/vvampirius/otp-crypt/coder v0.0.0-00010101000000-000000000000
)

require (
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/vvampirius/mygolibs/encryption v0.0.0-20211230132418-efcfd50bfb54 // indirect
	github.com/vvampirius/mygolibs/rwblocks v0.0.0-20211230132418-efcfd50bfb54 // indirect
)
