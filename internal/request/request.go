package request

type AgeUpdate struct {
	Age int `json:"age"`
}

type UserDelete struct {
	TargetId uint64 `json:"target_id"`
}

type FriendshipRequest struct {
	SourceId uint64 `json:"source_id"`
	TargetId uint64 `json:"target_id"`
}

type UserData struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []uint64 `json:"friends"`
}
