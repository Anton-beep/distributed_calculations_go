package availableservers

import (
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"sync"
)

type AvailableServers struct {
	servers     sync.Map
	expressions *expressionstorage.ExpressionStorage
}

func New(expressions *expressionstorage.ExpressionStorage) *AvailableServers {
	return &AvailableServers{
		expressions: expressions,
	}
}

func (a *AvailableServers) Add(server string) {
	a.servers.Store(server, true)
}

func (a *AvailableServers) Remove(server string) {
	a.servers.Delete(server)
}

func (a *AvailableServers) GetAll() []string {
	servers := make([]string, 0)
	a.servers.Range(func(key, _ interface{}) bool {
		servers = append(servers, key.(string))
		return true
	})
	return servers
}

func (a *AvailableServers) GetOperations(server string) []db.Expression {
	expressions := a.expressions.GetAll()
	operations := make([]db.Expression, 0)
	for _, expression := range expressions {
		if expression.Servername == server {
			operations = append(operations, expression)
		}
	}
	return operations
}
