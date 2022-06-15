package localStructs

// publicAccess represents data about open access request.
type PublicAccess struct {
	JWTtoken    string      `json:"jwtToken"`
	Email       string      `json:"email"`
	IP          string      `json:"ip"`
	Time        int64       `json:"time"`
	LevelAccess AccessLevel `json:"levelAccess"`
}

// privateAccess represents data about private access request.
type PrivateAccess struct {
	JWTtoken    string      `json:"jwtToken"`
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	IP          string      `json:"ip"`
	Time        int64       `json:"time"`
	LevelAccess AccessLevel `json:"levelAccess"`
}

// accessLevel represents data about access level enum.
type AccessLevel int

const (
	Global AccessLevel = iota
	Level1
	Level2
	Level3
)

// publicAccessTokenRecords slice to seed record publicAccess data.
var PublicAccessTokenRecords = []PublicAccess{}
var PrivateAccessTokenRecords = []PrivateAccess{}
