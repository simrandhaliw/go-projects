// Unit tests for inline parsing

package mmark

import (
	"regexp"
	"testing"

	"strings"
)

func runMarkdownInline(input string, extensions, htmlFlags int, params HtmlRendererParameters) string {
	extensions |= EXTENSION_AUTOLINK

	renderer := HtmlRendererWithParameters(htmlFlags, "", "", params)

	return Parse([]byte(input), renderer, extensions).String()
}

func doTestsInline(t *testing.T, tests []string) {
	doTestsInlineParam(t, tests, 0, 0, HtmlRendererParameters{})
}

func doLinkTestsInline(t *testing.T, tests []string) {
	doTestsInline(t, tests)

	prefix := "http://localhost"
	params := HtmlRendererParameters{AbsolutePrefix: prefix}
	transformTests := transformLinks(tests, prefix)
	doTestsInlineParam(t, transformTests, 0, 0, params)
	doTestsInlineParam(t, transformTests, 0, commonHtmlFlags, params)
}

func doSafeTestsInline(t *testing.T, tests []string) {
	doTestsInlineParam(t, tests, 0, HTML_SAFELINK, HtmlRendererParameters{})

	// All the links in this test should not have the prefix appended, so
	// just rerun it with different parameters and the same expectations.
	prefix := "http://localhost"
	params := HtmlRendererParameters{AbsolutePrefix: prefix}
	transformTests := transformLinks(tests, prefix)
	doTestsInlineParam(t, transformTests, 0, HTML_SAFELINK, params)
}

func doTestsInlineParam(t *testing.T, tests []string, extensions, htmlFlags int,
	params HtmlRendererParameters) {
	// catch and report panics
	var candidate string
	/*
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("\npanic while processing [%#v] (%v)\n", candidate, err)
			}
		}()
	*/

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runMarkdownInline(candidate, extensions, htmlFlags, params)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownInline(candidate, extensions, htmlFlags, params)
				}
			}
		}
	}
}

func transformLinks(tests []string, prefix string) []string {
	newTests := make([]string, len(tests))
	anchorRe := regexp.MustCompile(`<a href="/(.*?)"`)
	imgRe := regexp.MustCompile(`<img src="/(.*?)"`)
	for i, test := range tests {
		if i%2 == 1 {
			test = anchorRe.ReplaceAllString(test, `<a href="`+prefix+`/$1"`)
			test = imgRe.ReplaceAllString(test, `<img src="`+prefix+`/$1"`)
		}
		newTests[i] = test
	}
	return newTests
}

func TestEmphasis(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple *inline* test\n",
		"<p>simple <em>inline</em> test</p>\n",

		"*at the* beginning\n",
		"<p><em>at the</em> beginning</p>\n",

		"at the *end*\n",
		"<p>at the <em>end</em></p>\n",

		"*try two* in *one line*\n",
		"<p><em>try two</em> in <em>one line</em></p>\n",

		"over *two\nlines* test\n",
		"<p>over <em>two\nlines</em> test</p>\n",

		"odd *number of* markers* here\n",
		"<p>odd <em>number of</em> markers* here</p>\n",

		"odd *number\nof* markers* here\n",
		"<p>odd <em>number\nof</em> markers* here</p>\n",

		"simple _inline_ test\n",
		"<p>simple <em>inline</em> test</p>\n",

		"_at the_ beginning\n",
		"<p><em>at the</em> beginning</p>\n",

		"at the _end_\n",
		"<p>at the <em>end</em></p>\n",

		"_try two_ in _one line_\n",
		"<p><em>try two</em> in <em>one line</em></p>\n",

		"over _two\nlines_ test\n",
		"<p>over <em>two\nlines</em> test</p>\n",

		"odd _number of_ markers_ here\n",
		"<p>odd <em>number of</em> markers_ here</p>\n",

		"odd _number\nof_ markers_ here\n",
		"<p>odd <em>number\nof</em> markers_ here</p>\n",

		"mix of *markers_\n",
		"<p>mix of *markers_</p>\n",

		"This is un*believe*able! This d_oe_s not work!\n",
		"<p>This is un<em>believe</em>able! This d_oe_s not work!</p>\n",

		"This is un**believe**able! This d_oe_s not work!\n",
		"<p>This is un<strong>believe</strong>able! This d_oe_s not work!</p>\n",

		"_does_ work\n",
		"<p><em>does</em> work</p>\n",

		"This is a * literal asterisk as is this *\n",
		"<p>This is a * literal asterisk as is this *</p>\n",

		"These are ** two literal asterisk.\n",
		"<p>These are ** two literal asterisk.</p>\n",

		"As \\*are\\* these!\n",
		"<p>As *are* these!</p>\n",

		"This is an underscore _ and a _word_\n",
		"<p>This is an underscore _ and a <em>word</em></p>\n",

		"This is a not a literal one *hallo there!*\n",
		"<p>This is a not a literal one <em>hallo there!</em></p>\n",

		"**This**",
		"<p><strong>This</strong></p>\n",

		"*What is A\\* algorithm?*\n",
		"<p><em>What is A* algorithm?</em></p>\n",
	}
	doTestsInline(t, tests)
}

func TestStrong(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple **inline** test\n",
		"<p>simple <strong>inline</strong> test</p>\n",

		"**at the** beginning\n",
		"<p><strong>at the</strong> beginning</p>\n",

		"at the **end**\n",
		"<p>at the <strong>end</strong></p>\n",

		"**try two** in **one line**\n",
		"<p><strong>try two</strong> in <strong>one line</strong></p>\n",

		"over **two\nlines** test\n",
		"<p>over <strong>two\nlines</strong> test</p>\n",

		"odd **number of** markers** here\n",
		"<p>odd <strong>number of</strong> markers** here</p>\n",

		"odd **number\nof** markers** here\n",
		"<p>odd <strong>number\nof</strong> markers** here</p>\n",

		"simple __inline__ test\n",
		"<p>simple <strong>inline</strong> test</p>\n",

		"__at the__ beginning\n",
		"<p><strong>at the</strong> beginning</p>\n",

		"at the __end__\n",
		"<p>at the <strong>end</strong></p>\n",

		"__try two__ in __one line__\n",
		"<p><strong>try two</strong> in <strong>one line</strong></p>\n",

		"over __two\nlines__ test\n",
		"<p>over <strong>two\nlines</strong> test</p>\n",

		"odd __number of__ markers__ here\n",
		"<p>odd <strong>number of</strong> markers__ here</p>\n",

		"odd __number\nof__ markers__ here\n",
		"<p>odd <strong>number\nof</strong> markers__ here</p>\n",

		"mix of **markers__\n",
		"<p>mix of **markers__</p>\n",

		"**`/usr`** : this folder is named `usr`\n",
		"<p><strong><code>/usr</code></strong> : this folder is named <code>usr</code></p>\n",

		"**`/usr`** :\n\n this folder is named `usr`\n",
		"<p><strong><code>/usr</code></strong> :</p>\n\n<p>this folder is named <code>usr</code></p>\n",
	}
	doTestsInline(t, tests)
}

func TestEmphasisMix(t *testing.T) {
	var tests = []string{
		"***triple emphasis***\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"***triple\nemphasis***\n",
		"<p><strong><em>triple\nemphasis</em></strong></p>\n",

		"___triple emphasis___\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"***triple emphasis___\n",
		"<p>***triple emphasis___</p>\n",

		"*__triple emphasis__*\n",
		"<p><em><strong>triple emphasis</strong></em></p>\n",

		"__*triple emphasis*__\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"**improper *nesting** is* bad\n",
		"<p><strong>improper *nesting</strong> is* bad</p>\n",

		"*improper **nesting* is** bad\n",
		"<p>*improper <strong>nesting* is</strong> bad</p>\n",
	}
	doTestsInline(t, tests)
}

func TestEmphasisLink(t *testing.T) {
	var tests = []string{
		"[first](before) *text[second] (inside)text* [third](after)\n",
		"<p><a href=\"before\">first</a> <em>text<a href=\"inside\">second</a>text</em> <a href=\"after\">third</a></p>\n",

		"*incomplete [link] definition*\n",
		"<p><em>incomplete [link] definition</em></p>\n",

		"*it's [emphasis*] (not link)\n",
		"<p><em>it's [emphasis</em>] (not link)</p>\n",

		"*it's [emphasis*] and *[asterisk]\n",
		"<p><em>it's [emphasis</em>] and *[asterisk]</p>\n",
	}
	doTestsInline(t, tests)
}

func TestStrikeThrough(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple ~~inline~~ test\n",
		"<p>simple <del>inline</del> test</p>\n",

		"~~at the~~ beginning\n",
		"<p><del>at the</del> beginning</p>\n",

		"at the ~~end~~\n",
		"<p>at the <del>end</del></p>\n",

		"~~try two~~ in ~~one line~~\n",
		"<p><del>try two</del> in <del>one line</del></p>\n",

		"over ~~two\nlines~~ test\n",
		"<p>over <del>two\nlines</del> test</p>\n",

		"odd ~~number of~~ markers~~ here\n",
		"<p>odd <del>number of</del> markers~~ here</p>\n",

		"odd ~~number\nof~~ markers~~ here\n",
		"<p>odd <del>number\nof</del> markers~~ here</p>\n",
	}
	doTestsInline(t, tests)
}

func TestSubscript(t *testing.T) {
	var tests = []string{
		"H~2~O is a liquid. is 1024. but this is ~~strikethrough~~ text\n",
		"<p>H<sub>2</sub>O is a liquid. is 1024. but this is <del>strikethrough</del> text</p>\n",

		"~[hallo](http:/ddssd)~",
		"<p><sub><a href=\"http:/ddssd\">hallo</a></sub></p>\n",

		"~subtext~ and some other ~subtext~",
		"<p><sub>subtext</sub> and some other <sub>subtext</sub></p>\n",

		"no~thing inline\n",
		"<p>no~thing inline</p>\n",

		"simple ~~in~l~ine~~ test\n",
		"<p>simple <del>in<sub>l</sub>ine</del> test</p>\n",

		"~odd ~~number o~~ markers~~ here\n",
		"<p>~odd <del>number o</del> markers~~ here</p>\n",

		"~subtext~ and some other ~subtext~",
		"<p><sub>subtext</sub> and some other <sub>subtext</sub></p>\n",

		" ~subtext~ and some other ~subtext~ ",
		"<p><sub>subtext</sub> and some other <sub>subtext</sub></p>\n",

		"odd ~~number\nof~~ ma~~rkers~ her~e~\n",
		"<p>odd <del>number\nof</del> ma~<sub>rkers</sub> her<sub>e</sub></p>\n",

		"~boo ~bla",
		"<p>~boo ~bla</p>\n",

		"~boo\\ ~bla",
		"<p><sub>boo </sub>bla</p>\n",

		"\\^not sub\\^",
		"<p>^not sub^</p>\n",
	}
	doTestsInline(t, tests)
}

func TestSuperscript(t *testing.T) {
	var tests = []string{
		"H~2~O is a liquid. 2^10^ is 1024. but this is ~~strikethrough~~ text\n",
		"<p>H<sub>2</sub>O is a liquid. 2<sup>10</sup> is 1024. but this is <del>strikethrough</del> text</p>\n",

		"^[hallo](http:/ddssd)^",
		"<p><sup><a href=\"http:/ddssd\">hallo</a></sup></p>\n",

		"^[hallo]^ is superscript, not a footnote.",
		"<p><sup>[hallo]</sup> is superscript, not a footnote.</p>\n",

		"^[^]",
		"<p><sup>[</sup>]</p>\n",

		"^[a^]",
		"<p><sup>[a</sup>]</p>\n",

		"types ^inlinenote  and [^2] footnote",
		"<p>types ^inlinenote  and [^2] footnote</p>\n",

		"types ^[inlinenote] of ^[other inline with a [hallo](http://miek.nl) notes]",
		"<p>types ^[inlinenote] of ^[other inline with a <a href=\"http://miek.nl\">hallo</a> notes]</p>\n",

		"^[hallo](http://miek.nl)^",
		"<p><sup><a href=\"http://miek.nl\">hallo</a></sup></p>\n",

		"^boo ^bla",
		"<p>^boo ^bla</p>\n",

		"^boo\\ ^bla",
		"<p><sup>boo </sup>bla</p>\n",

		"P~a\\ cat~",
		"<p>P<sub>a cat</sub></p>\n",

		"\\^not sup\\^",
		"<p>^not sup^</p>\n",
	}
	doTestsInline(t, tests)
}

func TestCodeSpan(t *testing.T) {
	var tests = []string{
		"`source code`\n",
		"<p><code>source code</code></p>\n",

		"` source code with spaces `\n",
		"<p><code>source code with spaces</code></p>\n",

		"` source code with spaces `not here\n",
		"<p><code>source code with spaces</code>not here</p>\n",

		"a `single marker\n",
		"<p>a `single marker</p>\n",

		"a single multi-tick marker with ``` no text\n",
		"<p>a single multi-tick marker with ``` no text</p>\n",

		"markers with ` ` a space\n",
		"<p>markers with  a space</p>\n",

		"`source code` and a `stray\n",
		"<p><code>source code</code> and a `stray</p>\n",

		"`source *with* _awkward characters_ in it`\n",
		"<p><code>source *with* _awkward characters_ in it</code></p>\n",

		"`split over\ntwo lines`\n",
		"<p><code>split over\ntwo lines</code></p>\n",

		"```multiple ticks``` for the marker\n",
		"<p><code>multiple ticks</code> for the marker</p>\n",

		"```multiple ticks `with` ticks inside```\n",
		"<p><code>multiple ticks `with` ticks inside</code></p>\n",
	}
	doTestsInline(t, tests)
}

func TestLineBreak(t *testing.T) {
	var tests = []string{
		"this line  \nhas a break\n",
		"<p>this line<br>\nhas a break</p>\n",

		"this line \ndoes not\n",
		"<p>this line\ndoes not</p>\n",

		"this line\\\ndoes also!\n",
		"<p>this line<br>\ndoes also!</p>\n",

		"this line\\ \ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this has an   \nextra space\n",
		"<p>this has an<br>\nextra space</p>\n",
	}
	doTestsInline(t, tests)

	tests = []string{
		"this line  \nhas a break\n",
		"<p>this line<br>\nhas a break</p>\n",

		"this line \ndoes not\n",
		"<p>this line\ndoes not</p>\n",

		"this line\\\nhas a break\n",
		"<p>this line<br>\nhas a break</p>\n",

		"this line\\ \ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this has an   \nextra space\n",
		"<p>this has an<br>\nextra space</p>\n",
	}
	doTestsInlineParam(t, tests, EXTENSION_BACKSLASH_LINE_BREAK, 0, HtmlRendererParameters{})
}

func TestInlineLink(t *testing.T) {
	var tests = []string{
		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo with a title](/bar/ \"title\")\n",
		"<p><a href=\"/bar/\" title=\"title\">foo with a title</a></p>\n",

		"[foo with a title](/bar/\t\"title\")\n",
		"<p><a href=\"/bar/\" title=\"title\">foo with a title</a></p>\n",

		"[foo with a title](/bar/ \"title\"  )\n",
		"<p><a href=\"/bar/\" title=\"title\">foo with a title</a></p>\n",

		"[foo with a title](/bar/ title with no quotes)\n",
		"<p><a href=\"/bar/ title with no quotes\">foo with a title</a></p>\n",

		"[foo]()\n",
		"<p>[foo]()</p>\n",

		"![foo](/bar/)\n",
		"<p><figure><img src=\"/bar/\" alt=\"foo\"></figure></p>\n",

		"![foo with a title](/bar/ \"title\")\n",
		"<p><figure><img src=\"/bar/\" alt=\"foo with a title\" title=\"title\"><figcaption>title</figcaption></figure></p>\n",

		"![foo with a title](/bar/\t\"title\")\n",
		"<p><figure><img src=\"/bar/\" alt=\"foo with a title\" title=\"title\"><figcaption>title</figcaption></figure></p>\n",

		"![foo with a title](/bar/ \"title\"  )\n",
		"<p><figure><img src=\"/bar/\" alt=\"foo with a title\" title=\"title\"><figcaption>title</figcaption></figure></p>\n",

		"![foo with a title](/bar/ title with no quotes)\n",
		"<p><figure><img src=\"/bar/ title with no quotes\" alt=\"foo with a title\"></figure></p>\n",

		"![](img.jpg)\n",
		"<p><figure><img src=\"img.jpg\" alt=\"\"></figure></p>\n",

		"[link](url)\n",
		"<p><a href=\"url\">link</a></p>\n",

		"![foo]()\n",
		"<p>![foo]()</p>\n",

		"[a link]\t(/with_a_tab/)\n",
		"<p><a href=\"/with_a_tab/\">a link</a></p>\n",

		"[a link]  (/with_spaces/)\n",
		"<p><a href=\"/with_spaces/\">a link</a></p>\n",

		"[text (with) [[nested] (brackets)]](/url/)\n",
		"<p><a href=\"/url/\">text (with) [[nested] (brackets)]</a></p>\n",

		"[text (with) [broken nested] (brackets)]](/url/)\n",
		"<p>[text (with) <a href=\"brackets\">broken nested</a>]](/url/)</p>\n",

		"[text\nwith a newline](/link/)\n",
		"<p><a href=\"/link/\">text\nwith a newline</a></p>\n",

		"[text in brackets] [followed](/by a link/)\n",
		"<p>[text in brackets] <a href=\"/by a link/\">followed</a></p>\n",

		"[link with\\] a closing bracket](/url/)\n",
		"<p><a href=\"/url/\">link with] a closing bracket</a></p>\n",

		"[link with\\[ an opening bracket](/url/)\n",
		"<p><a href=\"/url/\">link with[ an opening bracket</a></p>\n",

		"[link with\\) a closing paren](/url/)\n",
		"<p><a href=\"/url/\">link with) a closing paren</a></p>\n",

		"[link with\\( an opening paren](/url/)\n",
		"<p><a href=\"/url/\">link with( an opening paren</a></p>\n",

		"[link](  with whitespace)\n",
		"<p><a href=\"with whitespace\">link</a></p>\n",

		"[link](  with whitespace   )\n",
		"<p><a href=\"with whitespace\">link</a></p>\n",

		"[![image](someimage)](with image)\n",
		"<p><a href=\"with image\"><figure><img src=\"someimage\" alt=\"image\"></figure></a></p>\n",

		"[link](url \"one quote)\n",
		"<p><a href=\"url &quot;one quote\">link</a></p>\n",

		"[link](url 'one quote)\n",
		"<p><a href=\"url 'one quote\">link</a></p>\n",

		"[link](<url>)\n",
		"<p><a href=\"url\">link</a></p>\n",

		"[link & ampersand](/url/)\n",
		"<p><a href=\"/url/\">link &amp; ampersand</a></p>\n",

		"[link &amp; ampersand](/url/)\n",
		"<p><a href=\"/url/\">link &amp; ampersand</a></p>\n",

		"[link](/url/&query)\n",
		"<p><a href=\"/url/&amp;query\">link</a></p>\n",

		"[[t]](/t)\n",
		"<p><a href=\"/t\">[t]</a></p>\n",

		// Issue #126 on blackfriday
		"[Chef](https://en.wikipedia.org/wiki/Chef_(film))\n",
		"<p><a href=\"https://en.wikipedia.org/wiki/Chef_(film)\">Chef</a></p>\n",

		"[Chef](https://en.wikipedia.org/wiki/Chef_(film)_moretext)\n",
		"<p><a href=\"https://en.wikipedia.org/wiki/Chef_(film)_moretext\">Chef</a></p>\n",

		// Issue 116 in blackfriday
		"![](http://www.broadgate.co.uk/Content/Upload/DetailImages/Cyclus700(1).jpg)",
		"<p><figure><img src=\"http://www.broadgate.co.uk/Content/Upload/DetailImages/Cyclus700(1).jpg\" alt=\"\"></figure></p>\n",

		// no closing ), autolinking detects the url next
		"[disambiguation](http://en.wikipedia.org/wiki/Disambiguation_(disambiguation) is the",
		"<p>[disambiguation](<a href=\"http://en.wikipedia.org/wiki/Disambiguation_(disambiguation\">http://en.wikipedia.org/wiki/Disambiguation_(disambiguation</a>) is the</p>\n",

		"[disambiguation](http://en.wikipedia.org/wiki/Disambiguation_(disambiguation)) is the",
		"<p><a href=\"http://en.wikipedia.org/wiki/Disambiguation_(disambiguation)\">disambiguation</a> is the</p>\n",

		"[disambiguation](http://en.wikipedia.org/wiki/Disambiguation_(disambiguation))",
		"<p><a href=\"http://en.wikipedia.org/wiki/Disambiguation_(disambiguation)\">disambiguation</a></p>\n",
	}
	doLinkTestsInline(t, tests)

}

func TestNofollowLink(t *testing.T) {
	var tests = []string{
		"[foo](http://bar.com/foo/)\n",
		"<p><a href=\"http://bar.com/foo/\" rel=\"nofollow\">foo</a></p>\n",

		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",
	}
	doTestsInlineParam(t, tests, 0, HTML_SAFELINK|HTML_NOFOLLOW_LINKS,
		HtmlRendererParameters{})
}

func TestMath(t *testing.T) {
	var tests = []string{
		"{#eq1}\n $$ E = MC^2 $$",
		"<p><span  id=\"eq1\" class=\"math\">\\[ E = MC^2 \\]</span></p>\n",

		"Another paragraph, with some inline math $$x^2$$",
		"<p>Another paragraph, with some inline math <span  class=\"math\">\\(x^2\\)</span></p>\n",

		"$$ E = MC^2 $$",
		"<p><span  class=\"math\">\\[ E = MC^2 \\]</span></p>\n",

		"$$ E = MC^2 $$ not so much $$ again yes $$",
		"<p><span  class=\"math\">\\( E = MC^2 \\)</span> not so much <span  class=\"math\">\\( again yes \\)</span></p>\n",

		"you can use $$\\Phi = \\Phi + 1$$ in your source code.",
		"<p>you can use <span  class=\"math\">\\(\\Phi = \\Phi + 1\\)</span> in your source code.</p>\n",
	}
	doTestsInlineParam(t, tests, EXTENSION_MATH, 0, HtmlRendererParameters{})
}

func TestHrefTargetBlank(t *testing.T) {
	var tests = []string{
		// internal link
		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo](http://example.com)\n",
		"<p><a href=\"http://example.com\" target=\"_blank\">foo</a></p>\n",
	}
	doTestsInlineParam(t, tests, 0, HTML_SAFELINK|HTML_HREF_TARGET_BLANK, HtmlRendererParameters{})
}

func TestSafeInlineLink(t *testing.T) {
	var tests = []string{
		"[foo](/bar/)\n",
		"<p><a href=\"/bar/\">foo</a></p>\n",

		"[foo](http://bar/)\n",
		"<p><a href=\"http://bar/\">foo</a></p>\n",

		"[foo](https://bar/)\n",
		"<p><a href=\"https://bar/\">foo</a></p>\n",

		"[foo](ftp://bar/)\n",
		"<p><a href=\"ftp://bar/\">foo</a></p>\n",

		"[foo](mailto://bar/)\n",
		"<p><a href=\"mailto://bar/\">foo</a></p>\n",

		// Not considered safe
		"[foo](baz://bar/)\n",
		"<p><tt>foo</tt></p>\n",
	}
	doSafeTestsInline(t, tests)
}

func TestReferenceLink(t *testing.T) {
	var tests = []string{
		"[link][ref]\n",
		"<p>[link][ref]</p>\n",

		"[link][ref]\n   [ref]: /url/ \"title\"\n",
		"<p><a href=\"/url/\" title=\"title\">link</a></p>\n",

		"[link][ref]\n   [ref]: /url/\n",
		"<p><a href=\"/url/\">link</a></p>\n",

		"   [ref]: /url/\n",
		"",

		"   [ref]: /url/\n[ref2]: /url/\n [ref3]: /url/\n",
		"",

		"   [ref]: /url/\n[ref2]: /url/\n [ref3]: /url/\n    [4spaces]: /url/\n",
		"<pre><code>[4spaces]: /url/\n</code></pre>\n",

		"[hmm](ref2)\n   [ref]: /url/\n[ref2]: /url/\n [ref3]: /url/\n",
		"<p><a href=\"ref2\">hmm</a></p>\n",

		"[ref]\n",
		"<p>[ref]</p>\n",

		"[ref]\n   [ref]: /url/ \"title\"\n",
		"<p><a href=\"/url/\" title=\"title\">ref</a></p>\n",

		"[link][ref]\n   [ref]: /url/",
		"<p><a href=\"/url/\">link</a></p>\n",

		// Issue 172 in blackfriday
		"[]:<",
		"<p>[]:&lt;</p>\n",
	}
	doLinkTestsInline(t, tests)
}

func TestTags(t *testing.T) {
	var tests = []string{
		"a <span>tag</span>\n",
		"<p>a <span>tag</span></p>\n",

		"<span>tag</span>\n",
		"<p><span>tag</span></p>\n",

		"<span>mismatch</spandex>\n",
		"<p><span>mismatch</spandex></p>\n",

		"a <singleton /> tag\n",
		"<p>a <singleton /> tag</p>\n",
	}
	doTestsInline(t, tests)
}

func TestAutoLink(t *testing.T) {
	var tests = []string{
		"http://foo.com/\n",
		"<p><a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1 http://foo.com/\n",
		"<p>1 <a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1http://foo.com/\n",
		"<p>1<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1.http://foo.com/\n",
		"<p>1.<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"1. http://foo.com/\n",
		"<ol>\n<li><a href=\"http://foo.com/\">http://foo.com/</a></li>\n</ol>\n",

		"-http://foo.com/\n",
		"<p>-<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"- http://foo.com/\n",
		"<ul>\n<li><a href=\"http://foo.com/\">http://foo.com/</a></li>\n</ul>\n",

		"_http://foo.com/\n",
		"<p>_<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"令狐http://foo.com/\n",
		"<p>令狐<a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"令狐 http://foo.com/\n",
		"<p>令狐 <a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"ahttp://foo.com/\n",
		"<p>ahttp://foo.com/</p>\n",

		">http://foo.com/\n",
		"<blockquote>\n<p><a href=\"http://foo.com/\">http://foo.com/</a></p>\n</blockquote>\n",

		"> http://foo.com/\n",
		"<blockquote>\n<p><a href=\"http://foo.com/\">http://foo.com/</a></p>\n</blockquote>\n",

		"go to <http://foo.com/>\n",
		"<p>go to <a href=\"http://foo.com/\">http://foo.com/</a></p>\n",

		"a secure <https://link.org>\n",
		"<p>a secure <a href=\"https://link.org\">https://link.org</a></p>\n",

		"an email <mailto:some@one.com>\n",
		"<p>an email <a href=\"mailto:some@one.com\">some@one.com</a></p>\n",

		"an email <mailto://some@one.com>\n",
		"<p>an email <a href=\"mailto://some@one.com\">some@one.com</a></p>\n",

		"an email <some@one.com>\n",
		"<p>an email <a href=\"mailto:some@one.com\">some@one.com</a></p>\n",

		"an ftp <ftp://old.com>\n",
		"<p>an ftp <a href=\"ftp://old.com\">ftp://old.com</a></p>\n",

		"an ftp <ftp:old.com>\n",
		"<p>an ftp <a href=\"ftp:old.com\">ftp:old.com</a></p>\n",

		"a link with <http://new.com?query=foo&bar>\n",
		"<p>a link with <a href=\"http://new.com?query=foo&amp;bar\">" +
			"http://new.com?query=foo&amp;bar</a></p>\n",

		"quotes mean a tag <http://new.com?query=\"foo\"&bar>\n",
		"<p>quotes mean a tag <http://new.com?query=\"foo\"&bar></p>\n",

		"quotes mean a tag <http://new.com?query='foo'&bar>\n",
		"<p>quotes mean a tag <http://new.com?query='foo'&bar></p>\n",

		"unless escaped <http://new.com?query=\\\"foo\\\"&bar>\n",
		"<p>unless escaped <a href=\"http://new.com?query=&quot;foo&quot;&amp;bar\">" +
			"http://new.com?query=&quot;foo&quot;&amp;bar</a></p>\n",

		"even a > can be escaped <http://new.com?q=\\>&etc>\n",
		"<p>even a &gt; can be escaped <a href=\"http://new.com?q=&gt;&amp;etc\">" +
			"http://new.com?q=&gt;&amp;etc</a></p>\n",

		"<a href=\"http://fancy.com\">http://fancy.com</a>\n",
		"<p><a href=\"http://fancy.com\">http://fancy.com</a></p>\n",

		"<a href=\"http://fancy.com\">This is a link</a>\n",
		"<p><a href=\"http://fancy.com\">This is a link</a></p>\n",

		"<a href=\"http://www.fancy.com/A_B.pdf\">http://www.fancy.com/A_B.pdf</a>\n",
		"<p><a href=\"http://www.fancy.com/A_B.pdf\">http://www.fancy.com/A_B.pdf</a></p>\n",

		"(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (\n",
		"<p>(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (</p>\n",

		"(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (part two: <a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a>)).\n",
		"<p>(<a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a> (part two: <a href=\"http://www.fancy.com/A_B\">http://www.fancy.com/A_B</a>)).</p>\n",

		"http://www.foo.com<br>\n",
		"<p><a href=\"http://www.foo.com\">http://www.foo.com</a><br></p>\n",

		"http://foo.com/viewtopic.php?f=18&amp;t=297",
		"<p><a href=\"http://foo.com/viewtopic.php?f=18&amp;t=297\">http://foo.com/viewtopic.php?f=18&amp;t=297</a></p>\n",

		"http://foo.com/viewtopic.php?param=&quot;18&quot;zz",
		"<p><a href=\"http://foo.com/viewtopic.php?param=&quot;18&quot;zz\">http://foo.com/viewtopic.php?param=&quot;18&quot;zz</a></p>\n",

		"http://foo.com/viewtopic.php?param=&quot;18&quot;",
		"<p><a href=\"http://foo.com/viewtopic.php?param=&quot;18&quot;\">http://foo.com/viewtopic.php?param=&quot;18&quot;</a></p>\n",
	}
	doLinkTestsInline(t, tests)
}

var footnoteTests = []string{
	"testing footnotes.[^a]\n\n[^a]: This is the note\n",
	`<p>testing footnotes.<sup class="footnote-ref" id="fnref:a"><a class="footnote" href="#fn:a">1</a></sup></p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:a">This is the note
</li>
</ol>
</div>
`,

	`testing long[^b] notes.

[^b]: Paragraph 1

	Paragraph 2

	` + "```\n\tsome code\n\t```" + `

	Paragraph 3

No longer in the footnote
`,
	`<p>testing long<sup class="footnote-ref" id="fnref:b"><a class="footnote" href="#fn:b">1</a></sup> notes.</p>

<p>No longer in the footnote</p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:b"><p>Paragraph 1</p>

<p>Paragraph 2</p>

<p><code>some code</code></p>

<p>Paragraph 3</p>
</li>
</ol>
</div>
`,

	`testing[^c] multiple[^d] notes.

[^c]: this is [note] c


omg

[^d]: this is note d

what happens here

[note]: /link/c

`,
	`<p>testing<sup class="footnote-ref" id="fnref:c"><a class="footnote" href="#fn:c">1</a></sup> multiple<sup class="footnote-ref" id="fnref:d"><a class="footnote" href="#fn:d">2</a></sup> notes.</p>

<p>omg</p>

<p>what happens here</p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:c">this is <a href="/link/c">note</a> c
</li>
<li id="fn:d">this is note d
</li>
</ol>
</div>
`,

	"testing inline^[this is the note] notes.\n",
	`<p>testing inline<sup class="footnote-ref" id="fnref:this-is-the-note"><a class="footnote" href="#fn:this-is-the-note">1</a></sup> notes.</p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:this-is-the-note">this is the note</li>
</ol>
</div>
`,

	"testing multiple[^1] types^[inline note] of notes[^2]\n\n[^2]: the second deferred note\n[^1]: the first deferred note\n\n\twhich happens to be a block\n",
	`<p>testing multiple<sup class="footnote-ref" id="fnref:1"><a class="footnote" href="#fn:1">1</a></sup> types<sup class="footnote-ref" id="fnref:inline-note"><a class="footnote" href="#fn:inline-note">2</a></sup> of notes<sup class="footnote-ref" id="fnref:2"><a class="footnote" href="#fn:2">3</a></sup></p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:1"><p>the first deferred note</p>

<p>which happens to be a block</p>
</li>
<li id="fn:inline-note">inline note</li>
<li id="fn:2">the second deferred note
</li>
</ol>
</div>
`,

	`This is a footnote[^1]^[and this is an inline footnote]

[^1]: the footnote text.

    may be multiple paragraphs.
`,
	`<p>This is a footnote<sup class="footnote-ref" id="fnref:1"><a class="footnote" href="#fn:1">1</a></sup><sup class="footnote-ref" id="fnref:and-this-is-an-i"><a class="footnote" href="#fn:and-this-is-an-i">2</a></sup></p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:1"><p>the footnote text.</p>

<p>may be multiple paragraphs.</p>
</li>
<li id="fn:and-this-is-an-i">and this is an inline footnote</li>
</ol>
</div>
`,

	"empty footnote[^]\n\n[^]: fn text",
	"<p>empty footnote<sup class=\"footnote-ref\" id=\"fnref:\"><a class=\"footnote\" href=\"#fn:\">1</a></sup></p>\n<div class=\"footnotes\">\n\n<hr>\n\n<ol>\n<li id=\"fn:\">fn text\n</li>\n</ol>\n</div>\n",

	"Some text.[^note1]\n\n[^note1]: fn1",
	"<p>Some text.<sup class=\"footnote-ref\" id=\"fnref:note1\"><a class=\"footnote\" href=\"#fn:note1\">1</a></sup></p>\n<div class=\"footnotes\">\n\n<hr>\n\n<ol>\n<li id=\"fn:note1\">fn1\n</li>\n</ol>\n</div>\n",

	"Some text.[^note1][^note2]\n\n[^note1]: fn1\n[^note2]: fn2\n",
	"<p>Some text.<sup class=\"footnote-ref\" id=\"fnref:note1\"><a class=\"footnote\" href=\"#fn:note1\">1</a></sup><sup class=\"footnote-ref\" id=\"fnref:note2\"><a class=\"footnote\" href=\"#fn:note2\">2</a></sup></p>\n<div class=\"footnotes\">\n\n<hr>\n\n<ol>\n<li id=\"fn:note1\">fn1\n</li>\n<li id=\"fn:note2\">fn2\n</li>\n</ol>\n</div>\n",
}

func TestFootnotes(t *testing.T) {
	doTestsInlineParam(t, footnoteTests, EXTENSION_FOOTNOTES, 0, HtmlRendererParameters{})
}

func TestFootnotesWithParameters(t *testing.T) {
	tests := make([]string, len(footnoteTests))

	prefix := "testPrefix"
	returnText := "ret"
	re := regexp.MustCompile(`(?ms)<li id="fn:(\S+?)">(.*?)</li>`)

	// Transform the test expectations to match the parameters we're using.
	for i, test := range footnoteTests {
		if i%2 == 1 {
			test = strings.Replace(test, "fn:", "fn:"+prefix, -1)
			test = strings.Replace(test, "fnref:", "fnref:"+prefix, -1)
			test = re.ReplaceAllString(test, `<li id="fn:$1">$2 <a class="footnote-return" href="#fnref:$1">ret</a></li>`)
		}
		tests[i] = test
	}

	params := HtmlRendererParameters{
		FootnoteAnchorPrefix:       prefix,
		FootnoteReturnLinkContents: returnText,
	}

	doTestsInlineParam(t, tests, EXTENSION_FOOTNOTES, HTML_FOOTNOTE_RETURN_LINKS, params)
}

func runMarkdownInlineXML(input string, extensions, xmlFlags int) string {
	extensions |= EXTENSION_AUTOLINK
	extensions |= EXTENSION_CITATION
	extensions |= EXTENSION_SHORT_REF

	renderer := XmlRenderer(xmlFlags)

	return Parse([]byte(input), renderer, extensions).String()
}

func doTestsInlineXML(t *testing.T, tests []string) {
	doTestsInlineParamXML(t, tests, 0, 0)
}

func doTestsInlineParamXML(t *testing.T, tests []string, extensions, xmlFlags int) {
	var candidate string

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runMarkdownInlineXML(candidate, extensions, xmlFlags)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownInlineXML(candidate, extensions, xmlFlags)
				}
			}
		}
	}
}

func TestIndexXML(t *testing.T) {
	var tests = []string{
		"(((Tiger, Cats)))\n",
		"<t>\n<iref item=\"Tiger\" subitem=\"Cats\"/>\n</t>\n",

		"`(((Tiger, Cats)))`\n",
		"<t>\n<tt>(((Tiger, Cats)))</tt>\n</t>\n",

		"(((Tiger, Cats))\n",
		"<t>\n(((Tiger, Cats))\n</t>\n",
	}
	doTestsInlineXML(t, tests)
}

func TestNestedFootnotes(t *testing.T) {
	var tests = []string{
		`Paragraph.[^fn1]

[^fn1]:
  Asterisk[^fn2]

[^fn2]:
  Obelisk`,
		`<p>Paragraph.<sup class="footnote-ref" id="fnref:fn1"><a class="footnote" href="#fn:fn1">1</a></sup></p>
<div class="footnotes">

<hr>

<ol>
<li id="fn:fn1">Asterisk<sup class="footnote-ref" id="fnref:fn2"><a class="footnote" href="#fn:fn2">2</a></sup>
</li>
<li id="fn:fn2">Obelisk
</li>
</ol>
</div>
`,
	}
	doTestsInlineParam(t, tests, EXTENSION_FOOTNOTES, 0, HtmlRendererParameters{})
}

func TestInlineComments(t *testing.T) {
	var tests = []string{
		"Hello <!-- there ->\n",
		"<p>Hello &lt;!&mdash; there &ndash;&gt;</p>\n",

		"Hello <!-- there -->\n",
		"<p>Hello <!-- there --></p>\n",

		"Hello <!-- there -->",
		"<p>Hello <!-- there --></p>\n",

		"Hello <!---->\n",
		"<p>Hello <!----></p>\n",

		"Hello <!-- there -->\na",
		"<p>Hello <!-- there -->\na</p>\n",

		"* list <!-- item -->\n",
		"<ul>\n<li>list <!-- item --></li>\n</ul>\n",

		"<!-- Front --> comment\n",
		"<p><!-- Front --> comment</p>\n",

		"blahblah\n<!--- foo -->\nrhubarb\n",
		"<p>blahblah\n<!--- foo -->\nrhubarb</p>\n",
	}
	doTestsInlineParam(t, tests, 0, HTML_USE_SMARTYPANTS|HTML_SMARTYPANTS_DASHES, HtmlRendererParameters{})
}

func TestCitationXML(t *testing.T) {
	var tests = []string{
		"[@RFC2525]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"RFC2525\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2525.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@!RFC1024]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"RFC1024\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Normative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.1024.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@?RFC3024]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"RFC3024\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.3024.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[-@RFC3024]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.3024.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@?I-D.6man-udpzero]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"I-D.6man-udpzero\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml3/reference.I-D.6man-udpzero.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@?I-D.6man-udpzero#06]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"I-D.6man-udpzero\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml3/reference.I-D.draft-6man-udpzero-06.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@?I-D.6man-udpzero p. 23]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"I-D.6man-udpzero\" section=\"p. 23\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml3/reference.I-D.6man-udpzero.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@RFC2525], as you can see from @RFC2525, as @miekg says",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"RFC2525\"/>, as you can see from <xref target=\"RFC2525\"/>, as @miekg says\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2525.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@?RFC2525], as you can see from [@!RFC2525]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"RFC2525\"/>, as you can see from <xref target=\"RFC2525\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Normative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2525.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"[@!RFC2525], as you can see from [@?RFC2525]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\n<xref target=\"RFC2525\"/>, as you can see from <xref target=\"RFC2525\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Normative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2525.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"More of that in [@?ISO.8571-2.1988]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\nMore of that in <xref target=\"ISO.8571-2.1988\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Informative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml2/reference.ISO.8571-2.1988.xml\"/>\n</references>\n\n</back>\n</rfc>\n",

		"Hey, what about [@!ISO.8571-2.1988]",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>\nHey, what about <xref target=\"ISO.8571-2.1988\"/>\n</t>\n\n</middle>\n<back>\n<references>\n<name>Normative References</name>\n<xi:include href=\"https://xml2rfc.tools.ietf.org/public/rfc/bibxml2/reference.ISO.8571-2.1988.xml\"/>\n</references>\n\n</back>\n</rfc>\n",
	}
	doTestsInlineParamXML(t, tests, 0, XML_STANDALONE)
}

func TestRFC2119XML(t *testing.T) {
	var tests = []string{
		"MUST",
		"<t>\nMUST\n</t>\n",

		"*MUST*",
		"<t>\n<em>MUST</em>\n</t>\n",

		"**MUST**",
		"<t>\n<bcp14>MUST</bcp14>\n</t>\n",
	}
	doTestsInlineXML(t, tests)
}

func TestShortReferenceXML(t *testing.T) {
	var tests = []string{
		"(#ref)",
		"<t>\n<xref target=\"ref\"/>\n</t>\n",
	}
	doTestsInlineXML(t, tests)
}

func TestSmartDoubleQuotes(t *testing.T) {
	var tests = []string{
		"this should be normal \"quoted\" text.\n",
		"<p>this should be normal &ldquo;quoted&rdquo; text.</p>\n",
		"this \" single double\n",
		"<p>this &ldquo; single double</p>\n",
		"two pair of \"some\" quoted \"text\".\n",
		"<p>two pair of &ldquo;some&rdquo; quoted &ldquo;text&rdquo;.</p>\n"}

	doTestsInlineParam(t, tests, 0, HTML_USE_SMARTYPANTS, HtmlRendererParameters{})
}

func TestSmartAngledDoubleQuotes(t *testing.T) {
	var tests = []string{
		"this should be angled \"quoted\" text.\n",
		"<p>this should be angled &laquo;quoted&raquo; text.</p>\n",
		"this \" single double\n",
		"<p>this &laquo; single double</p>\n",
		"two pair of \"some\" quoted \"text\".\n",
		"<p>two pair of &laquo;some&raquo; quoted &laquo;text&raquo;.</p>\n"}

	doTestsInlineParam(t, tests, 0, HTML_USE_SMARTYPANTS|HTML_SMARTYPANTS_ANGLED_QUOTES, HtmlRendererParameters{})
}

func TestSmartFractions(t *testing.T) {
	var tests = []string{
		"1/2, 1/4 and 3/4; 1/4th and 3/4ths\n",
		"<p>&frac12;, &frac14; and &frac34;; &frac14;th and &frac34;ths</p>\n",
		"1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.\n",
		"<p>1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.</p>\n"}

	doTestsInlineParam(t, tests, 0, HTML_USE_SMARTYPANTS, HtmlRendererParameters{})

	tests = []string{
		"1/2, 2/3, 81/100 and 1000000/1048576.\n",
		"<p><sup>1</sup>&frasl;<sub>2</sub>, <sup>2</sup>&frasl;<sub>3</sub>, <sup>81</sup>&frasl;<sub>100</sub> and <sup>1000000</sup>&frasl;<sub>1048576</sub>.</p>\n",
		"1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.\n",
		"<p>1/2/2015, 1/4/2015, 3/4/2015; 2015/1/2, 2015/1/4, 2015/3/4.</p>\n"}

	doTestsInlineParam(t, tests, 0, HTML_USE_SMARTYPANTS|HTML_SMARTYPANTS_FRACTIONS, HtmlRendererParameters{})
}

func TestDisableSmartDashes(t *testing.T) {
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo - bar</p>\n",
		"foo -- bar\n",
		"<p>foo -- bar</p>\n",
		"foo --- bar\n",
		"<p>foo --- bar</p>\n",
	}, 0, 0, HtmlRendererParameters{})
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo &ndash; bar</p>\n",
		"foo -- bar\n",
		"<p>foo &mdash; bar</p>\n",
		"foo --- bar\n",
		"<p>foo &mdash;&ndash; bar</p>\n",
	}, 0, HTML_USE_SMARTYPANTS|HTML_SMARTYPANTS_DASHES, HtmlRendererParameters{})
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo - bar</p>\n",
		"foo -- bar\n",
		"<p>foo &ndash; bar</p>\n",
		"foo --- bar\n",
		"<p>foo &mdash; bar</p>\n",
	}, 0, HTML_USE_SMARTYPANTS|HTML_SMARTYPANTS_LATEX_DASHES|HTML_SMARTYPANTS_DASHES,
		HtmlRendererParameters{})
	doTestsInlineParam(t, []string{
		"foo - bar\n",
		"<p>foo - bar</p>\n",
		"foo -- bar\n",
		"<p>foo -- bar</p>\n",
		"foo --- bar\n",
		"<p>foo --- bar</p>\n",
	}, 0,
		HTML_USE_SMARTYPANTS|HTML_SMARTYPANTS_LATEX_DASHES,
		HtmlRendererParameters{})
}
