package client

type V2BaseMessage struct {
	Id            string
	Timestamp     string
	StreamId      string
	V2messageType string
}

type Attachment struct {
	// TODO
}

type V2Message struct {
	Message     string
	FromUserId  int64
	Attachments []Attachment
	StreamId    string
}

type UserJoinedRoomMessage struct {
	AddedByUserId     int64
	MemberAddedUserId int64
	Id                string
	Timestamp         string
	StreamId          string
}

type UserLeftRoomMessage struct {
	RemovedByUserId  int64
	MemberLeftUserId int64
	Id               string
	Timestamp        string
	StreamId         string
}

type RoomMemberPromotedToOwnerMessage struct {
	PromotedByUserId  int64
	PromotedUserId int64
	Id               string
	Timestamp        string
	StreamId         string
}

type RoomMemberDemotedFromOwnerMessage struct {
	DemotedByUserId  int64
	DemotedUserId int64
	Id               string
	Timestamp        string
	StreamId         string
}

