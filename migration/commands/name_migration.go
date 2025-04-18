package commands

type NameMigration struct{}

func NewNameMigration() *NameMigration {
	return &NameMigration{}
}

func (n *NameMigration) Name() string {
	return "Name Migration"
}

func (n *NameMigration) Execute() error {
	return MigrateJSONToMongo("name", "migration/name.json")
}
