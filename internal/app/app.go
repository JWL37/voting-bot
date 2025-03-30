package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"voting-bot/internal/config"
	"voting-bot/internal/handlers"
	"voting-bot/internal/storage/tarantool"

	"github.com/mattermost/mattermost-server/v6/model"
)

type App struct {
	MattermostClient          *model.Client4
	MattermostWebSocketClient *model.WebSocketClient
	MattermostChannel         *model.Channel
	MattermostTeam            *model.Team
	log                       *slog.Logger
	Storage                   *tarantool.Storage
}

func NewApp(log *slog.Logger, cfg *config.Config) *App {
	storage, err := tarantool.NewConnection(log, &cfg.Database)
	if err != nil {
		log.Error(fmt.Errorf("failed to create DB: %w", err).Error())
		os.Exit(1)
	}
	client := model.NewAPIv4Client(cfg.MattermostServer)
	client.SetToken(cfg.MattermostToken)

	team, _, err := client.GetTeamByName(cfg.MattermostTeamName, "")
	if err != nil {
		log.Error(fmt.Errorf("error fetching team: %v", err).Error())
	}

	channel, _, err := client.GetChannelByName(cfg.MattermostChannel, team.Id, "")
	if err != nil {
		log.Error(fmt.Errorf("error fetching channel: %v", err).Error())
	}

	wsURL := strings.Replace(cfg.MattermostServer, "http", "ws", 1)
	wsClient, err := model.NewWebSocketClient4(wsURL, client.AuthToken)
	if err != nil {
		log.Error(fmt.Errorf("error connecting to WebSocket: %v", err).Error())
	}

	return &App{
		MattermostClient:          client,
		MattermostWebSocketClient: wsClient,
		MattermostTeam:            team,
		MattermostChannel:         channel,
		log:                       log,
		Storage:                   storage,
	}
}

func (a *App) Run() error {

	a.log.Info("Starting bot")

	a.MattermostWebSocketClient.Listen()

	for event := range a.MattermostWebSocketClient.EventChannel {
		if event.EventType() == model.WebsocketEventPosted {
			postData, ok := event.GetData()["post"].(string)
			if !ok {
				continue
			}

			var post model.Post
			if err := json.Unmarshal([]byte(postData), &post); err != nil {
				a.log.Error(fmt.Errorf("error decoding post JSON: %v", err).Error())
				continue
			}

			if post.ChannelId == a.MattermostChannel.Id {
				a.SwitchCommands(&post)
			}

		}
	}
	return nil
}

func (a *App) Stop() {
	a.log.Info("Stopping bot")

	a.Storage.DB.Close()
	a.MattermostWebSocketClient.Close()
}

func (a *App) SwitchCommands(post *model.Post) {
	channel := a.MattermostChannel
	client := a.MattermostClient
	switch {
	case strings.HasPrefix(post.Message, "poll create"):
		{
			response := handlers.HandleCreateCommand(a.log, post.Message, post.UserId, a.Storage)
			client.CreatePost(&model.Post{
				ChannelId: channel.Id,
				Message:   response,
			})
		}
	case strings.HasPrefix(post.Message, "poll vote"):
		{
			response := handlers.HandleVoteCommand(a.log, post.Message, post.UserId, a.Storage)
			client.CreatePost(&model.Post{
				ChannelId: channel.Id,
				Message:   response,
			})
		}
	case strings.HasPrefix(post.Message, "poll result"):
		{
			response := handlers.HandleResultCommand(a.log, post.Message, a.Storage)
			client.CreatePost(&model.Post{
				ChannelId: channel.Id,
				Message:   response,
			})
		}
	case strings.HasPrefix(post.Message, "poll close"):
		{
			response := handlers.HandleCloseCommand(a.log, post.Message, post.UserId, a.Storage)
			client.CreatePost(&model.Post{
				ChannelId: channel.Id,
				Message:   response,
			})
		}
	case strings.HasPrefix(post.Message, "poll delete"):
		{
			response := handlers.HandleDeleteCommand(a.log, post.Message, post.UserId, a.Storage)
			client.CreatePost(&model.Post{
				ChannelId: channel.Id,
				Message:   response,
			})
		}
	}
}
