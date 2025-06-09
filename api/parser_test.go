package goodreads_parser

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseBook(t *testing.T) {
	html := `<table><tbody><tr itemscope="" itemtype="http://schema.org/Book">
		<td width="5%" valign="top">
			<div id="42844155" class="u-anchorTarget"></div>
				<a title="Harry Potter and the Sorcerer’s Stone" href="/book/show/42844155-harry-potter-and-the-sorcerer-s-stone?from_search=true&amp;from_srp=true&amp;qid=MgxJ8chY0D&amp;rank=1">
					<img alt="Harry Potter and the Sorcer..." class="bookCover" itemprop="image" src="https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1598823299i/42844155._SX50_.jpg">
</a>    </td>
		<td width="100%" valign="top">
			<a class="bookTitle" itemprop="url" href="/book/show/42844155-harry-potter-and-the-sorcerer-s-stone?from_search=true&amp;from_srp=true&amp;qid=MgxJ8chY0D&amp;rank=1">
				<span itemprop="name" role="heading" aria-level="4">Harry Potter and the Sorcerer’s Stone (Harry Potter, #1)</span>
</a>      <br>
				<span class="by">by</span>
<span itemprop="author" itemscope="" itemtype="http://schema.org/Person">
<div class="authorName__container">
<a class="authorName" itemprop="url" href="https://www.goodreads.com/author/show/1077326.J_K_Rowling?from_search=true&amp;from_srp=true"><span itemprop="name">J.K. Rowling</span></a>,
</div>
<div class="authorName__container">
<a class="authorName" itemprop="url" href="https://www.goodreads.com/author/show/6458971.Olly_Moss?from_search=true&amp;from_srp=true"><span itemprop="name">Olly Moss</span></a> <span class="authorName greyText smallText role">(Illustrator)</span>
</div>
</span>

				<br>
				<div>
					<span class="greyText smallText uitext">
								<span class="minirating"><span class="stars staticStars notranslate"><span size="12x12" class="staticStar p10"></span><span size="12x12" class="staticStar p10"></span><span size="12x12" class="staticStar p10"></span><span size="12x12" class="staticStar p10"></span><span size="12x12" class="staticStar p3"></span></span> 4.47 avg rating — 9,852,011 ratings</span>
							—
								published
							 1997
							—
							<a class="greyText" rel="nofollow" href="/work/editions/4640799-harry-potter-and-the-philosopher-s-stone">1328 editions</a>
					</span>
				</div>




					<div style="float: left">
						<div class="wtrButtonContainer wtrSignedOut" id="1_book_42844155">
<div class="wtrUp wtrLeft">
<form action="/shelf/add_to_shelf" accept-charset="UTF-8" method="post"><input name="utf8" type="hidden" value="✓"><input type="hidden" name="authenticity_token" value="pjQ7EZxYxIBJRm+/ghIQkLB5pUdtuWIKepKasbLCV3j7+opu54PTrxOAyZlJkSf5tS6dbKIB3YBHNhYFbQWTpQ==">
<input type="hidden" name="book_id" id="book_id" value="42844155">
<input type="hidden" name="name" id="name" value="to-read">
<input type="hidden" name="unique_id" id="unique_id" value="1_book_42844155">
<input type="hidden" name="wtr_new" id="wtr_new" value="true">
<input type="hidden" name="from_choice" id="from_choice" value="false">
<input type="hidden" name="from_home_module" id="from_home_module" value="false">
<input type="hidden" name="ref" id="ref" value="" class="wtrLeftUpRef">
<input type="hidden" name="existing_review" id="existing_review" value="false" class="wtrExisting">
<input type="hidden" name="page_url" id="page_url">
<input type="hidden" name="from_search" id="from_search" value="true">
<input type="hidden" name="qid" id="qid" value="MgxJ8chY0D">
<input type="hidden" name="rank" id="rank" value="1">
<button class="wtrToRead" type="submit">
<span class="progressTrigger">Want to Read</span>
<span class="progressIndicator">saving…</span>
</button>
</form>

</div>

<div class="wtrRight wtrUp">
<form class="hiddenShelfForm" action="/shelf/add_to_shelf" accept-charset="UTF-8" method="post"><input name="utf8" type="hidden" value="✓"><input type="hidden" name="authenticity_token" value="pjQ7EZxYxIBJRm+/ghIQkLB5pUdtuWIKepKasbLCV3j7+opu54PTrxOAyZlJkSf5tS6dbKIB3YBHNhYFbQWTpQ==">
<input type="hidden" name="unique_id" id="unique_id" value="1_book_42844155">
<input type="hidden" name="book_id" id="book_id" value="42844155">
<input type="hidden" name="a" id="a">
<input type="hidden" name="name" id="name">
<input type="hidden" name="from_choice" id="from_choice" value="false">
<input type="hidden" name="from_home_module" id="from_home_module" value="false">
<input type="hidden" name="page_url" id="page_url">
<input type="hidden" name="from_search" id="from_search" value="true">
<input type="hidden" name="qid" id="qid" value="MgxJ8chY0D">
<input type="hidden" name="rank" id="rank" value="1">
</form>

<button class="wtrShelfButton"></button>
<div class="wtrShelfMenu">
<ul class="wtrExclusiveShelves">
<li><button class="wtrExclusiveShelf" name="name" type="submit" value="to-read">
<span class="progressTrigger">Want to Read</span>
<img alt="saving…" class="progressIndicator" src="https://s.gr-assets.com/assets/loading-trans-ced157046184c3bc7c180ffbfc6825a4.gif">
</button>
</li>
<li><button class="wtrExclusiveShelf" name="name" type="submit" value="currently-reading">
<span class="progressTrigger">Currently Reading</span>
<img alt="saving…" class="progressIndicator" src="https://s.gr-assets.com/assets/loading-trans-ced157046184c3bc7c180ffbfc6825a4.gif">
</button>
</li>
<li><button class="wtrExclusiveShelf" name="name" type="submit" value="read">
<span class="progressTrigger">Read</span>
<img alt="saving…" class="progressIndicator" src="https://s.gr-assets.com/assets/loading-trans-ced157046184c3bc7c180ffbfc6825a4.gif">
</button>
</li>
</ul>
</div>
</div>

<div class="ratingStars wtrRating">
<div class="starsErrorTooltip hidden">
Error rating book. Refresh and try again.
</div>
<div class="myRating uitext greyText">Rate this book</div>
<div class="clearRating uitext" style="display: none;">Clear rating</div>
<div class="stars" data-resource-id="42844155" data-user-id="0" data-submit-url="/review/rate/42844155?from_search=true&amp;from_srp=true&amp;qid=MgxJ8chY0D&amp;rank=1&amp;stars_click=true&amp;wtr_button_id=1_book_42844155" data-rating="0" data-restore-rating="null"><a class="star off" title="did not like it" href="#" ref="">1 of 5 stars</a><a class="star off" title="it was ok" href="#" ref="">2 of 5 stars</a><a class="star off" title="liked it" href="#" ref="">3 of 5 stars</a><a class="star off" title="really liked it" href="#" ref="">4 of 5 stars</a><a class="star off" title="it was amazing" href="#" ref="">5 of 5 stars</a></div>
</div>

</div>

					</div>
					<div class="getACopyButtonWrapper getACopyButtonWrapper--desktop">
						<div data-react-class="ReactComponents.GetACopyButton" data-react-props="{&quot;getACopyDataUrl&quot;:&quot;/book/42844155/buy_buttons&quot;}"><div data-reactid=".jpxu06lili" data-react-checksum="250495492"><button class="gr-button gr-button--fullWidth u-paddingTopTiny u-paddingBottomTiny u-defaultType" data-reactid=".jpxu06lili.0">Get a copy</button></div></div>
					</div>


						</td>
						<td width="130px">

			</td>

	</tr></table></tbody>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}

	book, err := parseSearchBook(doc.Find("tr").First())
	if err != nil {
		t.Fatal(err)
	}

	if book.Title != "Harry Potter and the Sorcerer’s Stone (Harry Potter, #1)" {
		t.Errorf("expected title: %s, got: %s", "Harry Potter and the Sorcerer’s Stone (Harry Potter, #1)", book.Title)
	}

	if book.Author != "J.K. Rowling" {
		t.Errorf("expected author: %s, got: %s", "J.K. Rowling", book.Author)
	}

	if book.Url != "/book/show/42844155-harry-potter-and-the-sorcerer-s-stone?from_search=true&from_srp=true&qid=MgxJ8chY0D&rank=1" {
		t.Errorf("expected url: %s, got: %s", "https://www.goodreads.com/book/show/42844155-harry-potter-and-the-sorcerer-s-stone?from_search=true&from_srp=true&qid=MgxJ8chY0D&rank=1", book.Url)
	}

	if book.Id != 42844155 {
		t.Errorf("expected id: %d, got: %d", 42844155, book.Id)
	}
}
