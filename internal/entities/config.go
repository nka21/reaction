package entities

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config - Discord botの設定を保持する
type Config struct {
	DiscordBotToken      string
	TransferChannelID    string
	TriggerReactionEmoji string
}

// LoadConfig - 環境変数から設定を読み込む
func LoadConfig() (*Config, error) {
	// WHY: 本番環境では.envファイルが存在せず環境変数を直接設定するため、エラーは無視する
	_ = godotenv.Load()

	discordBotToken := os.Getenv("DISCORD_BOT_TOKEN")
	if discordBotToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN が設定されていません")
	}

	transferChannelID := os.Getenv("TRANSFER_CHANNEL_ID")
	if transferChannelID == "" {
		return nil, fmt.Errorf("TRANSFER_CHANNEL_ID が設定されていません")
	}

	triggerReactionEmoji := os.Getenv("TRIGGER_REACTION_EMOJI")
	if triggerReactionEmoji == "" {
		return nil, fmt.Errorf("TRIGGER_REACTION_EMOJI が設定されていません")
	}

	return &Config{
		DiscordBotToken:      discordBotToken,
		TransferChannelID:    transferChannelID,
		TriggerReactionEmoji: triggerReactionEmoji,
	}, nil
}
