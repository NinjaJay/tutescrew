package commands

import (
	"errors"
	"net/url"

	"github.com/albshin/tutescrew/config"
)

// Verify struct class
type Verify struct {
	Config config.CASConfig
}

func (r *Verify) handle(ctx Context) error {
	// Check if student is already verified
	ch, err := ctx.Sess.State.Channel(ctx.Msg.ChannelID)
	if err != nil {
		return err
	}
	g, _ := ctx.Sess.State.Guild(ch.GuildID)

	if UserIDHasRoleByGuild("Verified", ctx.Msg.Author.ID, g) {
		ctx.Sess.ChannelMessageSend(ctx.Msg.ChannelID, "You are already verified!")
		return errors.New("already verified")
	}

	// Build the full login URL
	u, err := url.Parse(r.Config.AuthURL)
	if err != nil {
		return err
	}
	q := u.Query()

	// Encode Discord values into the redirect
	re, err := url.Parse(r.Config.RedirectURL)
	if err != nil {
		return err
	}
	reque := re.Query()
	reque.Add("guild", ch.GuildID)
	reque.Add("discord_id", ctx.Msg.Author.ID)
	re.RawQuery = reque.Encode()

	// Add redirect to url
	q.Add("service", re.String())
	u.RawQuery = q.Encode()

	usrch, err := ctx.Sess.UserChannelCreate(ctx.Msg.Author.ID)
	if err != nil {
		return err
	}
	ctx.Sess.ChannelMessageSend(usrch.ID, "Please go to "+u.String()+" to start the verification process.")

	return nil
}

func (r *Verify) description() string {
	return "Allows the user to start the student validation process. Upon success, the user with receive the \"Verified\" role."
}
func (r *Verify) usage() string { return "" }
func (r *Verify) canDM() bool   { return false }
