package middlewares

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

// Now this is peak cybersec.
type RequiresRoleCommand interface {
	RequiresRole() string
}

type PermissionsMiddleware struct {
	count int
}

var (
	_ ken.MiddlewareBefore = (*PermissionsMiddleware)(nil)
)

func (c *PermissionsMiddleware) Before(ctx *ken.Ctx) (next bool, err error) {
	cmd, ok := ctx.Command.(RequiresRoleCommand)
	if !ok {
		next = true
		return
	}

	guildRoles, err := ctx.GetSession().GuildRoles(ctx.GetEvent().GuildID)
	if err != nil {
		return
	}

	// will I ever use labels again? probably not. but they're cool.
roleLoop:
	for _, rid := range ctx.GetEvent().Member.Roles {
		for _, r := range guildRoles {
			if rid == r.ID && r.Name == cmd.RequiresRole() {
				next = true
				break roleLoop
			}
		}
	}

	if !next {
		ctx.SetEphemeral(true)
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color: 0xF44336,
						Description: fmt.Sprintf(
							"You must have the role \"%s\" to perform this command!",
							cmd.RequiresRole()),
					},
				},
			},
		})
		ctx.SetEphemeral(false)
	}
	return
}
