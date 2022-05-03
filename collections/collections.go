package collections

// https://app.quicktype.io/
type Collections []CollectionElement

type CollectionElement struct {
	ID                     string                  `json:"id"`
	CreatedAt              string                  `json:"created_at"`
	UpdatedAt              string                  `json:"updated_at"`
	Width                  int64                   `json:"width"`
	Height                 int64                   `json:"height"`
	Color                  string                  `json:"color"`
	BlurHash               string                  `json:"blur_hash"`
	Likes                  int64                   `json:"likes"`
	LikedByUser            bool                    `json:"liked_by_user"`
	Description            string                  `json:"description"`
	User                   User                    `json:"user"`
	CurrentUserCollections []CurrentUserCollection `json:"current_user_collections"`
	Urls                   Urls                    `json:"urls"`
	Links                  WelcomeLinks            `json:"links"`
}

type CurrentUserCollection struct {
	ID              int64       `json:"id"`
	Title           string      `json:"title"`
	PublishedAt     string      `json:"published_at"`
	LastCollectedAt string      `json:"last_collected_at"`
	UpdatedAt       string      `json:"updated_at"`
	CoverPhoto      interface{} `json:"cover_photo"`
	User            interface{} `json:"user"`
}

type WelcomeLinks struct {
	Self             string `json:"self"`
	HTML             string `json:"html"`
	Download         string `json:"download"`
	DownloadLocation string `json:"download_location"`
}

type Urls struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

type User struct {
	ID                string       `json:"id"`
	Username          string       `json:"username"`
	Name              string       `json:"name"`
	PortfolioURL      string       `json:"portfolio_url"`
	Bio               string       `json:"bio"`
	Location          string       `json:"location"`
	TotalLikes        int64        `json:"total_likes"`
	TotalPhotos       int64        `json:"total_photos"`
	TotalCollections  int64        `json:"total_collections"`
	InstagramUsername string       `json:"instagram_username"`
	TwitterUsername   string       `json:"twitter_username"`
	ProfileImage      ProfileImage `json:"profile_image"`
	Links             UserLinks    `json:"links"`
}

type UserLinks struct {
	Self      string `json:"self"`
	HTML      string `json:"html"`
	Photos    string `json:"photos"`
	Likes     string `json:"likes"`
	Portfolio string `json:"portfolio"`
}

type ProfileImage struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}
