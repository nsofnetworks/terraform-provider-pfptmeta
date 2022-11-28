package content_category

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id"}

const (
	description = "Web content categories allow for tracking and regulating access to websites according " +
		"to their content type to ensure acceptable use. " +
		"These categories include websites that are classified as pornography, gambling, entertainment etc. " +
		"Once defined, the categories are added to access rules, enabling the administrators to manage users' " +
		"access to websites within these categories. " +
		"For administrator convenience, the solution provides three default content categories " +
		"(permissive, moderate and strict) which include pre-defined sets of website types. " +
		"In addition to creating a category-based object, " +
		"the administrator can define specific URLs that may trigger a rule violation."
	confidenceLevelDesc = "ENUM: `LOW`, `MEDIUM`, `HIGH`." +
		"The classification engines classify URLs under certain categories with some degree of confidence based on various factors. " +
		"The higher this confidence value is, the more certain is the engine in stating that the URL is indeed classified under that content type."
	forbidUncategorizedUrlDesc = "Whether to forbid access to uncategorized URLs."
	typesDesc                  = "Enum:`Abortion`,`AbusedDrugs`,`AdultandPornography`,`AlcoholandTobacco`," +
		"`Auctions`,`BusinessandEconomy`,`Cheating`,`ComputerandInternetInfo`,`ComputerandInternetSecurity`," +
		"`ContentDeliveryNetworks`,`CultandOccult`,`Dating`,`DeadSites`,`DynamicallyGeneratedContent`," +
		"`EducationalInstitutions`,`EntertainmentandArts`,`FashionandBeauty`,`FinancialServices`," +
		"`Gambling`,`Games`,`Government`,`Gross`,`Hacking`,`HateandRacism`,`HealthandMedicine`," +
		"`HomeandGarden`,`HuntingandFishing`,`Illegal`,`ImageandVideoSearch`," +
		"`IndividualStockAdviceandTools`,`InternetPortals`,`InternetCommunications`,`JobSearch`,`Kids`," +
		"`Legal`,`LocalInformation`,`Marijuana`,`Military`,`MotorVehicles`,`Music`,`NewsandMedia`,`Nudity`," +
		"`OnlineGreetingCards`,`ParkedDomains`,`PaytoSurf`,`PersonalsitesandBlogs`,`PersonalStorage`," +
		"`PhilosophyandPoliticalAdvocacy`,`Questionable`,`RealEstate`,`RecreationandHobbies`," +
		"`ReferenceandResearch`,`Religion`,`SearchEngines`,`SexEducation`,`SharewareandFreeware`,`Shopping`," +
		"`SocialNetworking`,`Society`,`Sports`,`StreamingMedia`,`SwimsuitsandIntimateApparel`," +
		"`TrainingandTools`,`Translation`,`Travel`,`Violence`,`Weapons`,`WebAdvertisements`,`Web-basedEmail`," +
		"`WebHosting`"
	urlsDesc = "A list of URLs to put under this custom content category."
)

func contentCategoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	cc, err := client.GetContentCategory(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing device_settings %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId(cc.ID)
	err = client.MapResponseToResource(cc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
func contentCategoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewContentCategory(d)
	cc, err := client.CreateContentCategory(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cc.ID)
	err = client.MapResponseToResource(cc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func contentCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewContentCategory(d)
	cc, err := client.UpdateContentCategory(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cc.ID)
	err = client.MapResponseToResource(cc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func contentCategoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteContentCategory(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
