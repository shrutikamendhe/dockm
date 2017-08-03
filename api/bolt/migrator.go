package bolt

import "github.com/shrutikamendhe/dockm/api"

// Migrator defines a service to migrate data after a DockM version update.
type Migrator struct {
	UserService            *UserService
	EndpointService        *EndpointService
	ResourceControlService *ResourceControlService
	VersionService         *VersionService
	CurrentDBVersion       int
	store                  *Store
}

// NewMigrator creates a new Migrator.
func NewMigrator(store *Store, version int) *Migrator {
	return &Migrator{
		UserService:            store.UserService,
		EndpointService:        store.EndpointService,
		ResourceControlService: store.ResourceControlService,
		VersionService:         store.VersionService,
		CurrentDBVersion:       version,
		store:                  store,
	}
}

// Migrate checks the database version and migrate the existing data to the most recent data model.
func (m *Migrator) Migrate() error {

	// DockM < 1.12
	if m.CurrentDBVersion == 0 {
		err := m.updateAdminUserToDBVersion1()
		if err != nil {
			return err
		}
	}

	// DockM 1.12.x
	if m.CurrentDBVersion == 1 {
		err := m.updateResourceControlsToDBVersion2()
		if err != nil {
			return err
		}
		err = m.updateEndpointsToDBVersion2()
		if err != nil {
			return err
		}
	}

	err := m.VersionService.StoreDBVersion(dockm.DBVersion)
	if err != nil {
		return err
	}
	return nil
}
