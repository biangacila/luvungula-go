package microservice

type ServiceCassOrm struct { // Object-Relational Mapping (ORM)
	DbHost  string
	DbName  string
	DbTable string

	Body            interface{}
	WhereConditions []struct {
		Key  string
		Val  interface{}
		Type string
	}
	UpdateConditions map[string]interface{}
	SelectFieldList  []string
	QueryOnly        bool
}

type ServiceSaveEventsCassandra struct {
	DbHost  string
	DbName  string
	DbTable string

	AppName string
	Topic   string
	Payload interface{}
}
