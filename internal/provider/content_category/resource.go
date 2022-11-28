package content_category

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: contentCategoryCreate,
		ReadContext:   contentCategoryRead,
		UpdateContext: contentCategoryUpdate,
		DeleteContext: contentCategoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"confidence_level": {
				Description:      confidenceLevelDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("LOW", "MEDIUM", "HIGH"),
			},
			"forbid_uncategorized_urls": {
				Description: forbidUncategorizedUrlDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"types": {
				Description: typesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("Abortion", "Abused Drugs",
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
						"Travel", "Violence", "Weapons", "Web Advertisements", "Web-based Email", "Web Hosting"),
				},
				Optional: true,
			},
			"urls": {
				Description: urlsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateCustomUrlOrIPV4(),
				},
				Optional: true,
			},
		},
	}
}
