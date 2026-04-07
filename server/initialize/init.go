package initialize

func StartInitSystemConfig(deps Deps) {
	Migrate(deps)
	Site(deps)
	Node(deps)
	Email(deps)
	Device(deps)
	Invite(deps)
	Verify(deps)
	Subscribe(deps)
	Register(deps)
	Mobile(deps)
	Currency(deps)
	if !deps.currentConfig().Debug {
		Telegram(deps)
	}
}
