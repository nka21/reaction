package interfaces

import (
	"log"

	"reaction/internal/entities"
	"reaction/internal/usecases"

	"github.com/bwmarrin/discordgo"
)

// DiscordHandler - Discordイベントを処理するハンドラー
type DiscordHandler struct {
	transferUseCase *usecases.TransferMessageUseCase
	config          *entities.Config
}

// NewDiscordHandler - 新しいDiscordHandlerを作成する
// WHY: UseCase を注入することで、DiscordHandler の責務を分離する
func NewDiscordHandler(transferUseCase *usecases.TransferMessageUseCase, config *entities.Config) *DiscordHandler {
	return &DiscordHandler{
		transferUseCase: transferUseCase,
		config:          config,
	}
}

// HandleReactionAdd - リアクション追加イベントを処理する
func (h *DiscordHandler) HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// WHY: Bot自身のリアクションは処理しない（無限ループ防止）
	if r.UserID == s.State.User.ID {
		return
	}

	if !h.isTriggerReactionEmoji(r.Emoji) {
		return
	}

	originalMsg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("メッセージ取得に失敗: %v", err)
		return
	}

	// WHY: すでにトリガーリアクションが付与されている場合は転送しない（重複転送を防ぐ）
	reactionCount := h.getTriggerReactionCount(originalMsg)
	if reactionCount > 1 {
		log.Printf("メッセージ %s には既にトリガーリアクションが %d 個ついているため、転送をスキップします", r.MessageID, reactionCount)
		return
	}

	// 転送メッセージを作成し、転送メッセージIDを保存
	err = h.transferUseCase.TransferMessage(s, originalMsg)
	if err != nil {
		log.Printf("メッセージ転送に失敗: %v", err)
		return
	}
}

// HandleReactionRemove - リアクション削除イベントを処理する
func (h *DiscordHandler) HandleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	// WHY: Bot自身のリアクションは処理しない（無限ループ防止）
	if r.UserID == s.State.User.ID {
		return
	}

	if !h.isTriggerReactionEmoji(r.Emoji) {
		return
	}

	originalMsg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("メッセージ取得に失敗: %v", err)
		return
	}

	// WHY: トリガーリアクションが0個になった場合のみ転送メッセージを削除
	triggerReactionCount := h.getTriggerReactionCount(originalMsg)
	if triggerReactionCount > 0 {
		log.Printf("メッセージ %s にはまだトリガーリアクションが %d 個残っているため、削除をスキップします", r.MessageID, triggerReactionCount)
		return
	}

	// 転送メッセージを削除（保存してあった、転送メッセージIDをキーにして削除）
	err = h.transferUseCase.DeleteTransferredMessage(s, r.MessageID)
	if err != nil {
		log.Printf("転送メッセージの削除に失敗: %v", err)
		return
	}
}

func (h *DiscordHandler) isTriggerReactionEmoji(emoji discordgo.Emoji) bool {
	return emoji.ID == h.config.TriggerReactionEmoji
}

func (h *DiscordHandler) getTriggerReactionCount(msg *discordgo.Message) int {
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.ID == h.config.TriggerReactionEmoji {
			return reaction.Count
		}
	}
	return 0
}
