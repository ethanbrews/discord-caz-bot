package bot

import (
	"context"
	"log"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/terrors"
)

var commands = []*discordgo.ApplicationCommand{
	&assumeCommand,
}

func MakeBot(authToken string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + authToken)
	if err != nil {
		return nil, terrors.Augment(err, "creating discordgo session", nil)
	}
	return discord, nil
}

func StartBot(ctx context.Context, session *discordgo.Session, guildID string) error {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Logged in", "username", s.State.User.Username, "discriminator", s.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		return terrors.Augment(err, "opening discord-go session", nil)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, v)
		if err != nil {
			return terrors.Augment(err, "creating command", map[string]string{
				"command": v.Name,
			})
		}
		registeredCommands[i] = cmd
	}

	defer session.Close()

	<-ctx.Done()

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, guildID, v.ID)
		if err != nil {
			return terrors.Augment(err, "deleting command", map[string]string{
				"command": v.Name,
			})
		}
	}

	slog.Info("Gracefully shutting down.")
	return nil
}
