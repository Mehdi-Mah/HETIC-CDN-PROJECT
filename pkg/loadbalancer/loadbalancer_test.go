package loadbalancer

import "testing"

func TestNextServer(t *testing.T) {
	lb := Instance()
	server1 := lb.NextServer()
	server2 := lb.NextServer()

	if server1 == server2 {
		t.Errorf("expected different servers, got %v and %v", server1, server2)
	}
}
