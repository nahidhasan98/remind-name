package commands

type PlatformMigration struct{}

func NewPlatformMigration() *PlatformMigration {
	return &PlatformMigration{}
}

func (p *PlatformMigration) Name() string {
	return "Platform Migration"
}

func (p *PlatformMigration) Execute() error {
	return MigrateJSONToMongo("platform", "migration/platform.json")
}
