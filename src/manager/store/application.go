package store

import ()

func (zk *ZKStore) CreateApp(app *Application) error {
	if zk.GetApp(app.ID) != nil {
		return ErrAppAlreadyExists
	}

	op := &AtomicOp{
		Op:      OP_ADD,
		Entity:  ENTITY_APP,
		Param1:  app.ID,
		Payload: app,
	}

	return zk.Apply(op, true)
}

func (zk *ZKStore) UpdateApp(app *Application) error {
	if zk.GetApp(app.ID) == nil {
		return ErrAppNotFound
	}

	op := &AtomicOp{
		Op:      OP_UPDATE,
		Entity:  ENTITY_APP,
		Param1:  app.ID,
		Payload: app,
	}

	return zk.Apply(op, true)
}

func (zk *ZKStore) GetApp(appId string) *Application {
	zk.mu.RLock()
	defer zk.mu.RUnlock()

	appStore, found := zk.Storage.Apps[appId]
	if !found {
		return nil
	}

	return appStore.App
}

func (zk *ZKStore) ListApps() []*Application {
	zk.mu.RLock()
	defer zk.mu.RUnlock()

	apps := make([]*Application, 0)
	for _, app := range zk.Storage.Apps {
		apps = append(apps, app.App)
	}

	return apps
}

func (zk *ZKStore) DeleteApp(appId string) error {
	if zk.GetApp(appId) == nil {
		return ErrAppNotFound
	}

	op := &AtomicOp{
		Op:     OP_REMOVE,
		Entity: ENTITY_APP,
		Param1: appId,
	}

	return zk.Apply(op, true)
}
