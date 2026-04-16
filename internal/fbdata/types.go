package fbdata

// FBAccountData represents the full rich data for a Facebook Account
type FBAccountData struct {
	UID     string     `json:"-"` // Not serialized in the outer JSON since it's the key
	Info    FBInfo     `json:"info"`
	Friends []FBFriend `json:"friends"`
	Groups  []FBGroup  `json:"groups"`
}

// FBInfo holds personal info and secrets like cookies
type FBInfo struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Cookie   string `json:"cookie,omitempty"`
	FbDtsg   string `json:"fb_dtsg,omitempty"`
}

// FBFriend holds basic friend information
type FBFriend struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

// FBGroup holds basic group information
type FBGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
