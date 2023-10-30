package inventory

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal/internal/data"
	"github.com/mccune1224/betrayal/internal/discord"
	"github.com/mccune1224/betrayal/internal/util"
	"github.com/zekrotja/ken"
)

// errors that can occur
var (
	ErrNotAuthorized = errors.New("you are not an admin role")
)

// TODO: Maybe make these configurable?
const (
	defaultCoins      = 200
	defaultItemsLimit = 4
	defaultLuck       = 0
)

var optional = discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionBoolean,
	Name:        "hidden",
	Description: "Make view hidden or public (default hidden)",
	Required:    false,
}

type Inventory struct {
	models data.Models
}

// Components implements main.BetrayalCommand.
func (*Inventory) Components() []*discordgo.MessageComponent {
	return nil
}

var _ ken.SlashCommand = (*Inventory)(nil)

func (i *Inventory) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (i *Inventory) SetModels(models data.Models) {
	i.models = models
}

// Description implements ken.SlashCommand.
func (*Inventory) Description() string {
	return "Command for managing inventory"
}

// Name implements ken.SlashCommand.
func (*Inventory) Name() string {
	return discord.DebugCmd + "inv"
}

// Options implements ken.SlashCommand.
func (*Inventory) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "get",
			Description: "get player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "create a new player",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
				discord.StringCommandArg("role", "Role to assign to player", true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "delete inventory",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "whitelist",
			Description: "whitelist channel for inventory",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "add whitelist channel",
					Options: []*discordgo.ApplicationCommandOption{
						discord.ChannelCommandArg(true),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "remove whitelist channel",
					Options: []*discordgo.ApplicationCommandOption{
						discord.ChannelCommandArg(true),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "list whitelist channels",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "add",
			Description: "add to player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "ability",
					Description: "add a base ability",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the ability", true),
						discord.IntCommandArg("charges", "Number of charges", false),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "aa",
					Description: "add an any ability",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the ability", true),
						discord.IntCommandArg("charges", "Number of charges", false),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "perk",
					Description: "add a perk",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the perk", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "item",
					Description: "add an item",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the item", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "status",
					Description: "add a status",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the status", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "immunity",
					Description: "add an immunity",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the immunity", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "effect",
					Description: "add an effect",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the effect", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "coins",
					Description: "add coins",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("amount", "Amount of coins to add", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "bonus",
					Description: "add coin bonus",
					Options: []*discordgo.ApplicationCommandOption{
						// Discord is fucking stupid and doesn't allow decimals
						discord.StringCommandArg("amount", "Amount of coin bonus to add", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "luck",
					Description: "Add onto of the current luck level",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("amount", "amount of luck to add", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "note",
					Description: "add a note",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("message", "Note to add", true),
						discord.UserCommandArg(false),
					},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "remove",
			Description: "remove to player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "ability",
					Description: "remove a base ability",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the ability", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "aa",
					Description: "remove an any ability",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the ability", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "perk",
					Description: "remove a perk",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the perk", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "item",
					Description: "remove an item",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the item", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "status",
					Description: "remove a status",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the status", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "immunity",
					Description: "remove an immunity",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the immunity", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "effect",
					Description: "remove an effect",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the effect", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "coins",
					Description: "remove coins",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("amount", "Amount of coins to remove", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "luck",
					Description: "Amount of luck to remove",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("amount", "Amount of luck to remove", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "bonus",
					Description: "remove coin bonus",
					Options: []*discordgo.ApplicationCommandOption{
						// Discord is fucking stupid and doesn't take decimals...need to use string arg
						discord.StringCommandArg("amount", "Amount of coin bonus to remove", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "note",
					Description: "remove a note by index number",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("index", "Index # to remove", true),
						discord.UserCommandArg(false),
					},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "set",
			Description: "set to player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "ability",
					Description: "set a base ability",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the ability", true),
						discord.IntCommandArg("charges", "Number of charges", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "aa",
					Description: "set an any ability",
					Options: []*discordgo.ApplicationCommandOption{
						discord.StringCommandArg("name", "Name of the ability", true),
						discord.IntCommandArg("charges", "Number of charges", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "coins",
					Description: "set coins",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("amount", "Amount of coins to set", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "bonus",
					Description: "set coin bonus",
					Options: []*discordgo.ApplicationCommandOption{
						// Discord is fucking stupid and doesn't take decimals...need to use string arg
						discord.StringCommandArg("amount", "Amount of coin bonus to set", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "limit",
					Description: "Set item limit",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("size", "New size to set", true),
						discord.UserCommandArg(false),
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "luck",
					Description: "set luck level for player",
					Options: []*discordgo.ApplicationCommandOption{
						discord.IntCommandArg("level", "Set Luck level", true),
						discord.UserCommandArg(false),
					},
				},
			},
		},
	}
}

// Run implements ken.SlashCommand.
func (i *Inventory) Run(ctx ken.Context) (err error) {
	return ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "get", Run: i.get},
		ken.SubCommandHandler{Name: "create", Run: i.create},
		ken.SubCommandHandler{Name: "delete", Run: i.delete},
		ken.SubCommandGroup{Name: "whitelist", SubHandler: []ken.CommandHandler{
			ken.SubCommandHandler{Name: "add", Run: i.addWhitelist},
			ken.SubCommandHandler{Name: "remove", Run: i.removeWhitelist},
			ken.SubCommandHandler{Name: "list", Run: i.listWhitelist},
		}},
		ken.SubCommandGroup{Name: "add", SubHandler: []ken.CommandHandler{
			ken.SubCommandHandler{Name: "ability", Run: i.addAbility},
			ken.SubCommandHandler{Name: "aa", Run: i.addAnyAbility},
			ken.SubCommandHandler{Name: "perk", Run: i.addPerk},
			ken.SubCommandHandler{Name: "item", Run: i.addItem},
			ken.SubCommandHandler{Name: "status", Run: i.addStatus},
			ken.SubCommandHandler{Name: "immunity", Run: i.addImmunity},
			ken.SubCommandHandler{Name: "effect", Run: i.addEffect},
			ken.SubCommandHandler{Name: "coins", Run: i.addCoins},
			ken.SubCommandHandler{Name: "bonus", Run: i.addCoinBonus},
			ken.SubCommandHandler{Name: "luck", Run: i.addLuck},
			ken.SubCommandHandler{Name: "note", Run: i.addNote},
		}},
		ken.SubCommandGroup{Name: "remove", SubHandler: []ken.CommandHandler{
			ken.SubCommandHandler{Name: "ability", Run: i.removeAbility},
			ken.SubCommandHandler{Name: "aa", Run: i.removeAnyAbility},
			ken.SubCommandHandler{Name: "perk", Run: i.removePerk},
			ken.SubCommandHandler{Name: "item", Run: i.removeItem},
			ken.SubCommandHandler{Name: "status", Run: i.removeStatus},
			ken.SubCommandHandler{Name: "immunity", Run: i.removeImmunity},
			ken.SubCommandHandler{Name: "effect", Run: i.removeEffect},
			ken.SubCommandHandler{Name: "coins", Run: i.removeCoins},
			ken.SubCommandHandler{Name: "bonus", Run: i.removeCoinBonus},
			ken.SubCommandHandler{Name: "luck", Run: i.removeLuck},
			ken.SubCommandHandler{Name: "note", Run: i.removeNote},
		}},
		ken.SubCommandGroup{Name: "set", SubHandler: []ken.CommandHandler{
			ken.SubCommandHandler{Name: "ability", Run: i.setAbility},
			ken.SubCommandHandler{Name: "aa", Run: i.setAnyAbility},
			ken.SubCommandHandler{Name: "coins", Run: i.setCoins},
			ken.SubCommandHandler{Name: "bonus", Run: i.setCoinBonus},
			ken.SubCommandHandler{Name: "limit", Run: i.setItemsLimit},
			ken.SubCommandHandler{Name: "luck", Run: i.setLuckLevel},
		}},
	)
}

func (i *Inventory) get(ctx ken.SubCommandContext) (err error) {
	ctx.SetEphemeral(true)

	player := ctx.Options().GetByName("user").UserValue(ctx)
	inv, err := i.models.Inventories.GetByDiscordID(player.ID)
	if err != nil {
		discord.ErrorMessage(
			ctx,
			"Failed to Find Inventory",
			fmt.Sprintf("Are you sure there's an inventory for %s?", player.Username),
		)
		return err
	}

	allowed := i.inventoryAuthorized(ctx, inv)

	if !allowed {
		ctx.SetEphemeral(true)
		err = discord.ErrorMessage(ctx, "Unauthorized",
			"You are not authorized to use this command.")
		ctx.SetEphemeral(false)
		return err
	}

	host := discord.IsAdminRole(ctx, discord.AdminRoles...)
	embd := InventoryEmbedBuilder(inv, host)
	err = ctx.RespondEmbed(embd)
	if err != nil {
		return err
	}
	return nil
}

func (i *Inventory) delete(ctx ken.SubCommandContext) (err error) {
	authed := discord.IsAdminRole(ctx, discord.AdminRoles...)
	if !authed {
		err = discord.ErrorMessage(
			ctx,
			"Unauthorized",
			"You are not authorized to use this command.",
		)
		return err
	}

	userArg := ctx.Options().GetByName("user").UserValue(ctx)
	inv, err := i.models.Inventories.GetByDiscordID(userArg.ID)
	if err != nil {
		log.Println(err)
		return discord.ErrorMessage(ctx, "Failed to Find Inventory",
			fmt.Sprintf("Failed to find inventory for %s", userArg.Username))
	}
	sesh := ctx.GetSession()
	err = sesh.ChannelMessageDelete(inv.UserPinChannel, inv.UserPinMessage)
	if err != nil {
		channel, _ := sesh.Channel(inv.UserPinChannel)
		return discord.ErrorMessage(ctx, "Failed to Delete Message",
			fmt.Sprintf("Failed to delete message for %s, could not find message in channel %s",
				userArg.Username, channel.Name))
	}
	err = i.models.Inventories.Delete(userArg.ID)
	if err != nil {
		return discord.ErrorMessage(ctx, "Failed to Delete Inventory",
			fmt.Sprintf("Failed to delete inventory for %s", userArg.Username))
	}
	return discord.SuccessfulMessage(
		ctx,
		"Inventory Deleted",
		fmt.Sprintf("Removed inventory for channel %s", discord.MentionChannel(inv.UserPinChannel)),
	)
}

// Version implements ken.SlashCommand.
func (*Inventory) Version() string {
	return "1.0.0"
}

// In order to use the inventory channel you must meet one of the following criteria:
// 1. Call inventory command in confessional channel
// 2. Have the role "Host", "Co-Host", or "Bot Developer" AND
//   - Be in the same channel as the pinned inventory message
//   - Be within a whiteilsted channel (admin only channel...etc)
func (i *Inventory) inventoryAuthorized(ctx ken.SubCommandContext, inv *data.Inventory) bool {
	event := ctx.GetEvent()
	invokeChannelID := event.ChannelID
	invoker := event.Member

	// Base case of user is in confessional channel and is the owner of the inventory
	if inv.DiscordID == invoker.User.ID && inv.UserPinChannel == invokeChannelID {
		return true
	}

	// If not in confessional channel, check if in whitelist
	whitelistChannels, _ := i.models.Whitelists.GetAll()
	if invokeChannelID != inv.UserPinChannel {
		for _, whitelist := range whitelistChannels {
			if whitelist.ChannelID == invokeChannelID {
				return true
			}
		}
		return false
	}

	// Go through and make sure user has one of the allowed roles:
	for _, role := range invoker.Roles {
		for _, allowedRole := range discord.AdminRoles {
			if role == allowedRole {
				return true
			}
		}
	}
	return true
}

func (i *Inventory) listWhitelist(ctx ken.SubCommandContext) (err error) {
	wishlists, _ := i.models.Whitelists.GetAll()
	if len(wishlists) == 0 {
		err = discord.ErrorMessage(ctx, "No whitelisted channels", "Nothing here...")
		return err
	}

	fields := []*discordgo.MessageEmbedField{}
	for _, v := range wishlists {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   v.ChannelName,
			Inline: false,
		})
	}
	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Whitelisted Channels",
		Description: "Whitelisted channels for inventory",
		Fields:      fields,
	})
	return err
}

// Helper to handle getting the pinned message for inventory and updating it
func (i *Inventory) updateInventoryMessage(
	ctx ken.SubCommandContext,
	inventory *data.Inventory,
) (err error) {
	sesh := ctx.GetSession()
	_, err = sesh.ChannelMessageEditEmbed(
		inventory.UserPinChannel,
		inventory.UserPinMessage,
		InventoryEmbedBuilder(inventory, false),
	)
	if err != nil {
		log.Println(err)
		return discord.ErrorMessage(
			ctx,
			"Failed to update inventory message",
			"Alex is a bad programmer, and this is his fault.",
		)
	}
	return nil
}

func UpdateInventoryMessage(ctx ken.Context, i *data.Inventory) (err error) {
	sesh := ctx.GetSession()
	_, err = sesh.ChannelMessageEditEmbed(
		i.UserPinChannel,
		i.UserPinMessage,
		InventoryEmbedBuilder(i, false),
	)
	if err != nil {
		return err
	}
	return nil
}

// Helper to determine if user is authorized to use inventory command based on:
// 1. In their confessional (and owner of inventory)
// 2. In a whitelisted channel (and an admin)
func InventoryAuthorized(
	ctx ken.SubCommandContext,
	inv *data.Inventory,
	wl []*data.Whitelist,
) bool {
	event := ctx.GetEvent()
	invokeChannelID := event.ChannelID
	invoker := event.Member

	// Base case of user is in confessional channel and is the owner of the inventory
	if inv.DiscordID == invoker.User.ID && inv.UserPinChannel == invokeChannelID {
		return true
	}

	// If not in confessional channel, check if in whitelist
	if invokeChannelID != inv.UserPinChannel {
		for _, whitelist := range wl {
			if whitelist.ChannelID == invokeChannelID {
				return true
			}
		}
		return false
	}

	// Go through and make sure user has one of the allowed roles:
	for _, role := range invoker.Roles {
		for _, allowedRole := range discord.AdminRoles {
			if role == allowedRole {
				return true
			}
		}
	}
	return true
}

// Helper to attempt to fetch given user's inventory from user command option
func Fetch(ctx ken.SubCommandContext, m data.Models, adminOnly bool) (inv *data.Inventory, err error) {
	if adminOnly && !discord.IsAdminRole(ctx, discord.AdminRoles...) {
		return nil, ErrNotAuthorized
	}
	userArg, ok := ctx.Options().GetByNameOptional("user")
	event := ctx.GetEvent()
	channelID := event.ChannelID
	if !ok {
		inv, err = m.Inventories.GetByPinChannel(channelID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	if inv == nil {
		inv, err = m.Inventories.GetByDiscordID(userArg.UserValue(ctx).ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	wl, err := m.Whitelists.GetAll()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !InventoryAuthorized(ctx, inv, wl) {
		return nil, ErrNotAuthorized
	}
	if inv == nil {
		return nil, errors.New("somehow inventory is nil in middleware")
	}
	return inv, nil
}

// Helper to build inventory embed message based off if user is host or not
func InventoryEmbedBuilder(
	inv *data.Inventory,
	host bool,
) *discordgo.MessageEmbed {
	roleField := &discordgo.MessageEmbedField{
		Name:   "Role",
		Value:  inv.RoleName,
		Inline: true,
	}
	alignmentEmoji := discord.EmojiAlignment
	alignmentField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Alignment", alignmentEmoji),
		Value:  inv.Alignment,
		Inline: true,
	}

	// show coin bonus x100
	cb := inv.CoinBonus * 100
	coinStr := fmt.Sprintf("%d", inv.Coins) + " [" + fmt.Sprintf("%.2f", cb) + "%]"
	coinField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Coins", discord.EmojiCoins),
		Value:  coinStr,
		Inline: true,
	}
	abilitiesField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Abilities", discord.EmojiAbility),
		Value:  strings.Join(inv.Abilities, "\n"),
		Inline: true,
	}
	perksField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Perks", discord.EmojiPerk),
		Value:  strings.Join(inv.Perks, "\n"),
		Inline: true,
	}
	anyAbilitiesField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Any Abilities", discord.EmojiAnyAbility),
		Value:  strings.Join(inv.AnyAbilities, "\n"),
		Inline: true,
	}
	itemsField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Items (%d/%d)", discord.EmojiItem, len(inv.Items), inv.ItemLimit),
		Value:  strings.Join(inv.Items, "\n"),
		Inline: true,
	}
	statusesField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Statuses", discord.EmojiStatus),
		Value:  strings.Join(inv.Statuses, "\n"),
		Inline: true,
	}

	immunitiesField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Immunities", discord.EmojiImmunity),
		Value:  strings.Join(inv.Immunities, "\n"),
		Inline: true,
	}
	effectsField := &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Effects", discord.EmojiEffect),
		Value:  strings.Join(inv.Effects, "\n"),
		Inline: true,
	}
	isAlive := ""
	if inv.IsAlive {
		isAlive = fmt.Sprintf("%s Alive", discord.EmojiAlive)
	} else {
		isAlive = fmt.Sprintf("%s Dead", discord.EmojiDead)
	}

	deadField := &discordgo.MessageEmbedField{
		Name:   isAlive,
		Inline: true,
	}

	embd := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Inventory %s", discord.EmojiInventory),
		Fields: []*discordgo.MessageEmbedField{
			roleField,
			alignmentField,
			coinField,
			abilitiesField,
			anyAbilitiesField,
			perksField,
			itemsField,
			statusesField,
			immunitiesField,
			effectsField,
			deadField,
		},
		Color: discord.ColorThemeDiamond,
	}

	humanReqTime := util.GetEstTimeStamp()
	embd.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Last updated: %s", humanReqTime),
	}

	if host {

		embd.Fields = append(embd.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s Luck", discord.EmojiLuck),
			Value:  fmt.Sprintf("%d", inv.Luck),
			Inline: true,
		})

		noteListString := ""
		for i, note := range inv.Notes {
			noteListString += fmt.Sprintf("%d. %s\n", i+1, note)
		}

		embd.Fields = append(embd.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s Notes", discord.EmojiNote),
			Value:  noteListString,
			Inline: false,
		})

		embd.Color = discord.ColorThemeAmethyst

	}

	return embd
}

// Ability strings follow the format of 'Name [#]'
func ParseAbilityString(raw string) (name string, charges int, err error) {
	// Check if there's a charge amount
	charges = 1
	split := strings.Split(raw, " ")
	if len(split) > 1 {
		charges, err = strconv.Atoi(split[len(split)-1])
		if err != nil {
			return "", 0, err
		}
	}
	name = strings.Join(split[:len(split)-1], " ")
	return name, charges, nil
}
