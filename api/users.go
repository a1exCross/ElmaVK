package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/a1exCross/ElmaVK/ApiErrors"
)

type UserParams struct {
	ID                     int    `json:"id"`
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	CanAccessClosed        bool   `json:"can_access_closed"`
	IsClosed               bool   `json:"is_closed"`
	Sex                    int    `json:"sex,omitempty"`
	ScreenName             string `json:"screen_name,omitempty"`
	Photo50                string `json:"photo_50,omitempty"`
	Photo100               string `json:"photo_100,omitempty"`
	Online                 int    `json:"online,omitempty"`
	Verified               int    `json:"verified,omitempty"`
	FriendStatus           int    `json:"friend_status,omitempty"`
	Nickname               string `json:"nickname,omitempty"`
	Domain                 string `json:"domain,omitempty"`
	Bdate                  string `json:"bdate,omitempty"`
	Photo200               string `json:"photo_200,omitempty"`
	PhotoMax               string `json:"photo_max,omitempty"`
	Photo200Orig           string `json:"photo_200_orig,omitempty"`
	Photo400Orig           string `json:"photo_400_orig,omitempty"`
	PhotoMaxOrig           string `json:"photo_max_orig,omitempty"`
	PhotoID                string `json:"photo_id,omitempty"`
	HasPhoto               int    `json:"has_photo,omitempty"`
	HasMobile              int    `json:"has_mobile,omitempty"`
	IsFriend               int    `json:"is_friend,omitempty"`
	CanPost                int    `json:"can_post,omitempty"`
	CanSeeAllPosts         int    `json:"can_see_all_posts,omitempty"`
	CanSeeAudio            int    `json:"can_see_audio,omitempty"`
	CanWritePrivateMessage int    `json:"can_write_private_message,omitempty"`
	CanSendFriendRequest   int    `json:"can_send_friend_request,omitempty"`
	CommonCount            int    `json:"common_count,omitempty"`
	Site                   string `json:"site,omitempty"`
	Status                 string `json:"status,omitempty"`
	LastSeen               *struct {
		Platform int `json:"platform,omitempty"`
		Time     int `json:"time,omitempty"`
	} `json:"last_seen,omitempty"`
	CropPhoto *struct {
		Photo *struct {
			AlbumID int `json:"album_id,omitempty"`
			Date    int `json:"date,omitempty"`
			ID      int `json:"id,omitempty"`
			OwnerID int `json:"owner_id,omitempty"`
			PostID  int `json:"post_id,omitempty"`
			Sizes   *[]struct {
				Height int    `json:"height,omitempty"`
				URL    string `json:"url,omitempty"`
				Type   string `json:"type,omitempty"`
				Width  int    `json:"width,omitempty"`
			} `json:"sizes,omitempty"`
			Text    string `json:"text,omitempty"`
			HasTags bool   `json:"has_tags,omitempty"`
		} `json:"photo,omitempty"`
		Crop *struct {
			X  float64 `json:"x,omitempty"`
			Y  float64 `json:"y,omitempty"`
			X2 float64 `json:"x2,omitempty"`
			Y2 float64 `json:"y2,omitempty"`
		} `json:"crop,omitempty"`
		Rect *struct {
			X  float64 `json:"x,omitempty"`
			Y  float64 `json:"y,omitempty"`
			X2 float64 `json:"x2,omitempty"`
			Y2 float64 `json:"y2,omitempty"`
		} `json:"rect,omitempty"`
	} `json:"crop_photo,omitempty"`
	FollowersCount   int `json:"followers_count,omitempty"`
	Blacklisted      int `json:"blacklisted,omitempty"`
	BlacklistedByMe  int `json:"blacklisted_by_me,omitempty"`
	IsHiddenFromFeed int `json:"is_hidden_from_feed,omitempty"`
	Occupation       *struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"occupation,omitempty"`
}

type User struct {
	Response []UserParams `json:"response"`
}

//https://dev.vk.com/reference/objects/user
type UserGetFields string

const (
	Activities             UserGetFields = "activities"
	About                  UserGetFields = "about"
	Blacklisted            UserGetFields = "blacklisted"
	BlacklistedByMe        UserGetFields = "blacklisted_by_me"
	Books                  UserGetFields = "books"
	Bdate                  UserGetFields = "bdate"
	CanBeInvitedGroup      UserGetFields = "can_be_invited_group"
	CanPost                UserGetFields = "can_post"
	CanSeeAllPosts         UserGetFields = "can_see_all_posts"
	CanSeeAudio            UserGetFields = "can_see_audio"
	CanSendFriendRequest   UserGetFields = "can_send_friend_request"
	CanWritePrivateMessage UserGetFields = "can_write_private_message"
	Career                 UserGetFields = "career"
	CommonCount            UserGetFields = "common_count"
	Connections            UserGetFields = "connections"
	Contacts               UserGetFields = "contacts"
	City                   UserGetFields = "city"
	Country                UserGetFields = "country"
	CropPhoto              UserGetFields = "crop_photo"
	Domain                 UserGetFields = "domain"
	Education              UserGetFields = "education"
	Exports                UserGetFields = "exports"
	FollowersCount         UserGetFields = "followers_count"
	FriendStatus           UserGetFields = "friend_status"
	HasPhoto               UserGetFields = "has_photo"
	HasMobile              UserGetFields = "has_mobile"
	HomeTown               UserGetFields = "home_town"
	Photo100               UserGetFields = "photo_100"
	Photo200               UserGetFields = "photo_200"
	Photo200Orig           UserGetFields = "photo_200_orig"
	Photo400Orig           UserGetFields = "photo_400_orig"
	Photo50                UserGetFields = "photo_50"
	Sex                    UserGetFields = "sex"
	Site                   UserGetFields = "site"
	Schools                UserGetFields = "schools"
	ScreenName             UserGetFields = "screen_name"
	Status                 UserGetFields = "status"
	Verified               UserGetFields = "verified"
	Games                  UserGetFields = "games"
	Interests              UserGetFields = "interests"
	IsFavorite             UserGetFields = "is_favorite"
	IsFriend               UserGetFields = "is_friend"
	IsHiddenFromFeed       UserGetFields = "is_hidden_from_feed"
	LastSeen               UserGetFields = "last_seen"
	MaidenName             UserGetFields = "maiden_name"
	Military               UserGetFields = "military"
	Movies                 UserGetFields = "movies"
	Music                  UserGetFields = "music"
	Nickname               UserGetFields = "nickname"
	Occupation             UserGetFields = "occupation"
	Online                 UserGetFields = "online"
	Personal               UserGetFields = "personal"
	PhotoID                UserGetFields = "photo_id"
	PhotoMax               UserGetFields = "photo_max"
	PhotoMaxOrig           UserGetFields = "photo_max_orig"
	Quotes                 UserGetFields = "quotes"
	Relation               UserGetFields = "relation"
	Relatives              UserGetFields = "relatives"
	Timezone               UserGetFields = "timezone"
	TV                     UserGetFields = "tv"
	Universities           UserGetFields = "universities"
)

type UserNameCases string

const (
	Nom UserNameCases = "nom" //именительный
	Gen UserNameCases = "gen" //родительный
	Dat UserNameCases = "dat" //дательный
	Acc UserNameCases = "acc" //винительный
	Ins UserNameCases = "ins" //творительный
	Abl UserNameCases = "abl" //предложный
)

type UserGetParams struct {
	UserIDs  []int
	Fields   []UserGetFields
	NameCase UserNameCases
}

func (v VK) GetUserByID(u UserGetParams) (User, error) {
	url := "access_token=" + v.Token + "&v=" + v.Version

	url += "&user_ids="
	for i, user := range u.UserIDs {
		if i > 0 {
			url += ","
		}
		url += strconv.Itoa(user)
	}

	url += "&fields="
	for i, f := range u.Fields {
		if i > 0 {
			url += ","
		}
		url += string(f)
	}

	url += "&name_case=" + string(u.NameCase)

	res, err := v.Reqeust_api_get("users.get?", url)

	if err != nil {
		return User{}, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return User{}, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return User{}, err
	}

	var user User

	err = json.Unmarshal(data, &user)

	if err != nil {
		return User{}, errors.New(err.Error())
	}

	return user, nil
}
