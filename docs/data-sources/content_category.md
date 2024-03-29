---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Data Source pfptmeta_content_category - terraform-provider-pfptmeta"
subcategory: "Web Security Resources"
description: |-
  Web content categories allow for tracking and regulating access to websites according to their content type to ensure acceptable use. These categories include websites that are classified as pornography, gambling, entertainment etc. Once defined, the categories are added to access rules, enabling the administrators to manage users' access to websites within these categories. For administrator convenience, the solution provides three default content categories (permissive, moderate and strict) which include pre-defined sets of website types. In addition to creating a category-based object, the administrator can define specific URLs that may trigger a rule violation.
---

# Data Source (pfptmeta_content_category)

Web content categories allow for tracking and regulating access to websites according to their content type to ensure acceptable use. These categories include websites that are classified as pornography, gambling, entertainment etc. Once defined, the categories are added to access rules, enabling the administrators to manage users' access to websites within these categories. For administrator convenience, the solution provides three default content categories (permissive, moderate and strict) which include pre-defined sets of website types. In addition to creating a category-based object, the administrator can define specific URLs that may trigger a rule violation.

## Example Usage

```terraform
data "pfptmeta_content_category" "cc" {
  id = "cc-123abc"
}

output "content_category" {
  value = data.pfptmeta_content_category.cc
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `confidence_level` (String) ENUM: `LOW`, `MEDIUM`, `HIGH`.The classification engines classify URLs under certain categories with some degree of confidence based on various factors. The higher this confidence value is, the more certain is the engine in stating that the URL is indeed classified under that content type.
- `description` (String)
- `forbid_uncategorized_urls` (Boolean) Whether to forbid access to uncategorized URLs.
- `id` (String) The ID of this resource.
- `name` (String)
- `types` (List of String) Enum:`Abortion`, `Abused Drugs`, `Adult and Pornography`, `Alcohol and Tobacco`, `Auctions`, `Business and Economy`, `Cheating`, `Computer and Internet Info`, `Computer and Internet Security`, `Content Delivery Networks`, `Cult and Occult`, `Dating`, `Dead Sites`, `Dynamically Generated Content`, `Educational Institutions`, `Entertainment and Arts`, `Fashion and Beauty`, `Financial Services`, `Gambling`, `Games`, `Government`, `Gross`, `Hacking`, `Hate and Racism`, `Health and Medicine`, `Home and Garden`, `Hunting and Fishing`, `Illegal`, `Image and Video Search`, `Individual Stock Advice and Tools`, `Internet Portals`, `Internet Communications`, `Job Search`, `Kids`, `Legal`, `Local Information`, `Marijuana`, `Military`, `Motor Vehicles`, `Music`, `News and Media`, `Nudity`, `Online Greeting Cards`, `Parked Domains`, `Pay to Surf`, `Personal sites and Blogs`, `Personal Storage`, `Philosophy and Political Advocacy`, `Questionable`, `Real Estate`, `Recreation and Hobbies`, `Reference and Research`, `Religion`, `Search Engines`, `Sex Education`, `Shareware and Freeware`, `Shopping`, `Social Networking`, `Society`, `Sports`, `Streaming Media`, `Swimsuits and Intimate Apparel`, `Training and Tools`, `Translation`, `Travel`, `Violence`, `Weapons`, `Web Advertisements`, `Web-based Email`, `Web Hosting`
- `urls` (List of String) A list of URLs to put under this custom content category.
