package service

import "github.com/google/wire"

// ProviderSet is service providers.

var ProviderSet = wire.NewSet(
    func() []string {
        return []string{"http://localhost:2379"} 
    },
    NewCardService,
)