package global

type FinInfo struct {
	Date      string
	WorkHour  float64
	TotHour   float64
	BonusHour float64
	Rate      float64
	Ph        bool
	Cost      float64
}
type WeekDate struct {
	Week      int
	StartDate string
	EndDate   string
	Mon       FinInfo
	Tue       FinInfo
	Wed       FinInfo
	Thu       FinInfo
	Fri       FinInfo
	Sat       FinInfo
	Sun       FinInfo
	TotOpen   int

	TotHours     float64
	TotHourNor   float64
	TotHourSat   float64
	TotHourSun   float64
	TotHourPh    float64
	TotHourOT    float64
	TotHourBonus float64

	TotHourRates   float64
	TotHourRateNor float64
	TotHourRateSat float64
	TotHourRateSun float64
	TotHourRatePh  float64
	TotHourRateOT  float64

	TotDays         float64
	TotDayNor       float64
	TotDaySat       float64
	TotDaySatCharge float64
	TotDaySun       float64
	TotDayPh        float64
	TotDayOT        float64

	TotDayRates   float64
	TotDayRateNor float64
	TotDayRateSat float64
	TotDayRateSun float64
	TotDayRatePh  float64
	TotDayRateOT  float64

	EmpNumber     string
	Site          string
	Designation   string
	Name          string
	PaymentMethod string
	BankName      string
	BankAccount   string
	BranchCode    string

	EmpIdNumber       string
	EmpGender         string
	EmpContractStart  string
	EmpContractStart2 string
	EmpContractEnd    string
	Gender            string

	RateNor float64
	RateSat float64
	RateSun float64
	RatePh  float64
	RateOT  float64
}
