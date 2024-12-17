package selector

type ApiToken struct {
	Token string
}

func NewEmptyApiTokenSearch() ApiToken {
	return ApiToken{}
}


func (s ApiToken) SetToken(token string) ApiToken {
	s.Token = token
	return s
}

func (s ApiToken) HasToken() bool {
	return s.Token != ""
}