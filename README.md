# github.com/KronusRodion/protocol
This is repository for application Layer Protocol in the OSI Model. Protocol presents 2 methods over TCP connection - adding key/value and getting value by key. An example of a example of server with simple in-memory storage is also presented here.

# Start
To start a server

```bash
go run cmd/server.main.go
```

For using Get and Send methods import pkg part of project and use Client:

```bash
client, err := NewClient(port)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	key := []byte("test_key")
	value := []byte("test_value")

	err = client.Send(key, value)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	// Проверяем метод Get
	value, err = client.Get([]byte("test_key"))
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
```

# Tests

To run the test:
```bash
go test .\pkg\client\ -v -count=1
```