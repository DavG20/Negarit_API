package friends

import (
	message "github.com/DavG20/Negarit_API/internal/pkg/Message"
)

type Friends struct {
	Friend_Id         string            `bson:"_id,omitempty"  json:"id,omitempty"`
	Friend_A_UserName string            `json:"friend_a_username"`
	Friend_B_UserName string            `json:"friend_b_username"`
	Message           []message.Message `json:"message,omitempty"`
	Block_By_A        bool              `json:"blocka"`
	Block_By_B        bool              `json:"blockb"`
}
