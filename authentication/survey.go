package authentication


type Collection struct{
	Exercise_sid string
	Instrument_id string
	Period string
}
type Metadata struct {
	User_id string
	Ru_ref string
}
type Data struct{
	Value []string
	Block_id string
	Answer_id string
	Group_id string
	Group_instance int
	Answer_instance int
}
type Survey struct{
	Tx_id string
	Type string
	Version string
	Origin string
	Survey_id string
	Flushed bool
	Collection Collection
	Submitted_at string
	Metadata Metadata
	Data []Data

}


