package orchiestrator

import "github.com/kuzxnia/loadcli/orchiestrator/resourcemanager"

type Orchiestrator struct {
	resourceManager resourcemanager.ResourceManager
}

func NewOrchiestrator() (*Orchiestrator, error) {
	// wybranie odpowiedniej strategi ? -
	resourceManager, err := resourcemanager.GetResourceManager()
	if err != nil {
		return nil, err
	}

	return &Orchiestrator{
		resourceManager: resourceManager,
	}, nil
}

func (o *Orchiestrator) createResources() {
	o.resourceManager.Install()
}

func (o *Orchiestrator) deleteResources() {
	o.resourceManager.UnInstall()
}

// nessesary? yes not allways we will have helms
