package config

// AllowedSites is a map that represents a domain->hips_prefix
// of all image domain/prefix combinations we are allowing.
var AllowedSites map[string]string

// NormalizeSite normalizes site name.
func NormalizeSite(site string) string {
	// We have both "cosmo" and "cosmopolitan" saved in our database as possible
	// pointers to cosmopolitan image bucket, so we must account for both.
	if site == "cosmo" || site == "cosmopolitan" {
		return "cosmopolitan"
	}

	return site
}

func Init(cdn bool) {

	// LOGIC A:
	// If we see this path
	// https://hips.one-stage-us-east-1.hearstapps.net/cosmo/assets/16/16/1461452031-bra4.jpg

	// HIPS will look for the first record with a value of "cosmo" and use the key (which is the domain)
	// as the image source

	// LOGIC B:
	// If, however we see this path
	// https://hips.one-stage-us-east-1.hearstapps.net/cng.h-cdn.co/assets/16/16/1461452031-bra4.jpg

	// HIPS will look for the first record where the KEY is "cng.h-cdn.co" and assume this
	// is a valid domain and use this domain as the image source.

	// cdnSite is a map of all cdn domains to hips-path
	cdnSite := map[string]string{
		"bpc.h-cdn.co":              "bestproducts",
		"cad.h-cdn.co":              "caranddriver",
		"clv.h-cdn.co":              "countryliving",
		"cos.h-cdn.co":              "cosmopolitan",
		"del.h-cdn.co":              "delish",
		"doz.h-cdn.co":              "drozthegoodlife",
		"edc.h-cdn.co":              "elledecor",
		"ell.h-cdn.co":              "elle",
		"esq.h-cdn.co":              "esquire",
		"evr.h-cdn.co":              "everything",
		"ghk.h-cdn.co":              "goodhousekeeping",
		"hbu.h-cdn.co":              "housebeautiful",
		"hbz.h-cdn.co":              "harpersbazaar",
		"lnl.h-cdn.co":              "lennyletter",
		"mac.h-cdn.co":              "marieclaire",
		"mos.h-cdn.co":              "mediaos",
		"pop.h-cdn.co":              "popularmechanics",
		"rbk.h-cdn.co":              "redbook",
		"roa.h-cdn.co":              "roadandtrack",
		"sev.h-cdn.co":              "seventeen",
		"ssp.h-cdn.co":              "sharedspaces",
		"swt.h-cdn.co":              "sweet",
		"toc.h-cdn.co":              "townandcountry",
		"ver.h-cdn.co":              "veranda",
		"wdy.h-cdn.co":              "womansday",
		"ssu.h-cdn.co":              "sharedspaces-uk",
		"cng.h-cdn.co":              "cosmopolitan-ng",
		"cosmouk.cdnds.net":         "cosmopolitan-uk",
		"cin.h-cdn.co":              "cosmopolitan-in",
		"arv.h-cdn.co":              "arrevista",
		"dm.h-cdn.co":               "diezminutos",
		"ellees.h-cdn.co":           "elle-es",
		"wds.h-cdn.co":              "womansday-es",
		"cit.h-cdn.co":              "cosmopolitan-it",
		"elleit.h-cdn.co":           "elle-it",
		"gioit.h-cdn.co":            "gioia-it",
		"cjp.h-cdn.co":              "cosmopolitan-jp",
		"cnl.h.cdn.cosmopolitan.nl": "cosmopolitan-nl",
		"cnl.h-cdn.co":              "cosmopolitan-nl",
		"ellnl.h-cdn.co":            "elle-nl",
		"esqnl.h-cdn.co":            "esquire-nl",
		"hbznl.h-cdn.co":            "harpersbazaar-nl",
		"cdk.h-cdn.co":              "cosmopolitan-dk",
		"cno.h-cdn.co":              "cosmopolitan-no",
		"cse.h-cdn.co":              "cosmopolitan-se",
		"ctw.h-cdn.co":              "cosmopolitan-tw",
		"hbztw.h-cdn.co":            "harpersbazaar-tw",
	}

	// s3Site is a map of all s3 domains to hips-path
	s3Site := map[string]string{
		"amv-prod-bpc.s3.amazonaws.com":            "bestproducts",
		"amv-prod-cad.s3.amazonaws.com":            "caranddriver",
		"amv-prod-clv.s3.amazonaws.com":            "countryliving",
		"amv-prod-cos.s3.amazonaws.com":            "cosmopolitan",
		"amv-prod-del.s3.amazonaws.com":            "delish",
		"amv-prod-doz.s3.amazonaws.com":            "drozthegoodlife",
		"amv-prod-edc.s3.amazonaws.com":            "elledecor",
		"amv-prod-ell.s3.amazonaws.com":            "elle",
		"amv-prod-esq.s3.amazonaws.com":            "esquire",
		"amv-prod-evr.s3.amazonaws.com":            "everything",
		"amv-prod-ghk.s3.amazonaws.com":            "goodhousekeeping",
		"amv-prod-hbu.s3.amazonaws.com":            "housebeautiful",
		"amv-prod-hbz.s3.amazonaws.com":            "harpersbazaar",
		"amv-prod-lnl.s3.amazonaws.com":            "lennyletter",
		"amv-prod-mac.s3.amazonaws.com":            "marieclaire",
		"amv-prod-mos.s3.amazonaws.com":            "mediaos",
		"amv-prod-pop.s3.amazonaws.com":            "popularmechanics",
		"amv-prod-rbk.s3.amazonaws.com":            "redbook",
		"amv-prod-roa.s3.amazonaws.com":            "roadandtrack",
		"amv-prod-sev.s3.amazonaws.com":            "seventeen",
		"amv-prod-ssp.s3.amazonaws.com":            "sharedspaces",
		"amv-prod-swt.s3.amazonaws.com":            "sweet",
		"amv-prod-toc.s3.amazonaws.com":            "townandcountry",
		"amv-prod-ver.s3.amazonaws.com":            "veranda",
		"amv-prod-wdy.s3.amazonaws.com":            "womansday",
		"amv-prod-ssu.s3.amazonaws.com":            "sharedspaces-uk",
		"ame-prod-cng.s3.amazonaws.com":            "cosmopolitan-ng",
		"ame-prod-cosmouk-assets.s3.amazonaws.com": "cosmopolitan-uk",
		"amap-prod-cin.s3.amazonaws.com":           "cosmopolitan-in",
		"ame-prod-arv.s3.amazonaws.com":            "arrevista",
		"ame-prod-dm.s3.amazonaws.com":             "diezminutos",
		"ame-prod-ellees.s3.amazonaws.com":         "elle-es",
		"ame-prod-wds.s3.amazonaws.com":            "womansday-es",
		"ame-prod-cit.s3.amazonaws.com":            "cosmopolitan-it",
		"ame-prod-elleit.s3.amazonaws.com":         "elle-it",
		"ame-prod-gioit.s3.amazonaws.com":          "gioia-it",
		"amapn-prod-cjp.s3.amazonaws.com":          "cosmopolitan-jp",
		"ame-prod-cnl.s3.amazonaws.com":            "cosmopolitan-nl",
		"ame-prod-ellnl.s3.amazonaws.com":          "elle-nl",
		"ame-prod-esqnl.s3.amazonaws.com":          "esquire-nl",
		"ame-prod-hbznl.s3.amazonaws.com":          "harpersbazaar-nl",
		"ame-prod-cdk.s3.amazonaws.com":            "cosmopolitan-dk",
		"ame-prod-cno.s3.amazonaws.com":            "cosmopolitan-no",
		"ame-prod-cse.s3.amazonaws.com":            "cosmopolitan-se",
		"amapn-prod-ctw.s3.amazonaws.com":          "cosmopolitan-tw",
		"amapn-prod-hbztw.s3.amazonaws.com":        "harpersbazaar-tw",
		"amv-prod-cad-assets.s3.amazonaws.com":     "caranddriver-assets",
	}

	AllowedSites = map[string]string{
		// naming convention here are added so we are also conistent with interunity themes/brand naming.
		"amv-games-prod-assets.s3.amazonaws.com": "games-prod",
		"s3.amazonaws.com/hmg-prod":              "hmg-prod",
		"hmg-prod.s3.amazonaws.com":              "hmg-prod",
		"quizapp-assets.s3.amazonaws.com":        "quizatio",

		// dev and stage
		"amv-games-stage-assets.s3.amazonaws.com":             "games-stage",
		"s3.amazonaws.com/hmg-dev":                            "hmg-dev",
		"hmg-dev.s3.amazonaws.com":                            "hmg-dev",
		"hmg-test.s3.amazonaws.com":                           "hmg-test",
		"ame-stage-cit.s3.amazonaws.com":                      "cosmopolitan-it-stage",
		"amv-dev-cos.s3.amazonaws.com":                        "cosmopolitan-dev",
		"ctsf-dev.s3-website-us-east-1.amazonaws.com":         "ctsf-dev",
		"dev-rover-htvdev-mediaos-hearst-io.s3.amazonaws.com": "htv-profile",
		"mediaos.s3.amazonaws.com":                            "htv-qa",
		"htv-prod-media.s3.amazonaws.com":                     "htv-prod",
		"mp-eui-test.s3.amazonaws.com":                        "hdm-dev",
		"cdn-dev.quizatio.us":                                 "quizatio-dev",
		"stg-quizapp-assets-b.s3.amazonaws.com":               "quizatio-stage",
		"amv-stg-cad-assets.s3.amazonaws.com":                 "caranddriver-assets-stg",
		"amv-qa-cad-assets.s3.amazonaws.com":                  "caranddriver-assets-qa",

		// corptech
		"dev-thumb-out-mediaos-hearst-io.s3.amazonaws.com":          "vidthumb-dev",
		"dev.hearst-gopher.thumbs.s3.amazonaws.com":                 "vidthumb-devnew",
		"qa.hearst-gopher.thumbs.s3.amazonaws.com":                  "vidthumb-stage",
		"htvqa-thumb-out-htvdev-mediaos-hearst-io.s3.amazonaws.com": "vidthumb-htvqa",
		"hearst-gopher.thumbs.s3.amazonaws.com":                     "vidthumb",
		"dev-rover-media-hearst-io.s3.amazonaws.com":                "rover-dev",
		"stage-rover-media-hearst-io.s3.amazonaws.com":              "rover-stage",
		"prod-rover-media-hearst-io.s3.amazonaws.com":               "rover",

		// partner-feed.
		"partnerfeeds-stage.hdmtools.com": "partner-feed-stage",
		"partnerfeeds.hdmtools.com":       "partner-feed",
	}

	// if cdn is true, all s3 sites will have their value set as "_" ensuring only CDN gets picked.
	// else CDN sites gets the "_" as value, ensuring s3 gets picked.
	if cdn == true {
		for k2, _ := range s3Site {
			AllowedSites[k2] = "_"
		}

		for k, v := range cdnSite {
			AllowedSites[k] = v
		}
	} else {
		for k, _ := range cdnSite {
			AllowedSites[k] = "_"
		}

		for k2, v2 := range s3Site {
			AllowedSites[k2] = v2
		}
	}
}
