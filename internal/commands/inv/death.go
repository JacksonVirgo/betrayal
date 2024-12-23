package inv

import (
	"context"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal/internal/commands/channels"
	"github.com/mccune1224/betrayal/internal/discord"
	"github.com/mccune1224/betrayal/internal/models"
	"github.com/mccune1224/betrayal/internal/services/inventory"
	"github.com/zekrotja/ken"
)

// The bot should be upset someone is revived or set alive
var playerSetAliveMessages = []string{
	"Bummer...",
	"Congrats I guess :/",
	"Sorry hosts...",
	"I'd be lying if I was happy...",
}

// All the messages here should celebrate the player's death
var playerSetDeadMessages = []string{
	"Yipee!!",
	"Ya love to see it",
	"I'm not really sorry for your loss...",
}

func getRandomItem[T any](slice []T) T {
	randomIndex := rand.Intn(len(slice))
	return slice[randomIndex]
}

func (i *Inv) deathCommandGroupBuilder() ken.SubCommandGroup {
	return ken.SubCommandGroup{Name: "death_status", SubHandler: []ken.CommandHandler{
		ken.SubCommandHandler{Name: "alive", Run: i.setAlive},
		ken.SubCommandHandler{Name: "dead", Run: i.setDead},
	}}
}
func (i *Inv) deathCommandArgBuilder() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        "death_status",
		Description: "set the player's death status",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "alive",
				Description: "set the player to alive",
				Options: []*discordgo.ApplicationCommandOption{
					discord.UserCommandArg(false),
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "dead",
				Description: "set the player to dead",
				Options: []*discordgo.ApplicationCommandOption{
					discord.UserCommandArg(false),
				},
			},
		},
	}
}

func (i *Inv) setAlive(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		log.Println(err)
		return err
	}
	if !discord.IsAdminRole(ctx, discord.AdminRoles...) {
		return discord.NotAdminError(ctx)
	}
	h, err := inventory.NewInventoryHandler(ctx, i.dbPool)
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "failed to init inv handler")
	}
	defer h.UpdateInventoryMessage(ctx.GetSession())
	if h.GetPlayer().Alive {
		return discord.ErrorMessage(ctx, "Already Alive", "Player is already alive, bummer...")
	}
	q := models.New(i.dbPool)
	_, err = q.UpdatePlayerAlive(context.Background(), models.UpdatePlayerAliveParams{
		ID:    h.GetPlayer().ID,
		Alive: true,
	})
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Failed to set player alive")
	}

	lifeboard, err := q.GetPlayerLifeboard(context.Background())
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Failed to get player lifeboard")
	}

	playerLifeStatuses, _ := q.ListPlayerLifeboard(context.Background())
	msg, err := channels.UserLifeboardMessageBuilder(ctx.GetSession(), playerLifeStatuses)
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Failed to build user lifeboard message")
	}

	ctx.GetSession().ChannelMessageEditEmbed(lifeboard.ChannelID, lifeboard.MessageID, msg)

	return discord.SuccessfulMessage(ctx, "Player Alive", "Player is now alive\n"+getRandomItem(playerSetAliveMessages))
}

func (i *Inv) setDead(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		log.Println(err)
		return err
	}
	if !discord.IsAdminRole(ctx, discord.AdminRoles...) {
		return discord.NotAdminError(ctx)
	}
	h, err := inventory.NewInventoryHandler(ctx, i.dbPool)
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "failed to init inv handler")
	}
	defer h.UpdateInventoryMessage(ctx.GetSession())
	if !h.GetPlayer().Alive {
		return discord.ErrorMessage(ctx, "Already Dead", "Player is already alive, Great!")
	}
	q := models.New(i.dbPool)
	_, err = q.UpdatePlayerAlive(context.Background(), models.UpdatePlayerAliveParams{
		ID:    h.GetPlayer().ID,
		Alive: false,
	})
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Failed to set player dead")
	}

	lifeboard, err := q.GetPlayerLifeboard(context.Background())
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Failed to get player lifeboard")
	}

	playerLifeStatuses, _ := q.ListPlayerLifeboard(context.Background())
	msg, err := channels.UserLifeboardMessageBuilder(ctx.GetSession(), playerLifeStatuses)
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Failed to build user lifeboard message")
	}

	inv, err := h.FetchInventory()
	if err != nil {
		ctx.GetSession().ChannelMessageEditEmbed(lifeboard.ChannelID, lifeboard.MessageID, msg)
		return discord.SuccessfulMessage(ctx, "Player Dead", "Player is now dead\n" + getRandomItem(playerSetDeadMessages) + "\n\n**Make sure to check for Lucky Coin**")
	}

	hasLuckyCoin := false
	for _, item := range inv.Items {
		if item.Name == "Lucky Coin" {
			hasLuckyCoin = true
			break
		}
	}

	messageContent := "Player is now dead\n" + getRandomItem(playerSetDeadMessages)
	if hasLuckyCoin {
		messageContent += "\n\n**Host/s, this inventory holds the Lucky Coin!**"
	}

	ctx.GetSession().ChannelMessageEditEmbed(lifeboard.ChannelID, lifeboard.MessageID, msg)
	return discord.SuccessfulMessage(ctx, "Player Dead", messageContent)

}
