package availableservers

import (
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"sync"
)

type AvailableServers struct {
	servers     []string
	expressions *expressionstorage.ExpressionStorage
	mu          sync.Mutex
}

func New(expressions *expressionstorage.ExpressionStorage) *AvailableServers {
	return &AvailableServers{
		expressions: expressions,
	}
}

func (a *AvailableServers) Add(server string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, s := range a.servers {
		if s == server {
			return
		}
	}
	a.servers = append(a.servers, server)
}

func (a *AvailableServers) Remove(server string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, s := range a.servers {
		if s == server {
			a.servers = append(a.servers[:i], a.servers[i+1:]...)
			break
		}
	}
}

func (a *AvailableServers) GetAll() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.servers
}

func (a *AvailableServers) GetExpressions(server string) []db.Expression {
	res := make([]db.Expression, 0)
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, s := range a.servers {
		if s == server {
			res = append(res, a.expressions.GetByServer(s)...)
		}
	}
	return res
}
