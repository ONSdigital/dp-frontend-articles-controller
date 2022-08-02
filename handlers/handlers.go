package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-renderer/helper"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// Bulletin handles bulletin requests
func SixteensBulletin(cfg config.Config, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		sixteensBulletin(w, r, accessToken, collectionID, lang, rc, zc, ac, cfg)
	})
}

func sixteensBulletin(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient, cfg config.Config) {
	ctx := req.Context()
	muxVars := mux.Vars(req)
	uri := muxVars["uri"]
	bulletin, err := ac.GetLegacyBulletin(ctx, userAccessToken, collectionID, lang, uri)

	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	breadcrumbs, err := zc.GetBreadcrumb(ctx, userAccessToken, collectionID, lang, bulletin.URI)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	model := mapper.CreateSixteensBulletinModel(basePage, *bulletin, breadcrumbs, lang)
	rc.BuildPage(w, model, "sixteens-bulletin")
}

// Bulletin handles bulletin requests
func Bulletin(cfg config.Config, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		bulletin(w, r, accessToken, collectionID, lang, rc, zc, ac, cfg)
	})
}

func bulletin(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient, cfg config.Config) {
	ctx := req.Context()
	bulletin, err := ac.GetLegacyBulletin(ctx, userAccessToken, collectionID, lang, req.URL.EscapedPath())

	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	/*
		type Bulletin struct {
			RelatedBulletins []Link      `json:"relatedBulletins"`
			Sections         []Section   `json:"sections"`
			Accordion        []Section   `json:"accordion"`
			RelatedData      []Link      `json:"relatedData"`
			Charts           []Figure    `json:"charts"`
			Tables           []Figure    `json:"tables"`
			Images           []Figure    `json:"images"`
			Equations        []Figure    `json:"equations"`
			Links            []Link      `json:"links"`
			Type             string      `json:"type"`
			URI              string      `json:"uri"`
			Description      Description `json:"description"`
			Versions         []Version   `json:"versions"`
			Alerts           []Alert     `json:"alerts"`
			LatestReleaseURI string      `json:"latestReleaseUri"`
		}

		type Section struct {
			Title    string `json:"title"`
			Markdown string `json:"markdown"`
		}

		babbage/src/main/web/templates/handlebars/partials/equation.handlebars
		babbage/src/main/java/com/github/onsdigital/babbage/template/handlebars/helpers/markdown/CustomMarkdownHelper.java
		babbage/src/main/java/com/github/onsdigital/babbage/template/handlebars/helpers/markdown/util/MathjaxTagReplacer.java

		Sidecar

		{
			"filename": "0369a98c",
			"content": "\\sqrt{x^2+1}",
			"files": [
				{
					"type": "generated-svg",
					"filename": "0369a98c.svg",
					"fileType": "svg"
				},
				{
					"type": "generated-mml",
					"filename": "0369a98c.mml",
					"fileType": "mml"
				},
				{
					"type": "generated-png",
					"filename": "0369a98c.png",
					"fileType": "png"
				}
			],
			"title": "something large",
			"type": "equation",
			"uri": "/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/0369a98c"
		}
	*/

	type SidecarFile struct {
		Type     string `json:"type"`
		Filename string `json:"filename"`
		FileType string `json:"fileType"`
	}

	type SidecarEquation struct {
		Type     string        `json:"type"`
		Filename string        `json:"filename"`
		Title    string        `json:"title"`
		Content  string        `json:"content"`
		Uri      string        `json:"uri"`
		Files    []SidecarFile `json:"files"`
	}

	type SidecarImage struct {
		Type     string        `json:"type"`
		Filename string        `json:"filename"`
		Title    string        `json:"title"`
		Subtitle string        `json:"subtitle"`
		Uri      string        `json:"uri"`
		Source   string        `json:"source"`
		Notes    string        `json:"notes"` /* Contains markdown */
		AltText  string        `json:"altText"`
		Files    []SidecarFile `json:"files"`
	}

	resolveEquations := func(markdown string) string {
		sidecarMap := make(map[string]string)
		sidecarMap["0369a98c"] = `{"filename":"0369a98c","content":"\\sqrt{x^2+1}","files":[{"type":"generated-svg","filename":"0369a98c.svg","fileType":"svg"},{"type":"generated-mml","filename":"0369a98c.mml","fileType":"mml"},{"type":"generated-png","filename":"0369a98c.png","fileType":"png"}],"title":"something large","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/0369a98c"}`
		sidecarMap["321d2877"] = `{"filename":"321d2877","content":"\\lim_{h \\rightarrow 0 } \\frac{f(x+h)-f(x)}{h}","files":[{"type":"generated-svg","filename":"321d2877.svg","fileType":"svg"},{"type":"generated-mml","filename":"321d2877.mml","fileType":"mml"},{"type":"generated-png","filename":"321d2877.png","fileType":"png"}],"title":"something huge","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/321d2877"}`
		sidecarMap["49eca6d1"] = `{"filename":"49eca6d1","content":"\\frac{1+\\frac{a}{b}}{1+\\frac{1}{1+\\frac{1}{a}}}","files":[{"type":"generated-svg","filename":"49eca6d1.svg","fileType":"svg"},{"type":"generated-mml","filename":"49eca6d1.mml","fileType":"mml"},{"type":"generated-png","filename":"49eca6d1.png","fileType":"png"}],"title":"something horrible","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/49eca6d1"}`
		sidecarMap["6220957a"] = `{"filename":"6220957a","content":"x^n","files":[{"type":"generated-svg","filename":"6220957a.svg","fileType":"svg"},{"type":"generated-mml","filename":"6220957a.mml","fileType":"mml"},{"type":"generated-png","filename":"6220957a.png","fileType":"png"}],"title":"something small","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/6220957a"}`
		sidecarMap["7fdf8aca"] = `{"filename":"7fdf8aca","content":"x^n \u003d y^n + z^n","files":[{"type":"generated-svg","filename":"7fdf8aca.svg","fileType":"svg"},{"type":"generated-mml","filename":"7fdf8aca.mml","fileType":"mml"},{"type":"generated-png","filename":"7fdf8aca.png","fileType":"png"}],"title":"something medium","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/7fdf8aca"}`
		sidecarMap["88d55720"] = `{"filename":"88d55720","content":"A_{m,n} \u003d \n \\begin{pmatrix}\n  a_{1,1} \u0026 a_{1,2} \u0026 \\cdots \u0026 a_{1,n} \\\\\n  a_{2,1} \u0026 a_{2,2} \u0026 \\cdots \u0026 a_{2,n} \\\\\n  \\vdots  \u0026 \\vdots  \u0026 \\ddots \u0026 \\vdots  \\\\\n  a_{m,1} \u0026 a_{m,2} \u0026 \\cdots \u0026 a_{m,n} \n \\end{pmatrix}","files":[{"type":"generated-svg","filename":"88d55720.svg","fileType":"svg"},{"type":"generated-mml","filename":"88d55720.mml","fileType":"mml"},{"type":"generated-png","filename":"88d55720.png","fileType":"png"}],"title":"Something bigger","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/88d55720"}`
		sidecarMap["a8fff4d1"] = `{"filename":"a8fff4d1","content":"P\\left(A\u003d2\\middle|\\frac{A^2}{B}\u003e4\\right)","files":[{"type":"generated-svg","filename":"a8fff4d1.svg","fileType":"svg"},{"type":"generated-mml","filename":"a8fff4d1.mml","fileType":"mml"},{"type":"generated-png","filename":"a8fff4d1.png","fileType":"png"}],"title":"something else","type":"equation","uri":"/economy/economicoutputandproductivity/output/articles/formulatest/formulatest/a8fff4d1"}`

		svgMap := make(map[string]string)
		svgMap["0369a98c"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="8.833ex" height="3ex" style="vertical-align: -0.667ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -1059.6 3771.5 1307.9" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E1-MJMATHI-78" d="M52 289Q59 331 106 386T222 442Q257 442 286 424T329 379Q371 442 430 442Q467 442 494 420T522 361Q522 332 508 314T481 292T458 288Q439 288 427 299T415 328Q415 374 465 391Q454 404 425 404Q412 404 406 402Q368 386 350 336Q290 115 290 78Q290 50 306 38T341 26Q378 26 414 59T463 140Q466 150 469 151T485 153H489Q504 153 504 145Q504 144 502 134Q486 77 440 33T333 -11Q263 -11 227 52Q186 -10 133 -10H127Q78 -10 57 16T35 71Q35 103 54 123T99 143Q142 143 142 101Q142 81 130 66T107 46T94 41L91 40Q91 39 97 36T113 29T132 26Q168 26 194 71Q203 87 217 139T245 247T261 313Q266 340 266 352Q266 380 251 392T217 404Q177 404 142 372T93 290Q91 281 88 280T72 278H58Q52 284 52 289Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-32" d="M109 429Q82 429 66 447T50 491Q50 562 103 614T235 666Q326 666 387 610T449 465Q449 422 429 383T381 315T301 241Q265 210 201 149L142 93L218 92Q375 92 385 97Q392 99 409 186V189H449V186Q448 183 436 95T421 3V0H50V19V31Q50 38 56 46T86 81Q115 113 136 137Q145 147 170 174T204 211T233 244T261 278T284 308T305 340T320 369T333 401T340 431T343 464Q343 527 309 573T212 619Q179 619 154 602T119 569T109 550Q109 549 114 549Q132 549 151 535T170 489Q170 464 154 447T109 429Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-2B" d="M56 237T56 250T70 270H369V420L370 570Q380 583 389 583Q402 583 409 568V270H707Q722 262 722 250T707 230H409V-68Q401 -82 391 -82H389H387Q375 -82 369 -68V230H70Q56 237 56 250Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-31" d="M213 578L200 573Q186 568 160 563T102 556H83V602H102Q149 604 189 617T245 641T273 663Q275 666 285 666Q294 666 302 660V361L303 61Q310 54 315 52T339 48T401 46H427V0H416Q395 3 257 3Q121 3 100 0H88V46H114Q136 46 152 46T177 47T193 50T201 52T207 57T213 61V578Z"></path>
		<path stroke-width="10" id="E1-MJSZ1-221A" d="M263 249Q264 249 315 130T417 -108T470 -228L725 302Q981 837 982 839Q989 850 1001 850Q1008 850 1013 844T1020 832V826L741 243Q645 43 540 -176Q479 -303 469 -324T453 -348Q449 -350 436 -350L424 -349L315 -96Q206 156 205 156L171 130Q138 104 137 104L111 130L263 249Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		 <use xlink:href="#E1-MJSZ1-221A" x="0" y="125"></use>
		<rect stroke="none" width="2766" height="60" x="1005" y="925"></rect>
		<g transform="translate(1005,0)">
		 <use xlink:href="#E1-MJMATHI-78" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="816" y="583"></use>
		 <use xlink:href="#E1-MJMAIN-2B" x="1256" y="0"></use>
		 <use xlink:href="#E1-MJMAIN-31" x="2261" y="0"></use>
		</g>
		</g>
		</svg>`
		svgMap["321d2877"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="20.667ex" height="5.5ex" style="vertical-align: -1.833ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -1553.3 8881.1 2364.4" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E2-MJMAIN-6C" d="M42 46H56Q95 46 103 60V68Q103 77 103 91T103 124T104 167T104 217T104 272T104 329Q104 366 104 407T104 482T104 542T103 586T103 603Q100 622 89 628T44 637H26V660Q26 683 28 683L38 684Q48 685 67 686T104 688Q121 689 141 690T171 693T182 694H185V379Q185 62 186 60Q190 52 198 49Q219 46 247 46H263V0H255L232 1Q209 2 183 2T145 3T107 3T57 1L34 0H26V46H42Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-69" d="M69 609Q69 637 87 653T131 669Q154 667 171 652T188 609Q188 579 171 564T129 549Q104 549 87 564T69 609ZM247 0Q232 3 143 3Q132 3 106 3T56 1L34 0H26V46H42Q70 46 91 49Q100 53 102 60T104 102V205V293Q104 345 102 359T88 378Q74 385 41 385H30V408Q30 431 32 431L42 432Q52 433 70 434T106 436Q123 437 142 438T171 441T182 442H185V62Q190 52 197 50T232 46H255V0H247Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-6D" d="M41 46H55Q94 46 102 60V68Q102 77 102 91T102 122T103 161T103 203Q103 234 103 269T102 328V351Q99 370 88 376T43 385H25V408Q25 431 27 431L37 432Q47 433 65 434T102 436Q119 437 138 438T167 441T178 442H181V402Q181 364 182 364T187 369T199 384T218 402T247 421T285 437Q305 442 336 442Q351 442 364 440T387 434T406 426T421 417T432 406T441 395T448 384T452 374T455 366L457 361L460 365Q463 369 466 373T475 384T488 397T503 410T523 422T546 432T572 439T603 442Q729 442 740 329Q741 322 741 190V104Q741 66 743 59T754 49Q775 46 803 46H819V0H811L788 1Q764 2 737 2T699 3Q596 3 587 0H579V46H595Q656 46 656 62Q657 64 657 200Q656 335 655 343Q649 371 635 385T611 402T585 404Q540 404 506 370Q479 343 472 315T464 232V168V108Q464 78 465 68T468 55T477 49Q498 46 526 46H542V0H534L510 1Q487 2 460 2T422 3Q319 3 310 0H302V46H318Q379 46 379 62Q380 64 380 200Q379 335 378 343Q372 371 358 385T334 402T308 404Q263 404 229 370Q202 343 195 315T187 232V168V108Q187 78 188 68T191 55T200 49Q221 46 249 46H265V0H257L234 1Q210 2 183 2T145 3Q42 3 33 0H25V46H41Z"></path>
		<path stroke-width="10" id="E2-MJMATHI-68" d="M137 683Q138 683 209 688T282 694Q294 694 294 685Q294 674 258 534Q220 386 220 383Q220 381 227 388Q288 442 357 442Q411 442 444 415T478 336Q478 285 440 178T402 50Q403 36 407 31T422 26Q450 26 474 56T513 138Q516 149 519 151T535 153Q555 153 555 145Q555 144 551 130Q535 71 500 33Q466 -10 419 -10H414Q367 -10 346 17T325 74Q325 90 361 192T398 345Q398 404 354 404H349Q266 404 205 306L198 293L164 158Q132 28 127 16Q114 -11 83 -11Q69 -11 59 -2T48 16Q48 30 121 320L195 616Q195 629 188 632T149 637H128Q122 643 122 645T124 664Q129 683 137 683Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-2192" d="M56 237T56 250T70 270H835Q719 357 692 493Q692 494 692 496T691 499Q691 511 708 511H711Q720 511 723 510T729 506T732 497T735 481T743 456Q765 389 816 336T935 261Q944 258 944 250Q944 244 939 241T915 231T877 212Q836 186 806 152T761 85T740 35T732 4Q730 -6 727 -8T711 -11Q691 -11 691 0Q691 7 696 25Q728 151 835 230H70Q56 237 56 250Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-30" d="M96 585Q152 666 249 666Q297 666 345 640T423 548Q460 465 460 320Q460 165 417 83Q397 41 362 16T301 -15T250 -22Q224 -22 198 -16T137 16T82 83Q39 165 39 320Q39 494 96 585ZM321 597Q291 629 250 629Q208 629 178 597Q153 571 145 525T137 333Q137 175 145 125T181 46Q209 16 250 16Q290 16 318 46Q347 76 354 130T362 333Q362 478 354 524T321 597Z"></path>
		<path stroke-width="10" id="E2-MJMATHI-66" d="M118 -162Q120 -162 124 -164T135 -167T147 -168Q160 -168 171 -155T187 -126Q197 -99 221 27T267 267T289 382V385H242Q195 385 192 387Q188 390 188 397L195 425Q197 430 203 430T250 431Q298 431 298 432Q298 434 307 482T319 540Q356 705 465 705Q502 703 526 683T550 630Q550 594 529 578T487 561Q443 561 443 603Q443 622 454 636T478 657L487 662Q471 668 457 668Q445 668 434 658T419 630Q412 601 403 552T387 469T380 433Q380 431 435 431Q480 431 487 430T498 424Q499 420 496 407T491 391Q489 386 482 386T428 385H372L349 263Q301 15 282 -47Q255 -132 212 -173Q175 -205 139 -205Q107 -205 81 -186T55 -132Q55 -95 76 -78T118 -61Q162 -61 162 -103Q162 -122 151 -136T127 -157L118 -162Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-28" d="M94 250Q94 319 104 381T127 488T164 576T202 643T244 695T277 729T302 750H315H319Q333 750 333 741Q333 738 316 720T275 667T226 581T184 443T167 250T184 58T225 -81T274 -167T316 -220T333 -241Q333 -250 318 -250H315H302L274 -226Q180 -141 137 -14T94 250Z"></path>
		<path stroke-width="10" id="E2-MJMATHI-78" d="M52 289Q59 331 106 386T222 442Q257 442 286 424T329 379Q371 442 430 442Q467 442 494 420T522 361Q522 332 508 314T481 292T458 288Q439 288 427 299T415 328Q415 374 465 391Q454 404 425 404Q412 404 406 402Q368 386 350 336Q290 115 290 78Q290 50 306 38T341 26Q378 26 414 59T463 140Q466 150 469 151T485 153H489Q504 153 504 145Q504 144 502 134Q486 77 440 33T333 -11Q263 -11 227 52Q186 -10 133 -10H127Q78 -10 57 16T35 71Q35 103 54 123T99 143Q142 143 142 101Q142 81 130 66T107 46T94 41L91 40Q91 39 97 36T113 29T132 26Q168 26 194 71Q203 87 217 139T245 247T261 313Q266 340 266 352Q266 380 251 392T217 404Q177 404 142 372T93 290Q91 281 88 280T72 278H58Q52 284 52 289Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-2B" d="M56 237T56 250T70 270H369V420L370 570Q380 583 389 583Q402 583 409 568V270H707Q722 262 722 250T707 230H409V-68Q401 -82 391 -82H389H387Q375 -82 369 -68V230H70Q56 237 56 250Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-29" d="M60 749L64 750Q69 750 74 750H86L114 726Q208 641 251 514T294 250Q294 182 284 119T261 12T224 -76T186 -143T145 -194T113 -227T90 -246Q87 -249 86 -250H74Q66 -250 63 -250T58 -247T55 -238Q56 -237 66 -225Q221 -64 221 250T66 725Q56 737 55 738Q55 746 60 749Z"></path>
		<path stroke-width="10" id="E2-MJMAIN-2212" d="M84 237T84 250T98 270H679Q694 262 694 250T679 230H98Q84 237 84 250Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		<g transform="translate(37,0)">
		 <use xlink:href="#E2-MJMAIN-6C"></use>
		 <use xlink:href="#E2-MJMAIN-69" x="283" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-6D" x="566" y="0"></use>
		</g>
		<g transform="translate(0,-675)">
		 <use transform="scale(0.707)" xlink:href="#E2-MJMATHI-68" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E2-MJMAIN-2192" x="581" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E2-MJMAIN-30" x="1586" y="0"></use>
		</g>
		<g transform="translate(1478,0)">
		<g transform="translate(286,0)">
		<rect stroke="none" width="6995" height="60" x="0" y="220"></rect>
		<g transform="translate(60,779)">
		 <use xlink:href="#E2-MJMATHI-66" x="0" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-28" x="555" y="0"></use>
		 <use xlink:href="#E2-MJMATHI-78" x="949" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-2B" x="1748" y="0"></use>
		 <use xlink:href="#E2-MJMATHI-68" x="2753" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-29" x="3334" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-2212" x="3950" y="0"></use>
		 <use xlink:href="#E2-MJMATHI-66" x="4955" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-28" x="5510" y="0"></use>
		 <use xlink:href="#E2-MJMATHI-78" x="5904" y="0"></use>
		 <use xlink:href="#E2-MJMAIN-29" x="6481" y="0"></use>
		</g>
		 <use xlink:href="#E2-MJMATHI-68" x="3207" y="-724"></use>
		</g>
		</g>
		</g>
		</svg>`
		svgMap["49eca6d1"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="9.333ex" height="8.667ex" style="vertical-align: -4.667ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -1768.1 4029.8 3728.9" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E1-MJMAIN-31" d="M213 578L200 573Q186 568 160 563T102 556H83V602H102Q149 604 189 617T245 641T273 663Q275 666 285 666Q294 666 302 660V361L303 61Q310 54 315 52T339 48T401 46H427V0H416Q395 3 257 3Q121 3 100 0H88V46H114Q136 46 152 46T177 47T193 50T201 52T207 57T213 61V578Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-2B" d="M56 237T56 250T70 270H369V420L370 570Q380 583 389 583Q402 583 409 568V270H707Q722 262 722 250T707 230H409V-68Q401 -82 391 -82H389H387Q375 -82 369 -68V230H70Q56 237 56 250Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-61" d="M33 157Q33 258 109 349T280 441Q331 441 370 392Q386 422 416 422Q429 422 439 414T449 394Q449 381 412 234T374 68Q374 43 381 35T402 26Q411 27 422 35Q443 55 463 131Q469 151 473 152Q475 153 483 153H487Q506 153 506 144Q506 138 501 117T481 63T449 13Q436 0 417 -8Q409 -10 393 -10Q359 -10 336 5T306 36L300 51Q299 52 296 50Q294 48 292 46Q233 -10 172 -10Q117 -10 75 30T33 157ZM351 328Q351 334 346 350T323 385T277 405Q242 405 210 374T160 293Q131 214 119 129Q119 126 119 118T118 106Q118 61 136 44T179 26Q217 26 254 59T298 110Q300 114 325 217T351 328Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-62" d="M73 647Q73 657 77 670T89 683Q90 683 161 688T234 694Q246 694 246 685T212 542Q204 508 195 472T180 418L176 399Q176 396 182 402Q231 442 283 442Q345 442 383 396T422 280Q422 169 343 79T173 -11Q123 -11 82 27T40 150V159Q40 180 48 217T97 414Q147 611 147 623T109 637Q104 637 101 637H96Q86 637 83 637T76 640T73 647ZM336 325V331Q336 405 275 405Q258 405 240 397T207 376T181 352T163 330L157 322L136 236Q114 150 114 114Q114 66 138 42Q154 26 178 26Q211 26 245 58Q270 81 285 114T318 219Q336 291 336 325Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		<g transform="translate(120,0)">
		<rect stroke="none" width="3789" height="60" x="0" y="220"></rect>
		<g transform="translate(659,976)">
		 <use xlink:href="#E1-MJMAIN-31" x="0" y="0"></use>
		 <use xlink:href="#E1-MJMAIN-2B" x="727" y="0"></use>
		<g transform="translate(1510,0)">
		<g transform="translate(342,0)">
		<rect stroke="none" width="497" height="60" x="0" y="220"></rect>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-61" x="84" y="648"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-62" x="134" y="-626"></use>
		</g>
		</g>
		</g>
		<g transform="translate(60,-950)">
		 <use xlink:href="#E1-MJMAIN-31" x="0" y="0"></use>
		 <use xlink:href="#E1-MJMAIN-2B" x="727" y="0"></use>
		<g transform="translate(1510,0)">
		<g transform="translate(342,0)">
		<rect stroke="none" width="1697" height="60" x="0" y="220"></rect>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="947" y="638"></use>
		<g transform="translate(60,-710)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2B" x="505" y="0"></use>
		<g transform="translate(910,0)">
		<g transform="translate(120,0)">
		<rect stroke="none" width="426" height="60" x="0" y="146"></rect>
		 <use transform="scale(0.574)" xlink:href="#E1-MJMAIN-31" x="119" y="656"></use>
		 <use transform="scale(0.574)" xlink:href="#E1-MJMATHI-61" x="104" y="-482"></use>
		</g>
		</g>
		</g>
		</g>
		</g>
		</g>
		</g>
		</g>
		</svg>`
		svgMap["6220957a"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="2.5ex" height="1.833ex" style="vertical-align: -0.167ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -747.1 1104.8 782.1" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E1-MJMATHI-78" d="M52 289Q59 331 106 386T222 442Q257 442 286 424T329 379Q371 442 430 442Q467 442 494 420T522 361Q522 332 508 314T481 292T458 288Q439 288 427 299T415 328Q415 374 465 391Q454 404 425 404Q412 404 406 402Q368 386 350 336Q290 115 290 78Q290 50 306 38T341 26Q378 26 414 59T463 140Q466 150 469 151T485 153H489Q504 153 504 145Q504 144 502 134Q486 77 440 33T333 -11Q263 -11 227 52Q186 -10 133 -10H127Q78 -10 57 16T35 71Q35 103 54 123T99 143Q142 143 142 101Q142 81 130 66T107 46T94 41L91 40Q91 39 97 36T113 29T132 26Q168 26 194 71Q203 87 217 139T245 247T261 313Q266 340 266 352Q266 380 251 392T217 404Q177 404 142 372T93 290Q91 281 88 280T72 278H58Q52 284 52 289Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-6E" d="M21 287Q22 293 24 303T36 341T56 388T89 425T135 442Q171 442 195 424T225 390T231 369Q231 367 232 367L243 378Q304 442 382 442Q436 442 469 415T503 336T465 179T427 52Q427 26 444 26Q450 26 453 27Q482 32 505 65T540 145Q542 153 560 153Q580 153 580 145Q580 144 576 130Q568 101 554 73T508 17T439 -10Q392 -10 371 17T350 73Q350 92 386 193T423 345Q423 404 379 404H374Q288 404 229 303L222 291L189 157Q156 26 151 16Q138 -11 108 -11Q95 -11 87 -5T76 7T74 17Q74 30 112 180T152 343Q153 348 153 366Q153 405 129 405Q91 405 66 305Q60 285 60 284Q58 278 41 278H27Q21 284 21 287Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		 <use xlink:href="#E1-MJMATHI-78" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="816" y="583"></use>
		</g>
		</svg>`
		svgMap["7fdf8aca"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="13.167ex" height="2.333ex" style="vertical-align: -0.5ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -747.1 5704.5 976.1" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E1-MJMATHI-78" d="M52 289Q59 331 106 386T222 442Q257 442 286 424T329 379Q371 442 430 442Q467 442 494 420T522 361Q522 332 508 314T481 292T458 288Q439 288 427 299T415 328Q415 374 465 391Q454 404 425 404Q412 404 406 402Q368 386 350 336Q290 115 290 78Q290 50 306 38T341 26Q378 26 414 59T463 140Q466 150 469 151T485 153H489Q504 153 504 145Q504 144 502 134Q486 77 440 33T333 -11Q263 -11 227 52Q186 -10 133 -10H127Q78 -10 57 16T35 71Q35 103 54 123T99 143Q142 143 142 101Q142 81 130 66T107 46T94 41L91 40Q91 39 97 36T113 29T132 26Q168 26 194 71Q203 87 217 139T245 247T261 313Q266 340 266 352Q266 380 251 392T217 404Q177 404 142 372T93 290Q91 281 88 280T72 278H58Q52 284 52 289Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-6E" d="M21 287Q22 293 24 303T36 341T56 388T89 425T135 442Q171 442 195 424T225 390T231 369Q231 367 232 367L243 378Q304 442 382 442Q436 442 469 415T503 336T465 179T427 52Q427 26 444 26Q450 26 453 27Q482 32 505 65T540 145Q542 153 560 153Q580 153 580 145Q580 144 576 130Q568 101 554 73T508 17T439 -10Q392 -10 371 17T350 73Q350 92 386 193T423 345Q423 404 379 404H374Q288 404 229 303L222 291L189 157Q156 26 151 16Q138 -11 108 -11Q95 -11 87 -5T76 7T74 17Q74 30 112 180T152 343Q153 348 153 366Q153 405 129 405Q91 405 66 305Q60 285 60 284Q58 278 41 278H27Q21 284 21 287Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-3D" d="M56 347Q56 360 70 367H707Q722 359 722 347Q722 336 708 328L390 327H72Q56 332 56 347ZM56 153Q56 168 72 173H708Q722 163 722 153Q722 140 707 133H70Q56 140 56 153Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-79" d="M21 287Q21 301 36 335T84 406T158 442Q199 442 224 419T250 355Q248 336 247 334Q247 331 231 288T198 191T182 105Q182 62 196 45T238 27Q261 27 281 38T312 61T339 94Q339 95 344 114T358 173T377 247Q415 397 419 404Q432 431 462 431Q475 431 483 424T494 412T496 403Q496 390 447 193T391 -23Q363 -106 294 -155T156 -205Q111 -205 77 -183T43 -117Q43 -95 50 -80T69 -58T89 -48T106 -45Q150 -45 150 -87Q150 -107 138 -122T115 -142T102 -147L99 -148Q101 -153 118 -160T152 -167H160Q177 -167 186 -165Q219 -156 247 -127T290 -65T313 -9T321 21L315 17Q309 13 296 6T270 -6Q250 -11 231 -11Q185 -11 150 11T104 82Q103 89 103 113Q103 170 138 262T173 379Q173 380 173 381Q173 390 173 393T169 400T158 404H154Q131 404 112 385T82 344T65 302T57 280Q55 278 41 278H27Q21 284 21 287Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-2B" d="M56 237T56 250T70 270H369V420L370 570Q380 583 389 583Q402 583 409 568V270H707Q722 262 722 250T707 230H409V-68Q401 -82 391 -82H389H387Q375 -82 369 -68V230H70Q56 237 56 250Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-7A" d="M347 338Q337 338 294 349T231 360Q211 360 197 356T174 346T162 335T155 324L153 320Q150 317 138 317Q117 317 117 325Q117 330 120 339Q133 378 163 406T229 440Q241 442 246 442Q271 442 291 425T329 392T367 375Q389 375 411 408T434 441Q435 442 449 442H462Q468 436 468 434Q468 430 463 420T449 399T432 377T418 358L411 349Q368 298 275 214T160 106L148 94L163 93Q185 93 227 82T290 71Q328 71 360 90T402 140Q406 149 409 151T424 153Q443 153 443 143Q443 138 442 134Q425 72 376 31T278 -11Q252 -11 232 6T193 40T155 57Q111 57 76 -3Q70 -11 59 -11H54H41Q35 -5 35 -2Q35 13 93 84Q132 129 225 214T340 322Q352 338 347 338Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		 <use xlink:href="#E1-MJMATHI-78" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="816" y="583"></use>
		 <use xlink:href="#E1-MJMAIN-3D" x="1382" y="0"></use>
		<g transform="translate(2443,0)">
		 <use xlink:href="#E1-MJMATHI-79" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="712" y="583"></use>
		</g>
		 <use xlink:href="#E1-MJMAIN-2B" x="3697" y="0"></use>
		<g transform="translate(4702,0)">
		 <use xlink:href="#E1-MJMATHI-7A" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="670" y="583"></use>
		</g>
		</g>
		</svg>`
		svgMap["88d55720"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="35.667ex" height="14.5ex" style="vertical-align: -6.667ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -3371.4 15333.8 6242.7" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E1-MJMATHI-41" d="M208 74Q208 50 254 46Q272 46 272 35Q272 34 270 22Q267 8 264 4T251 0Q249 0 239 0T205 1T141 2Q70 2 50 0H42Q35 7 35 11Q37 38 48 46H62Q132 49 164 96Q170 102 345 401T523 704Q530 716 547 716H555H572Q578 707 578 706L606 383Q634 60 636 57Q641 46 701 46Q726 46 726 36Q726 34 723 22Q720 7 718 4T704 0Q701 0 690 0T651 1T578 2Q484 2 455 0H443Q437 6 437 9T439 27Q443 40 445 43L449 46H469Q523 49 533 63L521 213H283L249 155Q208 86 208 74ZM516 260Q516 271 504 416T490 562L463 519Q447 492 400 412L310 260L413 259Q516 259 516 260Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-6D" d="M21 287Q22 293 24 303T36 341T56 388T88 425T132 442T175 435T205 417T221 395T229 376L231 369Q231 367 232 367L243 378Q303 442 384 442Q401 442 415 440T441 433T460 423T475 411T485 398T493 385T497 373T500 364T502 357L510 367Q573 442 659 442Q713 442 746 415T780 336Q780 285 742 178T704 50Q705 36 709 31T724 26Q752 26 776 56T815 138Q818 149 821 151T837 153Q857 153 857 145Q857 144 853 130Q845 101 831 73T785 17T716 -10Q669 -10 648 17T627 73Q627 92 663 193T700 345Q700 404 656 404H651Q565 404 506 303L499 291L466 157Q433 26 428 16Q415 -11 385 -11Q372 -11 364 -4T353 8T350 18Q350 29 384 161L420 307Q423 322 423 345Q423 404 379 404H374Q288 404 229 303L222 291L189 157Q156 26 151 16Q138 -11 108 -11Q95 -11 87 -5T76 7T74 17Q74 30 112 181Q151 335 151 342Q154 357 154 369Q154 405 129 405Q107 405 92 377T69 316T57 280Q55 278 41 278H27Q21 284 21 287Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-2C" d="M78 35T78 60T94 103T137 121Q165 121 187 96T210 8Q210 -27 201 -60T180 -117T154 -158T130 -185T117 -194Q113 -194 104 -185T95 -172Q95 -168 106 -156T131 -126T157 -76T173 -3V9L172 8Q170 7 167 6T161 3T152 1T140 0Q113 0 96 17Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-6E" d="M21 287Q22 293 24 303T36 341T56 388T89 425T135 442Q171 442 195 424T225 390T231 369Q231 367 232 367L243 378Q304 442 382 442Q436 442 469 415T503 336T465 179T427 52Q427 26 444 26Q450 26 453 27Q482 32 505 65T540 145Q542 153 560 153Q580 153 580 145Q580 144 576 130Q568 101 554 73T508 17T439 -10Q392 -10 371 17T350 73Q350 92 386 193T423 345Q423 404 379 404H374Q288 404 229 303L222 291L189 157Q156 26 151 16Q138 -11 108 -11Q95 -11 87 -5T76 7T74 17Q74 30 112 180T152 343Q153 348 153 366Q153 405 129 405Q91 405 66 305Q60 285 60 284Q58 278 41 278H27Q21 284 21 287Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-3D" d="M56 347Q56 360 70 367H707Q722 359 722 347Q722 336 708 328L390 327H72Q56 332 56 347ZM56 153Q56 168 72 173H708Q722 163 722 153Q722 140 707 133H70Q56 140 56 153Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-28" d="M94 250Q94 319 104 381T127 488T164 576T202 643T244 695T277 729T302 750H315H319Q333 750 333 741Q333 738 316 720T275 667T226 581T184 443T167 250T184 58T225 -81T274 -167T316 -220T333 -241Q333 -250 318 -250H315H302L274 -226Q180 -141 137 -14T94 250Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-61" d="M33 157Q33 258 109 349T280 441Q331 441 370 392Q386 422 416 422Q429 422 439 414T449 394Q449 381 412 234T374 68Q374 43 381 35T402 26Q411 27 422 35Q443 55 463 131Q469 151 473 152Q475 153 483 153H487Q506 153 506 144Q506 138 501 117T481 63T449 13Q436 0 417 -8Q409 -10 393 -10Q359 -10 336 5T306 36L300 51Q299 52 296 50Q294 48 292 46Q233 -10 172 -10Q117 -10 75 30T33 157ZM351 328Q351 334 346 350T323 385T277 405Q242 405 210 374T160 293Q131 214 119 129Q119 126 119 118T118 106Q118 61 136 44T179 26Q217 26 254 59T298 110Q300 114 325 217T351 328Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-31" d="M213 578L200 573Q186 568 160 563T102 556H83V602H102Q149 604 189 617T245 641T273 663Q275 666 285 666Q294 666 302 660V361L303 61Q310 54 315 52T339 48T401 46H427V0H416Q395 3 257 3Q121 3 100 0H88V46H114Q136 46 152 46T177 47T193 50T201 52T207 57T213 61V578Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-32" d="M109 429Q82 429 66 447T50 491Q50 562 103 614T235 666Q326 666 387 610T449 465Q449 422 429 383T381 315T301 241Q265 210 201 149L142 93L218 92Q375 92 385 97Q392 99 409 186V189H449V186Q448 183 436 95T421 3V0H50V19V31Q50 38 56 46T86 81Q115 113 136 137Q145 147 170 174T204 211T233 244T261 278T284 308T305 340T320 369T333 401T340 431T343 464Q343 527 309 573T212 619Q179 619 154 602T119 569T109 550Q109 549 114 549Q132 549 151 535T170 489Q170 464 154 447T109 429Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-22EF" d="M78 250Q78 274 95 292T138 310Q162 310 180 294T199 251Q199 226 182 208T139 190T96 207T78 250ZM525 250Q525 274 542 292T585 310Q609 310 627 294T646 251Q646 226 629 208T586 190T543 207T525 250ZM972 250Q972 274 989 292T1032 310Q1056 310 1074 294T1093 251Q1093 226 1076 208T1033 190T990 207T972 250Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-22EE" d="M78 30Q78 54 95 72T138 90Q162 90 180 74T199 31Q199 6 182 -12T139 -30T96 -13T78 30ZM78 440Q78 464 95 482T138 500Q162 500 180 484T199 441Q199 416 182 398T139 380T96 397T78 440ZM78 840Q78 864 95 882T138 900Q162 900 180 884T199 841Q199 816 182 798T139 780T96 797T78 840Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-22F1" d="M133 760Q133 784 150 802T193 820Q217 820 235 804T254 761Q254 736 237 718T194 700T151 717T133 760ZM580 460Q580 484 597 502T640 520Q664 520 682 504T701 461Q701 436 684 418T641 400T598 417T580 460ZM1027 160Q1027 184 1044 202T1087 220Q1111 220 1129 204T1148 161Q1148 136 1131 118T1088 100T1045 117T1027 160Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-29" d="M60 749L64 750Q69 750 74 750H86L114 726Q208 641 251 514T294 250Q294 182 284 119T261 12T224 -76T186 -143T145 -194T113 -227T90 -246Q87 -249 86 -250H74Q66 -250 63 -250T58 -247T55 -238Q56 -237 66 -225Q221 -64 221 250T66 725Q56 737 55 738Q55 746 60 749Z"></path>
		<path stroke-width="10" id="E1-MJSZ4-239B" d="M837 1154Q843 1148 843 1145Q843 1141 818 1106T753 1002T667 841T574 604T494 299Q417 -84 417 -609Q417 -641 416 -647T411 -654Q409 -655 366 -655Q299 -655 297 -654Q292 -652 292 -643T291 -583Q293 -400 304 -242T347 110T432 470T574 813T785 1136Q787 1139 790 1142T794 1147T796 1150T799 1152T802 1153T807 1154T813 1154H819H837Z"></path>
		<path stroke-width="10" id="E1-MJSZ4-239D" d="M843 -635Q843 -638 837 -644H820Q801 -644 800 -643Q792 -635 785 -626Q684 -503 605 -363T473 -75T385 216T330 518T302 809T291 1093Q291 1144 291 1153T296 1164Q298 1165 366 1165Q409 1165 411 1164Q415 1163 416 1157T417 1119Q417 529 517 109T833 -617Q843 -631 843 -635Z"></path>
		<path stroke-width="10" id="E1-MJSZ4-239C" d="M413 -9Q412 -9 407 -9T388 -10T354 -10Q300 -10 297 -9Q294 -8 293 -5Q291 5 291 127V300Q291 602 292 605L296 609Q298 610 366 610Q382 610 392 610T407 610T412 609Q416 609 416 592T417 473V127Q417 -9 413 -9Z"></path>
		<path stroke-width="10" id="E1-MJSZ4-239E" d="M31 1143Q31 1154 49 1154H59Q72 1154 75 1152T89 1136Q190 1013 269 873T401 585T489 294T544 -8T572 -299T583 -583Q583 -634 583 -643T577 -654Q575 -655 508 -655Q465 -655 463 -654Q459 -653 458 -647T457 -609Q457 -58 371 340T100 1037Q87 1059 61 1098T31 1143Z"></path>
		<path stroke-width="10" id="E1-MJSZ4-23A0" d="M56 -644H50Q31 -644 31 -635Q31 -632 37 -622Q69 -579 100 -527Q286 -228 371 170T457 1119Q457 1161 462 1164Q464 1165 520 1165Q575 1165 577 1164Q582 1162 582 1153T583 1093Q581 910 570 752T527 400T442 40T300 -303T89 -626Q78 -640 75 -642T61 -644H56Z"></path>
		<path stroke-width="10" id="E1-MJSZ4-239F" d="M579 -9Q578 -9 573 -9T554 -10T520 -10Q466 -10 463 -9Q460 -8 459 -5Q457 5 457 127V300Q457 602 458 605L462 609Q464 610 532 610Q548 610 558 610T573 610T578 609Q582 609 582 592T583 473V127Q583 -9 579 -9Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		 <use xlink:href="#E1-MJMATHI-41" x="0" y="0"></use>
		<g transform="translate(755,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6D" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="883" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="1166" y="0"></use>
		</g>
		 <use xlink:href="#E1-MJMAIN-3D" x="2385" y="0"></use>
		<g transform="translate(3445,0)">
		<g transform="translate(0,3357)">
		 <use xlink:href="#E1-MJSZ4-239B" x="0" y="-1165"></use>
		<g transform="translate(0,-4348.860866874306) scale(1,4.143323978829526)">
		 <use xlink:href="#E1-MJSZ4-239C"></use>
		</g>
		 <use xlink:href="#E1-MJSZ4-239D" x="0" y="-5561"></use>
		</g>
		<g transform="translate(1047,0)">
		<g transform="translate(-11,0)">
		<g transform="translate(133,2557)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="505" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="788" y="0"></use>
		</g>
		</g>
		<g transform="translate(133,1062)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="505" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="788" y="0"></use>
		</g>
		</g>
		 <use xlink:href="#E1-MJMAIN-22EE" x="766" y="-1163"></use>
		<g transform="translate(0,-2563)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6D" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="883" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="1166" y="0"></use>
		</g>
		</g>
		</g>
		<g transform="translate(2805,0)">
		<g transform="translate(133,2557)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="505" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="788" y="0"></use>
		</g>
		</g>
		<g transform="translate(133,1062)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="505" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="788" y="0"></use>
		</g>
		</g>
		 <use xlink:href="#E1-MJMAIN-22EE" x="766" y="-1163"></use>
		<g transform="translate(0,-2563)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6D" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="883" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="1166" y="0"></use>
		</g>
		</g>
		</g>
		<g transform="translate(5620,0)">
		 <use xlink:href="#E1-MJMAIN-22EF" x="55" y="2557"></use>
		 <use xlink:href="#E1-MJMAIN-22EF" x="55" y="1062"></use>
		 <use xlink:href="#E1-MJMAIN-22F1" x="0" y="-1163"></use>
		 <use xlink:href="#E1-MJMAIN-22EF" x="55" y="-2563"></use>
		</g>
		<g transform="translate(7907,0)">
		<g transform="translate(133,2557)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-31" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="505" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="788" y="0"></use>
		</g>
		</g>
		<g transform="translate(133,1062)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="505" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="788" y="0"></use>
		</g>
		</g>
		 <use xlink:href="#E1-MJMAIN-22EE" x="801" y="-1163"></use>
		<g transform="translate(0,-2563)">
		 <use xlink:href="#E1-MJMATHI-61" x="0" y="0"></use>
		<g transform="translate(534,-150)">
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6D" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-2C" x="883" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMATHI-6E" x="1166" y="0"></use>
		</g>
		</g>
		</g>
		</g>
		<g transform="translate(11007,3357)">
		 <use xlink:href="#E1-MJSZ4-239E" x="0" y="-1164"></use>
		<g transform="translate(0,-4348.8281055638545) scale(1,4.1448840412320225)">
		 <use xlink:href="#E1-MJSZ4-239F"></use>
		</g>
		 <use xlink:href="#E1-MJSZ4-23A0" x="0" y="-5561"></use>
		</g>
		</g>
		</g>
		</svg>`
		svgMap["a8fff4d1"] = `<svg xmlns:xlink="http://www.w3.org/1999/xlink" width="20.167ex" height="6ex" style="vertical-align: -2.5ex; margin-bottom: 1px; margin-top: 1px;" viewBox="0 -1531.3 8701.9 2562.7" xmlns="http://www.w3.org/2000/svg">
		<defs>
		<path stroke-width="10" id="E1-MJMATHI-50" d="M287 628Q287 635 230 637Q206 637 199 638T192 648Q192 649 194 659Q200 679 203 681T397 683Q587 682 600 680Q664 669 707 631T751 530Q751 453 685 389Q616 321 507 303Q500 302 402 301H307L277 182Q247 66 247 59Q247 55 248 54T255 50T272 48T305 46H336Q342 37 342 35Q342 19 335 5Q330 0 319 0Q316 0 282 1T182 2Q120 2 87 2T51 1Q33 1 33 11Q33 13 36 25Q40 41 44 43T67 46Q94 46 127 49Q141 52 146 61Q149 65 218 339T287 628ZM645 554Q645 567 643 575T634 597T609 619T560 635Q553 636 480 637Q463 637 445 637T416 636T404 636Q391 635 386 627Q384 621 367 550T332 412T314 344Q314 342 395 342H407H430Q542 342 590 392Q617 419 631 471T645 554Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-28" d="M94 250Q94 319 104 381T127 488T164 576T202 643T244 695T277 729T302 750H315H319Q333 750 333 741Q333 738 316 720T275 667T226 581T184 443T167 250T184 58T225 -81T274 -167T316 -220T333 -241Q333 -250 318 -250H315H302L274 -226Q180 -141 137 -14T94 250Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-41" d="M208 74Q208 50 254 46Q272 46 272 35Q272 34 270 22Q267 8 264 4T251 0Q249 0 239 0T205 1T141 2Q70 2 50 0H42Q35 7 35 11Q37 38 48 46H62Q132 49 164 96Q170 102 345 401T523 704Q530 716 547 716H555H572Q578 707 578 706L606 383Q634 60 636 57Q641 46 701 46Q726 46 726 36Q726 34 723 22Q720 7 718 4T704 0Q701 0 690 0T651 1T578 2Q484 2 455 0H443Q437 6 437 9T439 27Q443 40 445 43L449 46H469Q523 49 533 63L521 213H283L249 155Q208 86 208 74ZM516 260Q516 271 504 416T490 562L463 519Q447 492 400 412L310 260L413 259Q516 259 516 260Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-3D" d="M56 347Q56 360 70 367H707Q722 359 722 347Q722 336 708 328L390 327H72Q56 332 56 347ZM56 153Q56 168 72 173H708Q722 163 722 153Q722 140 707 133H70Q56 140 56 153Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-32" d="M109 429Q82 429 66 447T50 491Q50 562 103 614T235 666Q326 666 387 610T449 465Q449 422 429 383T381 315T301 241Q265 210 201 149L142 93L218 92Q375 92 385 97Q392 99 409 186V189H449V186Q448 183 436 95T421 3V0H50V19V31Q50 38 56 46T86 81Q115 113 136 137Q145 147 170 174T204 211T233 244T261 278T284 308T305 340T320 369T333 401T340 431T343 464Q343 527 309 573T212 619Q179 619 154 602T119 569T109 550Q109 549 114 549Q132 549 151 535T170 489Q170 464 154 447T109 429Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-7C" d="M139 -249H137Q125 -249 119 -235V251L120 737Q130 750 139 750Q152 750 159 735V-235Q151 -249 141 -249H139Z"></path>
		<path stroke-width="10" id="E1-MJMATHI-42" d="M231 637Q204 637 199 638T194 649Q194 676 205 682Q206 683 335 683Q594 683 608 681Q671 671 713 636T756 544Q756 480 698 429T565 360L555 357Q619 348 660 311T702 219Q702 146 630 78T453 1Q446 0 242 0Q42 0 39 2Q35 5 35 10Q35 17 37 24Q42 43 47 45Q51 46 62 46H68Q95 46 128 49Q142 52 147 61Q150 65 219 339T288 628Q288 635 231 637ZM649 544Q649 574 634 600T585 634Q578 636 493 637Q473 637 451 637T416 636H403Q388 635 384 626Q382 622 352 506Q352 503 351 500L320 374H401Q482 374 494 376Q554 386 601 434T649 544ZM595 229Q595 273 572 302T512 336Q506 337 429 337Q311 337 310 336Q310 334 293 263T258 122L240 52Q240 48 252 48T333 46Q422 46 429 47Q491 54 543 105T595 229Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-3E" d="M84 520Q84 528 88 533T96 539L99 540Q106 540 253 471T544 334L687 265Q694 260 694 250T687 235Q685 233 395 96L107 -40H101Q83 -38 83 -20Q83 -19 83 -17Q82 -10 98 -1Q117 9 248 71Q326 108 378 132L626 250L378 368Q90 504 86 509Q84 513 84 520Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-34" d="M462 0Q444 3 333 3Q217 3 199 0H190V46H221Q241 46 248 46T265 48T279 53T286 61Q287 63 287 115V165H28V211L179 442Q332 674 334 675Q336 677 355 677H373L379 671V211H471V165H379V114Q379 73 379 66T385 54Q393 47 442 46H471V0H462ZM293 211V545L74 212L183 211H293Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-29" d="M60 749L64 750Q69 750 74 750H86L114 726Q208 641 251 514T294 250Q294 182 284 119T261 12T224 -76T186 -143T145 -194T113 -227T90 -246Q87 -249 86 -250H74Q66 -250 63 -250T58 -247T55 -238Q56 -237 66 -225Q221 -64 221 250T66 725Q56 737 55 738Q55 746 60 749Z"></path>
		<path stroke-width="10" id="E1-MJSZ3-28" d="M701 -940Q701 -943 695 -949H664Q662 -947 636 -922T591 -879T537 -818T475 -737T412 -636T350 -511T295 -362T250 -186T221 17T209 251Q209 962 573 1361Q596 1386 616 1405T649 1437T664 1450H695Q701 1444 701 1441Q701 1436 681 1415T629 1356T557 1261T476 1118T400 927T340 675T308 359Q306 321 306 250Q306 -139 400 -430T690 -924Q701 -936 701 -940Z"></path>
		<path stroke-width="10" id="E1-MJMAIN-2223" d="M139 -249H137Q125 -249 119 -235V251L120 737Q130 750 139 750Q152 750 159 735V-235Q151 -249 141 -249H139Z"></path>
		<path stroke-width="10" id="E1-MJSZ3-29" d="M34 1438Q34 1446 37 1448T50 1450H56H71Q73 1448 99 1423T144 1380T198 1319T260 1238T323 1137T385 1013T440 864T485 688T514 485T526 251Q526 134 519 53Q472 -519 162 -860Q139 -885 119 -904T86 -936T71 -949H56Q43 -949 39 -947T34 -937Q88 -883 140 -813Q428 -430 428 251Q428 453 402 628T338 922T245 1146T145 1309T46 1425Q44 1427 42 1429T39 1433T36 1436L34 1438Z"></path>
		</defs>
		<g stroke="currentColor" fill="currentColor" stroke-width="0" transform="matrix(1 0 0 -1 0 0)">
		 <use xlink:href="#E1-MJMATHI-50" x="0" y="0"></use>
		<g transform="translate(922,0)">
		 <use xlink:href="#E1-MJSZ3-28"></use>
		 <use xlink:href="#E1-MJMATHI-41" x="741" y="0"></use>
		 <use xlink:href="#E1-MJMAIN-3D" x="1773" y="0"></use>
		 <use xlink:href="#E1-MJMAIN-32" x="2834" y="0"></use>
		<g transform="translate(3339,1517)">
		 <use xlink:href="#E1-MJMAIN-2223" x="0" y="-760"></use>
		<g transform="translate(0,-1414.117134253976) scale(1,0.585680439807863)">
		 <use xlink:href="#E1-MJMAIN-2223"></use>
		</g>
		 <use xlink:href="#E1-MJMAIN-2223" x="0" y="-2276"></use>
		</g>
		<g transform="translate(3622,0)">
		<g transform="translate(120,0)">
		<rect stroke="none" width="1332" height="60" x="0" y="220"></rect>
		<g transform="translate(60,676)">
		 <use xlink:href="#E1-MJMATHI-41" x="0" y="0"></use>
		 <use transform="scale(0.707)" xlink:href="#E1-MJMAIN-32" x="1067" y="513"></use>
		</g>
		 <use xlink:href="#E1-MJMATHI-42" x="284" y="-713"></use>
		</g>
		</g>
		 <use xlink:href="#E1-MJMAIN-3E" x="5472" y="0"></use>
		 <use xlink:href="#E1-MJMAIN-34" x="6533" y="0"></use>
		 <use xlink:href="#E1-MJSZ3-29" x="7038" y="-1"></use>
		</g>
		</g>
		</svg>`

		re := regexp.MustCompile("<ons-equation\\spath=\"([-A-Za-z0-9+&@#/%?=~_|!:,.;()*$]+)\"?\\s?/>")
		return re.ReplaceAllStringFunc(markdown, func(matchedTag string) string {
			fmt.Printf("resolveEquations/replace() string: %s\n", matchedTag)
			submatches := re.FindStringSubmatch(matchedTag)
			path := submatches[1]
			fmt.Printf("resolveEquations/replace() submatches: %+v\n", submatches)
			fmt.Printf("resolveEquations/replace() path: %s\n", path)
			return fmt.Sprintf("<figure class=\"ons-figure\">%s</figure>", svgMap[path])
		})
	}

	/* Image sidecar
	{
		"type": "image",
		"filename": "1fb0627f",
		"title": "Dino 1",
		"uri": "/economy/economicoutputandproductivity/output/articles/foo1119/imageedition/1fb0627f",
		"subtitle": "Some dino",
		"source": "google",
		"notes": "notes for dino\nmore **notes** for dino",
		"altText": "yep, it's a dino",
		"files": [
			{
				"type": "uploaded-image",
				"filename": "1fb0627f.jpeg",
				"fileType": "jpeg"
			}
		]
	}
	*/

	subscript := func(s string) string {
		// TODO Babbage replaces  "~(?=\\S)(\\S*)~" with "<sub>$1</sub>"
		return s
	}

	superscript := func(s string) string {
		// TODO Babbage replaces "\\^(?=\\S)(\\S*)\\^" with "<sup>$1</sup>"
		return s
	}

	resolveImages := func(markdown string) string {
		sidecarMap := make(map[string][]byte)
		sidecarMap["1fb0627f"] = []byte(`{"type":"image","filename":"1fb0627f","title":"Dino 1","uri":"/economy/economicoutputandproductivity/output/articles/foo1119/imageedition/1fb0627f","subtitle":"Some dino","source":"google","notes":"notes for dino\nmore **notes** for dino","altText":"yep, it's a dino","files":[{"type":"uploaded-image","filename":"1fb0627f.jpeg","fileType":"jpeg"}]}`)
		sidecarMap["5967a85a"] = []byte(`{"type":"image","filename":"5967a85a","title":"More dinos","uri":"/economy/economicoutputandproductivity/output/articles/foo1119/imageedition/5967a85a","subtitle":"A bunch of em","source":"google","notes":"","altText":"roar","files":[{"type":"uploaded-image","filename":"5967a85a.jpeg","fileType":"jpeg"}]}`)

		re := regexp.MustCompile("<ons-image\\spath=\"([-A-Za-z0-9+&@#/%?=~_|!:,.;()*$]+)\"?\\s?/>")
		return re.ReplaceAllStringFunc(markdown, func(matchedTag string) string {
			fmt.Printf("resolveImages/replace() string: %s\n", matchedTag)
			submatches := re.FindStringSubmatch(matchedTag)
			path := submatches[1]
			fmt.Printf("resolveImages/replace() submatches: %+v\n", submatches)
			fmt.Printf("resolveImages/replace() path: %s\n", path)
			var sidecar SidecarImage
			err := json.Unmarshal(sidecarMap[path], &sidecar)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("resolveImages/replace() sidecar: %#v\n", sidecar)

			var template string
			if len(sidecar.Files) > 0 {
				title := subscript(superscript(sidecar.Title))
				subtitle := subscript(superscript(sidecar.Subtitle))
				template = `<div class="ons-u-mb-l">` // TODO In print view, page-break-inside: avoid;
				template = template + fmt.Sprintf(`<h4 class="ons-u-mt-m ons-u-pt-s ons-u-mb-xs">%s</h4>`, title)
				if len(subtitle) != 0 {
					template = template + fmt.Sprintf(`<h5 class="ons-u-fs-s">%s</h5>`, subtitle)
				}
				for _, file := range sidecar.Files {
					template = template + fmt.Sprintf(`<img src="/resource?uri=%s.%s" alt="%s">`, sidecar.Uri, file.FileType, sidecar.AltText)
				}
				if len(sidecar.Source) != 0 {
					template = template + fmt.Sprintf(`<h6 class="ons-u-fs-s--b ons-u-mt-s ons-u-mb-xs">Source: %s</h6>`, sidecar.Source)
				}
				if len(sidecar.Notes) != 0 {
					template = template + fmt.Sprintf(`<h6 class="ons-u-fs-s--b ons-u-mt-s ons-u-mb-xs">Notes:</h6>%s`, helper.Markdown(sidecar.Notes))
				}

				// TODO: Download option to be hidden in print view
				template = template + fmt.Sprintf(`
					<h6 class="ons-u-fs-s--b ons-u-mt-s ons-u-mb-xs">
						<span role="text">Download this image
							<span class="ons-u-vh">%s</span>
						</span>
					</h6>
				`, sidecar.Title)
				for _, file := range sidecar.Files {
					fileSize := "100kB" // TODO Calculate

					var ariaLabel string
					if file.Type == "uploaded-image" {
						ariaLabel = fmt.Sprintf("Download %s as %s (%s)", sidecar.Title, file.FileType, fileSize)
					} else if file.Type == "uploaded-data" {
						ariaLabel = fmt.Sprintf("Download %s (%s)", sidecar.Title, fileSize)
					}

					template = template + fmt.Sprintf(`
							<a
								class="ons-btn"
								title="Download as %s"
								data-gtm-title="%s"
								data-gtm-type="download-%s"
								href="/file?uri=%s.%s"
								aria-label="%s"
							>
								<span class="ons-btn__inner">
									<span class="ons-btn__text">.%s (%s)</span>
								</span>
							</a>
						`, file.FileType, title, file.FileType, sidecar.Uri, file.FileType, ariaLabel, file.FileType, fileSize)
				}
				template = template + "</div>"
			} else {
				template = fmt.Sprintf(`<img src="/resource?uri=%s">`, sidecar.Uri)
			}
			return template
		})
	}

	for index, section := range bulletin.Sections {
		bulletin.Sections[index].Markdown = resolveEquations(section.Markdown)
		fmt.Printf("section %d.Markdown: %s\n", index, bulletin.Sections[index].Markdown)
		bulletin.Sections[index].Markdown = resolveImages(section.Markdown)
		fmt.Printf("section %d.Markdown: %s\n", index, bulletin.Sections[index].Markdown)
	}

	breadcrumbs, err := zc.GetBreadcrumb(ctx, userAccessToken, collectionID, lang, bulletin.URI)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	requestProtocol := "http"
	if req.TLS != nil {
		requestProtocol = "https"
	}
	model := mapper.CreateBulletinModel(basePage, *bulletin, breadcrumbs, lang, requestProtocol)
	rc.BuildPage(w, model, "bulletin")
}
