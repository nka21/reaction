package usecases

import (
	"log"
	"sync"

	"reaction/internal/entities"

	"github.com/bwmarrin/discordgo"
)

// TransferMessageUseCase - メッセージ転送のビジネスロジックを担当する
type TransferMessageUseCase struct {
	config             *entities.Config
	transferMsgMapping map[string]string // 元メッセージID と 転送メッセージID のマッピング
	mappingMutex       sync.RWMutex      // マッピング操作の排他制御（並行処理を防ぐ）
}

// NewTransferMessageUseCase - 新しいTransferMessageUseCaseを作成する
func NewTransferMessageUseCase(config *entities.Config) *TransferMessageUseCase {
	return &TransferMessageUseCase{
		config:             config,
		transferMsgMapping: make(map[string]string),
	}
}

// TransferMessage - 転送メッセージを作成し、転送メッセージIDを保存
func (uc *TransferMessageUseCase) TransferMessage(
	session *discordgo.Session,
	originalMsg *discordgo.Message,
) error {
	// WHY: Discord APIのメッセージ転送機能を使用するため、元メッセージへの参照を作成
	transferRef := originalMsg.Forward()

	transferMsgSend := &discordgo.MessageSend{
		Reference: transferRef,
	}

	// 指定されたチャンネルに転送
	transferredMsg, err := session.ChannelMessageSendComplex(uc.config.TransferChannelID, transferMsgSend)
	if err != nil {
		log.Printf("メッセージ転送に失敗: %v", err)
		return err
	}

	// 転送メッセージIDを保存
	uc.mappingMutex.Lock()
	uc.transferMsgMapping[originalMsg.ID] = transferredMsg.ID
	uc.mappingMutex.Unlock()

	log.Printf("メッセージ %s をチャンネル %s に転送しました (転送メッセージID: %s)", originalMsg.ID, uc.config.TransferChannelID, transferredMsg.ID)
	return nil
}

// DeleteTransferredMessage - 転送されたメッセージを削除する
func (uc *TransferMessageUseCase) DeleteTransferredMessage(
	session *discordgo.Session,
	originalMsgID string,
) error {
	// マッピングから転送メッセージIDを取得
	uc.mappingMutex.RLock()
	transferredMsgID, exists := uc.transferMsgMapping[originalMsgID]
	uc.mappingMutex.RUnlock()

	if !exists {
		log.Printf("メッセージ %s の転送記録が見つかりません", originalMsgID)
		return nil
	}

	// 転送メッセージを削除
	err := session.ChannelMessageDelete(uc.config.TransferChannelID, transferredMsgID)
	if err != nil {
		log.Printf("転送メッセージの削除に失敗: %v", err)
		return err
	}

	// マッピングから削除
	uc.mappingMutex.Lock()
	delete(uc.transferMsgMapping, originalMsgID)
	uc.mappingMutex.Unlock()

	log.Printf("転送メッセージ %s を削除しました", transferredMsgID)
	return nil
}

// IsTransferredMessage - 指定されたメッセージIDが転送メッセージかどうかを判定する
func (uc *TransferMessageUseCase) IsTransferredMessage(msgID string) bool {
	uc.mappingMutex.RLock()
	defer uc.mappingMutex.RUnlock()

	// WHY: transferMsgMappingの値として存在する場合、そのメッセージは転送メッセージ
	for _, transferredMsgID := range uc.transferMsgMapping {
		if transferredMsgID == msgID {
			return true
		}
	}
	return false
}
