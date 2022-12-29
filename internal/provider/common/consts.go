package common

var Countries = []string{"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO",
	"AQ", "AR", "AS", "AT", "AU", "AW", "AX", "AZ", "BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI",
	"BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS", "BT", "BV", "BW", "BY", "BZ", "CA", "CC", "CD",
	"CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN", "CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ",
	"DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK",
	"FM", "FO", "FR", "GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ",
	"GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM", "HN", "HR", "HT", "HU", "ID", "IE", "IL", "IM",
	"IN", "IO", "IQ", "IR", "IS", "IT", "JE", "JM", "JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN",
	"KP", "KR", "KW", "KY", "KZ", "LA", "LB", "LC", "LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY",
	"MA", "MC", "MD", "ME", "MF", "MG", "MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS",
	"MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP",
	"NR", "NU", "NZ", "OM", "PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT",
	"PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW", "SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI",
	"SJ", "SK", "SL", "SM", "SN", "SO", "SR", "SS", "ST", "SV", "SX", "SY", "SZ", "TC", "TD", "TF",
	"TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO", "TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM",
	"US", "UY", "UZ", "VA", "VC", "VE", "VG", "VI", "VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM",
	"ZW"}
var ContentTypes = []string{"Abortion", "Abused Drugs",
	"Adult and Pornography", "Alcohol and Tobacco", "Auctions", "Business and Economy",
	"Cheating", "Computer and Internet Info", "Computer and Internet Security",
	"Content Delivery Networks", "Cult and Occult", "Dating", "Dead Sites",
	"Dynamically Generated Content", "Educational Institutions", "Entertainment and Arts",
	"Fashion and Beauty", "Financial Services", "Gambling", "Games", "Government", "Gross",
	"Hacking", "Hate and Racism", "Health and Medicine", "Home and Garden", "Hunting and Fishing",
	"Illegal", "Image and Video Search", "Individual Stock Advice and Tools", "Internet Portals",
	"Internet Communications", "Job Search", "Kids", "Legal", "Local Information", "Marijuana",
	"Military", "Motor Vehicles", "Music", "News and Media", "Nudity", "Online Greeting Cards",
	"Parked Domains", "Pay to Surf", "Personal sites and Blogs", "Personal Storage",
	"Philosophy and Political Advocacy", "Questionable", "Real Estate", "Recreation and Hobbies",
	"Reference and Research", "Religion", "Search Engines", "Sex Education",
	"Shareware and Freeware", "Shopping", "Social Networking", "Society", "Sports",
	"Streaming Media", "Swimsuits and Intimate Apparel", "Training and Tools", "Translation",
	"Travel", "Violence", "Weapons", "Web Advertisements", "Web-based Email", "Web Hosting"}

const (
	CountriesDoc = "`AD`,`AE`,`AF`,`AG`,`AI`,`AL`,`AM`,`AO`,`AQ`,`AR`,`AS`,`AT`,`AU`,`AW`,`AX`,`AZ`,`BA`,`BB`,`BD`,`BE`,`BF`," +
		"`BG`,`BH`,`BI`,`BJ`,`BL`,`BM`,`BN`,`BO`,`BQ`,`BR`,`BS`,`BT`,`BV`,`BW`,`BY`,`BZ`,`CA`,`CC`,`CD`,`CF`,`CG`,`CH`," +
		"`CI`,`CK`,`CL`,`CM`,`CN`,`CO`,`CR`,`CU`,`CV`,`CW`,`CX`,`CY`,`CZ`,`DE`,`DJ`,`DK`,`DM`,`DO`,`DZ`,`EC`,`EE`,`EG`," +
		"`EH`,`ER`,`ES`,`ET`,`FI`,`FJ`,`FK`,`FM`,`FO`,`FR`,`GA`,`GB`,`GD`,`GE`,`GF`,`GG`,`GH`,`GI`,`GL`,`GM`,`GN`,`GP`," +
		"`GQ`,`GR`,`GS`,`GT`,`GU`,`GW`,`GY`,`HK`,`HM`,`HN`,`HR`,`HT`,`HU`,`ID`,`IE`,`IL`,`IM`,`IN`,`IO`,`IQ`,`IR`,`IS`," +
		"`IT`,`JE`,`JM`,`JO`,`JP`,`KE`,`KG`,`KH`,`KI`,`KM`,`KN`,`KP`,`KR`,`KW`,`KY`,`KZ`,`LA`,`LB`,`LC`,`LI`,`LK`,`LR`," +
		"`LS`,`LT`,`LU`,`LV`,`LY`,`MA`,`MC`,`MD`,`ME`,`MF`,`MG`,`MH`,`MK`,`ML`,`MM`,`MN`,`MO`,`MP`,`MQ`,`MR`,`MS`,`MT`," +
		"`MU`,`MV`,`MW`,`MX`,`MY`,`MZ`,`NA`,`NC`,`NE`,`NF`,`NG`,`NI`,`NL`,`NO`,`NP`,`NR`,`NU`,`NZ`,`OM`,`PA`,`PE`,`PF`," +
		"`PG`,`PH`,`PK`,`PL`,`PM`,`PN`,`PR`,`PS`,`PT`,`PW`,`PY`,`QA`,`RE`,`RO`,`RS`,`RU`,`RW`,`SA`,`SB`,`SC`,`SD`,`SE`," +
		"`SG`,`SH`,`SI`,`SJ`,`SK`,`SL`,`SM`,`SN`,`SO`,`SR`,`SS`,`ST`,`SV`,`SX`,`SY`,`SZ`,`TC`,`TD`,`TF`,`TG`,`TH`,`TJ`," +
		"`TK`,`TL`,`TM`,`TN`,`TO`,`TR`,`TT`,`TV`,`TW`,`TZ`,`UA`,`UG`,`UM`,`US`,`UY`,`UZ`,`VA`,`VC`,`VE`,`VG`,`VI`,`VN`," +
		"`VU`,`WF`,`WS`,`YE`,`YT`,`ZA`,`ZM`,`ZW`"
	ContentTypesDoc = "`Abortion`, `Abused Drugs`, " +
		"`Adult and Pornography`, `Alcohol and Tobacco`, `Auctions`, `Business and Economy`, " +
		"`Cheating`, `Computer and Internet Info`, `Computer and Internet Security`, " +
		"`Content Delivery Networks`, `Cult and Occult`, `Dating`, `Dead Sites`, " +
		"`Dynamically Generated Content`, `Educational Institutions`, `Entertainment and Arts`, " +
		"`Fashion and Beauty`, `Financial Services`, `Gambling`, `Games`, `Government`, `Gross`, " +
		"`Hacking`, `Hate and Racism`, `Health and Medicine`, `Home and Garden`, `Hunting and Fishing`, " +
		"`Illegal`, `Image and Video Search`, `Individual Stock Advice and Tools`, `Internet Portals`, " +
		"`Internet Communications`, `Job Search`, `Kids`, `Legal`, `Local Information`, `Marijuana`, " +
		"`Military`, `Motor Vehicles`, `Music`, `News and Media`, `Nudity`, `Online Greeting Cards`, " +
		"`Parked Domains`, `Pay to Surf`, `Personal sites and Blogs`, `Personal Storage`, " +
		"`Philosophy and Political Advocacy`, `Questionable`, `Real Estate`, `Recreation and Hobbies`, " +
		"`Reference and Research`, `Religion`, `Search Engines`, `Sex Education`, " +
		"`Shareware and Freeware`, `Shopping`, `Social Networking`, `Society`, `Sports`, " +
		"`Streaming Media`, `Swimsuits and Intimate Apparel`, `Training and Tools`, `Translation`, " +
		"`Travel`, `Violence`, `Weapons`, `Web Advertisements`, `Web-based Email`, `Web Hosting`"
)
