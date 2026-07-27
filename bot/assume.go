package bot

import "github.com/bwmarrin/discordgo"

var assumeCommand = discordgo.ApplicationCommand{
	Name:        "assume",
	Description: "Assume a role",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role-option",
			Description: "Role option",
			Required:    true,
		},
	},
}
