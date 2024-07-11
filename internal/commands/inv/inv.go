package inv

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mccune1224/betrayal/internal/discord"
	"github.com/zekrotja/ken"
)

type Inv struct {
	dbPool *pgxpool.Pool
}

func (i *Inv) Initialize(pool *pgxpool.Pool) {
	i.dbPool = pool
}

// Description implements ken.SlashCommand.
func (i *Inv) Description() string {
	return "Player Inventory"
}

// Name implements ken.SlashCommand.
func (i *Inv) Name() string {
	return "inv"
}

// Options implements ken.SlashCommand.
func (i *Inv) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		i.abilityCommandArgBuilder(),
		i.coinCommandArgBuilder(),
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("role", "Role to create inventory for", true),
				discord.UserCommandArg(true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "Delete a player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "get",
			Description: "Get player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
			},
		},
	}
}

// Run implements ken.SlashCommand.
func (i *Inv) Run(ctx ken.Context) (err error) {
	return ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "create", Run: i.create},
		ken.SubCommandHandler{Name: "delete", Run: i.delete},
		ken.SubCommandHandler{Name: "get", Run: i.get},
		i.abilityCommandGroupBuilder(),
		i.coinCommandGroupBuilder(),
		// ken.SubCommandGroup{Name: "immunity", SubHandler: []ken.CommandHandler{
		// 	ken.SubCommandHandler{Name: "add", Run: i.addImmunity},
		// 	ken.SubCommandHandler{Name: "remove", Run: i.removeImmunity},
		// }},
	)
}

// Version implements ken.SlashCommand.
func (i *Inv) Version() string {
	return "1.0.0"
}

var _ ken.SlashCommand = (*Inv)(nil)
