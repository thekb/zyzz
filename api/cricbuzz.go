package api

import (
	"encoding/xml"
	"io"
)

type MatchesData struct {
	XMLName   xml.Name `xml:"mchdata" json:"-"`
	MatchInfo []MatchInfo  `xml:"match" json:"MatchInfo"`
}

type MatchDetails struct {
	XMLName   xml.Name `xml:"mchDetails" json:"-"`
	MatchInfo MatchInfo `xml:"match" json:"MatchInfo"`
}


type MatchInfo struct {
	XMLName     xml.Name `xml:"match" json:"-"`
	Id          string `xml:"id,attr" json:"-"`
	Type        string `xml:"type,attr" json:"Type"`
	Srs         string `xml:"srs,attr" json:"Srs"`
	MatchDesc   string `xml:"mchDesc,attr" json:"MatchDesc"`
	MatchNumber string `xml:"mnum,attr" json:"MatchNumber"`
	Hostcity    string `xml:"vcity,attr" json:"HostCity"`
	Hostcountry string `xml:"vcountry,attr" json:"HostCountry"`
	Ground      string `xml:"grnd,attr" json:"Ground"`
	DataPath    string `xml:"datapath,attr" json:"DataPath"`
	InngCnt     string   `xml:"inngCnt,attr" json:"InngCnt"`
	MatchState  MatchState `xml:"state" json:"MatchState"`
	Team        []Team `xml:"Tm" json:"Team"`
	Schedule    Schedule `xml:"Tme" json:"Schedule"`
	Score 		Score `xml:"mscr" json:"Score"`
}

type MatchState struct {
	XMLName    xml.Name `xml:"state" json:"-"`
	MatchhState   string `xml:"mchState,attr" json:"MatchState"`
	Status     string `xml:"status,attr" json:"Status"`
	TossWon    string `xml:"TW,attr" json:"TossWon"`
	Decision   string `xml:"decisn,attr" json:"Decision"`
	AddnStatus string `xml:"addnStatus,attr" json:"AddnStatus"`
	SplStatus  string `xml:"splStatus,attr" json:"SplStatus"`
}

type Team struct {
	XMLName xml.Name `xml:"Tm" json:"-"`
	Id      string `xml:"id,attr" json:"-"`
	Name    string `xml:"Name,attr" json:"Name"`
	SName   string `xml:"sName,attr" json:"SName"`
	Flag    string `xml:"flag,attr" json:"Flag"`
}

type Schedule struct {
	XMLName   xml.Name `xml:"Tme" json:"-"`
	StartTime string `xml:"stTme,attr" json:"StartTime"`
	EndDate   string `xml:"enddt,attr" json:"EndDate"`
}

type Score struct {
	XMLName       xml.Name `xml:"mscr" json:"-"`
	InningsDetail InningsDetail `xml:"inngsdetail" json:"InningsDetail"`
	BattingTeam   BattingTeam `xml:"btTm" json:"BattingTeam"`
	BowlingTeam   BowlingTeam `xml:"blgTm" json:"BowlingTeam"`
	Batsmen       []Batsmen `xml:"btsmn" json:"Batsmen"`
	Bowler        []Bowler `xml:"blrs" json:"Bowler"`
}

type InningsDetail struct {
	XMLName            xml.Name `xml:"inngsdetail" json:"-"`
	NoOfOvers          string `xml:"noofovers,attr" json:"noOfOvers"`
	RequiredRunRate    string `xml:"rrr,attr" json:"RequiredRunRate"`
	CurrentRunRate     string `xml:"crr,attr" json:"CurrentRunRate"`
	CurrentPartnership string `xml:"cprtshp,attr" json:"CurrentPartnership"`
}

type BattingTeam struct {
	XMLName xml.Name `xml:"btTm" json:"-"`
	Id      string `xml:"id,attr" json:"-"`
	SName   string `xml:"sName,attr" json:"SName"`
	Innings []Innings `xml:"Inngs" json:"Innings"`
}

type BowlingTeam struct {
	XMLName xml.Name `xml:"blgTm" json:"-"`
	Id      string `xml:"id,attr" json:"-"`
	SName   string `xml:"sName,attr" json:"SName"`
	Innings []Innings `xml:"Inngs" json:"Innings"`
}

type Innings struct {
	XMLName     xml.Name `xml:"Inngs" json:"-"`
	Description string `xml:"desc,attr" json:"Description"`
	Runs        string `xml:"r,attr" json:"Runs"`
	Declared    string `xml:"Decl,attr" json:"Declared"`
	FollowOn    string `xml:"FollowOn,attr" json:"FollowOn"`
	Overs       string `xml:"ovrs,attr" json:"Overs"`
	Wickets     string `xml:"wkts,attr" json:"Wickets"`
}

type Batsmen struct {
	XMLName xml.Name `xml:"btsmn" json:"-"`
	Id      string `xml:"id,attr" json:"-"`
	SName   string `xml:"sName,attr" json:"SName"`
	Runs    string `xml:"r,attr" json:"Runs"`
	Balls   string `xml:"b,attr" json:"Balls"`
	Fours   string `xml:"frs,attr" json:"Fours"`
	Sixes   string `xml:"sxs,attr" json:"Sixes"`
}

type Bowler struct {
	XMLName xml.Name `xml:"blrs" json:"-"`
	Id      string `xml:"id,attr" json:"-"`
	SName   string `xml:"sName,attr" json:"SNames"`
	Runs    string `xml:"r,attr" json:"Runs"`
	Wickets string `xml:"b,wkts" json:"Wickets"`
	Overs   string `xml:"ovrs" json:"Overs"`
	Maidens string `xml:"mdns" json:"Maidens"`
}


func ReadMatchesData(reader io.Reader) (MatchesData, error) {
	var MatchesData MatchesData
	if err := xml.NewDecoder(reader).Decode(&MatchesData); err != nil {
		//return nil, err
	}

	return MatchesData, nil
}

func ReadMatchData(reader io.Reader) (MatchDetails, error) {
	var MatchDetails MatchDetails
	if err := xml.NewDecoder(reader).Decode(&MatchDetails); err != nil {
		//return nil, err
	}

	return MatchDetails, nil
}






