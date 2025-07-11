# Cara Kerjanya

suite.Run() akan mencari method dengan nama spesifik:

- SetupSuite()
- TearDownSuite()
- SetupTest()
- TearDownTest()

Urutan Eksekusi:
suite.Run() →
SetupSuite() →
(SetupTest() → TestA() → TearDownTest()) →
(SetupTest() → TestB() → TearDownTest()) →
TearDownSuite()
