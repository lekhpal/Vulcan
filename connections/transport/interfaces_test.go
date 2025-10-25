package transport

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockConnection is a mock implementation of the Connection interface for testing.
type MockConnection struct {
	Connected    bool
	SendError    error
	ReceiveData  []byte
	ReceiveError error
	ReceivedData [][]byte
}

func (m *MockConnection) Connect(ctx context.Context) error {
	m.Connected = true
	return nil
}

func (m *MockConnection) Disconnect(ctx context.Context) error {
	m.Connected = false
	return nil
}

func (m *MockConnection) Send(ctx context.Context, data []byte) error {
	if m.SendError != nil {
		return m.SendError
	}
	m.ReceivedData = append(m.ReceivedData, data)
	return nil
}

func (m *MockConnection) Receive(ctx context.Context) ([]byte, error) {
	if m.ReceiveError != nil {
		return nil, m.ReceiveError
	}
	return m.ReceiveData, nil
}

func (m *MockConnection) IsConnected() bool {
	return m.Connected
}

// MockMiddleware is a mock implementation of the Middleware interface.
type MockMiddleware struct {
	HandleFunc func(ctx context.Context, conn Connection, next Connection) (Connection, error)
}

func (m *MockMiddleware) Handle(ctx context.Context, conn Connection, next Connection) (Connection, error) {
	if m.HandleFunc != nil {
		return m.HandleFunc(ctx, conn, next)
	}
	return conn, next.Connect(ctx) // Default behavior: just connect
}

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	CallFunc func(ctx context.Context, method string, request []byte) ([]byte, error)
}

func (m *MockClient) Call(ctx context.Context, method string, request []byte) ([]byte, error) {
	if m.CallFunc != nil {
		return m.CallFunc(ctx, method, request)
	}
	return nil, nil
}

func TestConfig(t *testing.T) {
	cfg := Config{
		Timeout:    5 * time.Second,
		BufferSize: 1024,
	}

	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.Equal(t, 1024, cfg.BufferSize)
}

func TestMockConnection_ConnectDisconnect(t *testing.T) {
	conn := &MockConnection{}
	ctx := context.Background()

	assert.False(t, conn.IsConnected())

	err := conn.Connect(ctx)
	assert.NoError(t, err)
	assert.True(t, conn.IsConnected())

	err = conn.Disconnect(ctx)
	assert.NoError(t, err)
	assert.False(t, conn.IsConnected())
}

func TestMockConnection_SendReceive(t *testing.T) {
	conn := &MockConnection{ReceiveData: []byte("hello")}
	ctx := context.Background()

	err := conn.Connect(ctx)
	assert.NoError(t, err)

	sendData := []byte("world")
	err = conn.Send(ctx, sendData)
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{sendData}, conn.ReceivedData)

	receivedData, err := conn.Receive(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), receivedData)
}

func TestMockConnection_SendError(t *testing.T) {
	expectedErr := errors.New("send error")
	conn := &MockConnection{SendError: expectedErr}
	ctx := context.Background()

	err := conn.Send(ctx, []byte("data"))
	assert.EqualError(t, err, expectedErr.Error())
}

func TestMockConnection_ReceiveError(t *testing.T) {
	expectedErr := errors.New("receive error")
	conn := &MockConnection{ReceiveError: expectedErr}
	ctx := context.Background()

	_, err := conn.Receive(ctx)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestMockMiddleware_Handle(t *testing.T) {
	mockConn := &MockConnection{}
	ctx := context.Background()

	middleware := &MockMiddleware{
		HandleFunc: func(ctx context.Context, conn Connection, next Connection) (Connection, error) {
			err := conn.Connect(ctx)
			assert.NoError(t, err)
			return conn, nil
		},
	}

	_, err := middleware.Handle(ctx, mockConn, mockConn)
	assert.NoError(t, err)
	assert.True(t, mockConn.IsConnected())
}

func TestMockClient_Call(t *testing.T) {
	expectedResponse := []byte("client response")
	client := &MockClient{
		CallFunc: func(ctx context.Context, method string, request []byte) ([]byte, error) {
			assert.Equal(t, "testMethod", method)
			assert.Equal(t, []byte("testRequest"), request)
			return expectedResponse, nil
		},
	}
	ctx := context.Background()

	response, err := client.Call(ctx, "testMethod", []byte("testRequest"))
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}
