package schemas

type MicrosoftAction string

const (
	MicrosoftOutlookEventsAction MicrosoftAction = "get_outlook_events"
	MicrosoftTeamGroup           MicrosoftAction = "modify_team_group"
)

type MicrosoftReaction string

const (
	MicrosoftMailReaction MicrosoftReaction = "send_mail"
)

type MicrosoftUserInfo struct {
	Mail        string `json:"mail"`
	DisplayName string `json:"displayName"`
}

type MicrosoftResponseToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type MicrosoftOutlookEventsOptions struct {
	Subject string `json:"subject"`
}

type MicrosoftTeamsGroupOptionsInfos struct {
	Name string `json:"name"`
}

type MicrosoftTeamsChatResponse struct {
	IsOld               bool   `json:"is_old"`
	Name                string `json:"name"`
	LastUpdatedDateTime string `json:"lastUpdatedDateTime"`
}

type MicrosoftOutlookEventsResponse struct {
	Value []struct {
		Subject string `json:"subject"`
	} `json:"value"`
}

type MicrosoftSendMailBodyOptions struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type MicrosoftSendMailAdressOptions struct {
	Address string `json:"address"`
}

type MicrosoftSendMailRecipientsOptions struct {
	EmailAdress MicrosoftSendMailAdressOptions `json:"emailAddress"`
}

type MicrosoftSendMailMainMessageOptions struct {
	Subject      string                               `json:"subject"`
	Body         MicrosoftSendMailBodyOptions         `json:"body"`
	ToRecipients []MicrosoftSendMailRecipientsOptions `json:"toRecipients"`
}

type MicrosoftSendMailOptions struct {
	Message         MicrosoftSendMailMainMessageOptions `json:"message"`
	SaveToSentItems string                              `json:"saveToSentItems"`
}

type MicrosoftSendMailMainMessageOptionsSchema struct {
	Subject string                       `json:"subject"`
	Body    MicrosoftSendMailBodyOptions `json:"body"`
	Address string                       `json:"address"`
}

type MicrosoftSendMailOptionsSchema struct {
	Message         MicrosoftSendMailMainMessageOptionsSchema `json:"message"`
	SaveToSentItems string                                    `json:"saveToSentItems"`
}

type MicrosoftTeams struct {
	Id       string `json:"id"`
	Topic    string `json:"topic"`
	ChatType string `json:"chatType"`
}

type MicrosoftTeamsResponse struct {
	Value []MicrosoftTeams `json:"value"`
}
