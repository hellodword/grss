package grss

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_RssFeed_001(t *testing.T) {
	// RSS 0.90
	// https://www.rssboard.org/rss-0-9-0
	s := `
<?xml version="1.0"?>
<rdf:RDF
xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
xmlns="http://channel.netscape.com/rdf/simple/0.9/">

  <channel>
    <title>Mozilla Dot Org</title>
    <link>http://www.mozilla.org</link>
    <description>the Mozilla Organization
      web site</description>
  </channel>

  <image>
    <title>Mozilla</title>
    <url>http://www.mozilla.org/images/moz.gif</url>
    <link>http://www.mozilla.org</link>
  </image>

  <item>
    <title>New Status Updates</title>
    <link>http://www.mozilla.org/status/</link>
  </item>

  <item>
    <title>Bugzilla Reorganized</title>
    <link>http://www.mozilla.org/bugs/</link>
  </item>

  <item>
    <title>Mozilla Party, 2.0!</title>
    <link>http://www.mozilla.org/party/1999/</link>
  </item>

  <item>
    <title>Unix Platform Parity</title>
    <link>http://www.mozilla.org/build/unix.html</link>
  </item>

  <item>
    <title>NPL 1.0M published</title>
    <link>http://www.mozilla.org/NPL/NPL-1.0M.html</link>
  </item>

</rdf:RDF>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rdf:RDF", a)
	assert.Equal(t, a.Attributes[0].Name.Local, "xmlns:rdf", a)

}

func Test_RssFeed_002(t *testing.T) {
	// 0.91 (Netscape) Example 1 - Simple
	// https://www.rssboard.org/rss-0-9-1-netscape
	s := `
<?xml version="1.0"?>
<!DOCTYPE rss SYSTEM "http://my.netscape.com/publish/formats/rss-0.91.dtd">
<rss version="0.91">
	<channel>
		<language>en</language>
		<description>News and commentary from the cross-platform scripting community.</description>
		<link>http://www.scripting.com/</link>
		<title>Scripting News</title>
		<image>
			<link>http://www.scripting.com/</link>
			<title>Scripting News</title>
			<url>http://www.scripting.com/gifs/tinyScriptingNews.gif</url>
		</image>
	</channel>
</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "0.91", a)
	assert.Equal(t, a.Channel.Language, "en", a)
	assert.Equal(t, a.Channel.Image.Title, "Scripting News", a)

}

func Test_RssFeed_003(t *testing.T) {
	// 0.91 (Netscape) Example 2 - Complete
	// https://www.rssboard.org/rss-0-9-1-netscape
	s := `
<?xml version="1.0"?>
<!DOCTYPE rss SYSTEM "http://my.netscape.com/publish/formats/rss-0.91.dtd">
<rss version="0.91">
  <channel>
    <copyright>Copyright 1997-1999 UserLand Software, Inc.</copyright>
    <pubDate>Thu, 08 Jul 1999 07:00:00 GMT</pubDate>
    <lastBuildDate>Thu, 08 Jul 1999 16:20:26 GMT</lastBuildDate>
    <docs>http://my.userland.com/stories/storyReader$11</docs>
    <description>News and commentary from the cross-platform scripting community.</description>
    <link>http://www.scripting.com/</link>
    <title>Scripting News</title>
    <image>
      <link>http://www.scripting.com/</link>
      <title>Scripting News</title>
      <url>http://www.scripting.com/gifs/tinyScriptingNews.gif</url>
      <height>40</height>
      <width>78</width>
      <description>What is this used for?</description>
    </image>
    <managingEditor>dave@userland.com (Dave Winer)</managingEditor>
    <webMaster>dave@userland.com (Dave Winer)</webMaster>
    <language>en-us</language>
    <skipHours>
      <hour>6</hour>
      <hour>7</hour>
      <hour>8</hour>
      <hour>9</hour>
      <hour>10</hour>
      <hour>11</hour>
    </skipHours>
    <skipDays>
      <day>Sunday</day>
    </skipDays>
    <rating>(PICS-1.1 "http://www.rsac.org/ratingsv01.html" l gen true comment "RSACi North America Server" for "http://www.rsac.org" on "1996.04.16T08:15-0500" r (n 0 s 0 v 0 l 0))</rating>
    <item>
      <title>stuff</title>
      <link>http://bar</link>
      <description>This is an article about some stuff</description>
    </item>
    <textinput>
      <title>Search Now!</title>
      <description>Enter your search terms</description>
      <name>find</name>
      <link>http://my.site.com/search.cgi</link>
    </textinput>
  </channel>
</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "0.91", a)
	assert.Equal(t, a.Channel.Copyright, "Copyright 1997-1999 UserLand Software, Inc.", a)
	assert.Equal(t, a.Channel.Image.Height, "40", a)
	assert.Equal(t, a.Channel.Image.Description, "What is this used for?", a)
	assert.Equal(t, a.Channel.SkipHours.Hours[5], "11", a)
	assert.Equal(t, a.Channel.SkipDays.Days[0], "Sunday", a)
	assert.Equal(t, string(a.Channel.Rating.Content), "(PICS-1.1 \"http://www.rsac.org/ratingsv01.html\" l gen true comment \"RSACi North America Server\" for \"http://www.rsac.org\" on \"1996.04.16T08:15-0500\" r (n 0 s 0 v 0 l 0))", a)
	assert.Equal(t, a.Channel.Items[0].Link, "http://bar", a)
	assert.Equal(t, a.Channel.TextInput.Title, "Search Now!", a)

}

func Test_RssFeed_004(t *testing.T) {
	// 0.91 (Netscape) Example 3 - International
	// https://www.rssboard.org/rss-0-9-1-netscape
	s := `
<?xml version="1.0" encoding="EuC-JP"?>

<!DOCTYPE rss SYSTEM "http://my.netscape.com/publish/formats/rss-0.91.dtd">
<rss version="0.91">

<channel>
  <title> ... </title>
  <link>http://www.mozilla.org</link>
  <description> ... </description>
  <language>ja</language> <!-- tagged as Japanese content -->

  <item>
    <title> ... </title>
    <link>http://www.mozilla.org/status/</link>
    <description>This is an item description...</description>
  </item>

  <item>
    <title> ... </title>
    <link>http://www.mozilla.org/status/</link>
    <description>This is an item description...</description>
  </item>

  <item>
    <title> ... </title>
    <link>http://www.mozilla.org/status/</link>
    <description>This is an item description...</description>
  </item>
  <item>
    <title> ... </title>
    <link>http://www.mozilla.org/status/</link>
    <description>This is an item description...</description>
  </item>

</channel>
</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.Charset, "EuC-JP", a)
	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "0.91", a)
	assert.Equal(t, a.Channel.Language, "ja", a)
	assert.Equal(t, a.Channel.Items[3].Link, "http://www.mozilla.org/status/", a)

}

func Test_RssFeed_005(t *testing.T) {
	// 0.91
	// http://static.userland.com/gems/backend/sampleRss.xml
	s := `

<?xml version="1.0" encoding="ISO-8859-1" ?>
<rss version="0.91">
	<channel>
		<title>WriteTheWeb</title> 
		<link>http://writetheweb.com</link> 
		<description>News for web users that write back</description> 
		<language>en-us</language> 
		<copyright>Copyright 2000, WriteTheWeb team.</copyright> 
		<managingEditor>editor@writetheweb.com</managingEditor> 
		<webMaster>webmaster@writetheweb.com</webMaster> 
		<image>
			<title>WriteTheWeb</title> 
			<url>http://writetheweb.com/images/mynetscape88.gif</url> 
			<link>http://writetheweb.com</link> 
			<width>88</width> 
			<height>31</height> 
			<description>News for web users that write back</description> 
			</image>
		<item>
			<title>Giving the world a pluggable Gnutella</title> 
			<link>http://writetheweb.com/read.php?item=24</link> 
			<description>WorldOS is a framework on which to build programs that work like Freenet or Gnutella -allowing distributed applications using peer-to-peer routing.</description> 
			</item>
		<item>
			<title>Syndication discussions hot up</title> 
			<link>http://writetheweb.com/read.php?item=23</link> 
			<description>After a period of dormancy, the Syndication mailing list has become active again, with contributions from leaders in traditional media and Web syndication.</description> 
			</item>
		<item>
			<title>Personal web server integrates file sharing and messaging</title> 
			<link>http://writetheweb.com/read.php?item=22</link> 
			<description>The Magi Project is an innovative project to create a combined personal web server and messaging system that enables the sharing and synchronization of information across desktop, laptop and palmtop devices.</description> 
			</item>
		<item>
			<title>Syndication and Metadata</title> 
			<link>http://writetheweb.com/read.php?item=21</link> 
			<description>RSS is probably the best known metadata format around. RDF is probably one of the least understood. In this essay, published on my O'Reilly Network weblog, I argue that the next generation of RSS should be based on RDF.</description> 
			</item>
		<item>
			<title>UK bloggers get organised</title> 
			<link>http://writetheweb.com/read.php?item=20</link> 
			<description>Looks like the weblogs scene is gathering pace beyond the shores of the US. There's now a UK-specific page on weblogs.com, and a mailing list at egroups.</description> 
			</item>
		<item>
			<title>Yournamehere.com more important than anything</title> 
			<link>http://writetheweb.com/read.php?item=19</link> 
			<description>Whatever you're publishing on the web, your site name is the most valuable asset you have, according to Carl Steadman.</description> 
			</item>
		</channel>
	</rss>

`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.Charset, "ISO-8859-1", a)
	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "0.91", a)
	assert.Equal(t, len(a.Channel.Items), 6, a)

}

func Test_RssFeed_006(t *testing.T) {
	// RSS 0.92
	// http://static.userland.com/gems/backend/gratefulDead.xml
	s := `
<?xml version="1.0"?>
<!-- RSS generation done by 'Radio UserLand' on Fri, 13 Apr 2001 19:23:02 GMT -->
<rss version="0.92">
	<channel>
		<title>Dave Winer: Grateful Dead</title>
		<link>http://www.scripting.com/blog/categories/gratefulDead.html</link>
		<description>A high-fidelity Grateful Dead song every day. This is where we&apos;re experimenting with enclosures on RSS news items that download when you&apos;re not using your computer. If it works (it will) it will be the end of the Click-And-Wait multimedia experience on the Internet. </description>
		<lastBuildDate>Fri, 13 Apr 2001 19:23:02 GMT</lastBuildDate>
		<docs>http://backend.userland.com/rss092</docs>
		<managingEditor>dave@userland.com (Dave Winer)</managingEditor>
		<webMaster>dave@userland.com (Dave Winer)</webMaster>
		<cloud domain="data.ourfavoritesongs.com" port="80" path="/RPC2" registerProcedure="ourFavoriteSongs.rssPleaseNotify" protocol="xml-rpc"/>
		<item>
			<description>It&apos;s been a few days since I added a song to the Grateful Dead channel. Now that there are all these new Radio users, many of whom are tuned into this channel (it&apos;s #16 on the hotlist of upstreaming Radio users, there&apos;s no way of knowing how many non-upstreaming users are subscribing, have to do something about this..). Anyway, tonight&apos;s song is a live version of Weather Report Suite from Dick&apos;s Picks Volume 7. It&apos;s wistful music. Of course a beautiful song, oft-quoted here on Scripting News. &lt;i&gt;A little change, the wind and rain.&lt;/i&gt;
</description>
			<enclosure url="http://www.scripting.com/mp3s/weatherReportDicksPicsVol7.mp3" length="6182912" type="audio/mpeg"/>
			</item>
		<item>
			<description>Kevin Drennan started a &lt;a href=&quot;http://deadend.editthispage.com/&quot;&gt;Grateful Dead Weblog&lt;/a&gt;. Hey it&apos;s cool, he even has a &lt;a href=&quot;http://deadend.editthispage.com/directory/61&quot;&gt;directory&lt;/a&gt;. &lt;i&gt;A Frontier 7 feature.&lt;/i&gt;</description>
			<source url="http://scriptingnews.userland.com/xml/scriptingNews2.xml">Scripting News</source>
			</item>
		<item>
			<description>&lt;a href=&quot;http://arts.ucsc.edu/GDead/AGDL/other1.html&quot;&gt;The Other One&lt;/a&gt;, live instrumental, One From The Vault. Very rhythmic very spacy, you can listen to it many times, and enjoy something new every time.</description>
			<enclosure url="http://www.scripting.com/mp3s/theOtherOne.mp3" length="6666097" type="audio/mpeg"/>
			</item>
		<item>
			<description>This is a test of a change I just made. Still diggin..</description>
			</item>
		<item>
			<description>The HTML rendering almost &lt;a href=&quot;http://validator.w3.org/check/referer&quot;&gt;validates&lt;/a&gt;. Close. Hey I wonder if anyone has ever published a style guide for ALT attributes on images? What are you supposed to say in the ALT attribute? I sure don&apos;t know. If you&apos;re blind send me an email if u cn rd ths. </description>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.cs.cmu.edu/~mleone/gdead/dead-lyrics/Franklin&apos;s_Tower.txt&quot;&gt;Franklin&apos;s Tower&lt;/a&gt;, a live version from One From The Vault.</description>
			<enclosure url="http://www.scripting.com/mp3s/franklinsTower.mp3" length="6701402" type="audio/mpeg"/>
			</item>
		<item>
			<description>Moshe Weitzman says Shakedown Street is what I&apos;m lookin for for tonight. I&apos;m listening right now. It&apos;s one of my favorites. &quot;Don&apos;t tell me this town ain&apos;t got no heart.&quot; Too bright. I like the jazziness of Weather Report Suite. Dreamy and soft. How about The Other One? &quot;Spanish lady come to me..&quot;</description>
			<source url="http://scriptingnews.userland.com/xml/scriptingNews2.xml">Scripting News</source>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.scripting.com/mp3s/youWinAgain.mp3&quot;&gt;The news is out&lt;/a&gt;, all over town..&lt;p&gt;
You&apos;ve been seen, out runnin round. &lt;p&gt;
The lyrics are &lt;a href=&quot;http://www.cs.cmu.edu/~mleone/gdead/dead-lyrics/You_Win_Again.txt&quot;&gt;here&lt;/a&gt;, short and sweet. &lt;p&gt;
&lt;i&gt;You win again!&lt;/i&gt;
</description>
			<enclosure url="http://www.scripting.com/mp3s/youWinAgain.mp3" length="3874816" type="audio/mpeg"/>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.getlyrics.com/lyrics/grateful-dead/wake-of-the-flood/07.htm&quot;&gt;Weather Report Suite&lt;/a&gt;: &quot;Winter rain, now tell me why, summers fade, and roses die? The answer came. The wind and rain. Golden hills, now veiled in grey, summer leaves have blown away. Now what remains? The wind and rain.&quot;</description>
			<enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" type="audio/mpeg"/>
			</item>
		<item>
			<description>&lt;a href=&quot;http://arts.ucsc.edu/gdead/agdl/darkstar.html&quot;&gt;Dark Star&lt;/a&gt; crashes, pouring its light into ashes.</description>
			<enclosure url="http://www.scripting.com/mp3s/darkStar.mp3" length="10889216" type="audio/mpeg"/>
			</item>
		<item>
			<description>DaveNet: &lt;a href=&quot;http://davenet.userland.com/2001/01/21/theUsBlues&quot;&gt;The U.S. Blues&lt;/a&gt;.</description>
			</item>
		<item>
			<description>Still listening to the US Blues. &lt;i&gt;&quot;Wave that flag, wave it wide and high..&quot;&lt;/i&gt; Mistake made in the 60s. We gave our country to the assholes. Ah ah. Let&apos;s take it back. Hey I&apos;m still a hippie. &lt;i&gt;&quot;You could call this song The United States Blues.&quot;&lt;/i&gt;</description>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.sixties.com/html/garcia_stack_0.html&quot;&gt;&lt;img src=&quot;http://www.scripting.com/images/captainTripsSmall.gif&quot; height=&quot;51&quot; width=&quot;42&quot; border=&quot;0&quot; hspace=&quot;10&quot; vspace=&quot;10&quot; align=&quot;right&quot;&gt;&lt;/a&gt;In celebration of today&apos;s inauguration, after hearing all those great patriotic songs, America the Beautiful, even The Star Spangled Banner made my eyes mist up. It made my choice of Grateful Dead song of the night realllly easy. Here are the &lt;a href=&quot;http://searchlyrics2.homestead.com/gd_usblues.html&quot;&gt;lyrics&lt;/a&gt;. Click on the audio icon to the left to give it a listen. &quot;Red and white, blue suede shoes, I&apos;m Uncle Sam, how do you do?&quot; It&apos;s a different kind of patriotic music, but man I love my country and I love Jerry and the band. &lt;i&gt;I truly do!&lt;/i&gt;</description>
			<enclosure url="http://www.scripting.com/mp3s/usBlues.mp3" length="5272510" type="audio/mpeg"/>
			</item>
		<item>
			<description>Grateful Dead: &quot;Tennessee, Tennessee, ain&apos;t no place I&apos;d rather be.&quot;</description>
			<enclosure url="http://www.scripting.com/mp3s/tennesseeJed.mp3" length="3442648" type="audio/mpeg"/>
			</item>
		<item>
			<description>Ed Cone: &quot;Had a nice Deadhead experience with my wife, who never was one but gets the vibe and knows and likes a lot of the music. Somehow she made it to the age of 40 without ever hearing Wharf Rat. We drove to Jersey and back over Christmas with the live album commonly known as Skull and Roses in the CD player much of the way, and it was cool to see her discover one the band&apos;s finest moments. That song is unique and underappreciated. Fun to hear that disc again after a few years off -- you get Jerry as blues-guitar hero on Big Railroad Blues and a nice version of Bertha.&quot;</description>
			<enclosure url="http://www.scripting.com/mp3s/darkStarWharfRat.mp3" length="27503386" type="audio/mpeg"/>
			</item>
		<item>
			<description>&lt;a href=&quot;http://arts.ucsc.edu/GDead/AGDL/fotd.html&quot;&gt;Tonight&apos;s Song&lt;/a&gt;: &quot;If I get home before daylight I just might get some sleep tonight.&quot; </description>
			<enclosure url="http://www.scripting.com/mp3s/friendOfTheDevil.mp3" length="3219742" type="audio/mpeg"/>
			</item>
		<item>
			<description>&lt;a href=&quot;http://arts.ucsc.edu/GDead/AGDL/uncle.html&quot;&gt;Tonight&apos;s song&lt;/a&gt;: &quot;Come hear Uncle John&apos;s Band by the river side. Got some things to talk about here beside the rising tide.&quot;</description>
			<enclosure url="http://www.scripting.com/mp3s/uncleJohnsBand.mp3" length="4587102" type="audio/mpeg"/>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.cs.cmu.edu/~mleone/gdead/dead-lyrics/Me_and_My_Uncle.txt&quot;&gt;Me and My Uncle&lt;/a&gt;: &quot;I loved my uncle, God rest his soul, taught me good, Lord, taught me all I know. Taught me so well, I grabbed that gold and I left his dead ass there by the side of the road.&quot;
</description>
			<enclosure url="http://www.scripting.com/mp3s/meAndMyUncle.mp3" length="2949248" type="audio/mpeg"/>
			</item>
		<item>
			<description>Truckin, like the doo-dah man, once told me gotta play your hand. Sometimes the cards ain&apos;t worth a dime, if you don&apos;t lay em down.</description>
			<enclosure url="http://www.scripting.com/mp3s/truckin.mp3" length="4847908" type="audio/mpeg"/>
			</item>
		<item>
			<description>Two-Way-Web: &lt;a href=&quot;http://www.thetwowayweb.com/payloadsForRss&quot;&gt;Payloads for RSS&lt;/a&gt;. &quot;When I started talking with Adam late last year, he wanted me to think about high quality video on the Internet, and I totally didn&apos;t want to hear about it.&quot;</description>
			</item>
		<item>
			<description>A touch of gray, kinda suits you anyway..</description>
			<enclosure url="http://www.scripting.com/mp3s/touchOfGrey.mp3" length="5588242" type="audio/mpeg"/>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.sixties.com/html/garcia_stack_0.html&quot;&gt;&lt;img src=&quot;http://www.scripting.com/images/captainTripsSmall.gif&quot; height=&quot;51&quot; width=&quot;42&quot; border=&quot;0&quot; hspace=&quot;10&quot; vspace=&quot;10&quot; align=&quot;right&quot;&gt;&lt;/a&gt;In celebration of today&apos;s inauguration, after hearing all those great patriotic songs, America the Beautiful, even The Star Spangled Banner made my eyes mist up. It made my choice of Grateful Dead song of the night realllly easy. Here are the &lt;a href=&quot;http://searchlyrics2.homestead.com/gd_usblues.html&quot;&gt;lyrics&lt;/a&gt;. Click on the audio icon to the left to give it a listen. &quot;Red and white, blue suede shoes, I&apos;m Uncle Sam, how do you do?&quot; It&apos;s a different kind of patriotic music, but man I love my country and I love Jerry and the band. &lt;i&gt;I truly do!&lt;/i&gt;</description>
			<enclosure url="http://www.scripting.com/mp3s/usBlues.mp3" length="5272510" type="audio/mpeg"/>
			</item>
		</channel>
	</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "0.92", a)
	assert.Equal(t, a.Channel.Title, "Dave Winer: Grateful Dead", a.Channel)
	assert.Equal(t, a.Channel.Link, "http://www.scripting.com/blog/categories/gratefulDead.html", a.Channel)
	assert.Equal(t, a.Channel.Description, "A high-fidelity Grateful Dead song every day. This is where we're experimenting with enclosures on RSS news items that download when you're not using your computer. If it works (it will) it will be the end of the Click-And-Wait multimedia experience on the Internet. ", a.Channel)
	assert.Equal(t, a.Channel.Cloud.Attributes[0].Value, "data.ourfavoritesongs.com", a.Channel.Cloud)
	assert.Equal(t, a.Channel.Cloud.Attributes[4].Name.Local, "protocol", a.Channel.Cloud)
	assert.Equal(t, len(a.Channel.Items), 22, a.Channel.Items)
	assert.Equal(t, a.Channel.Items[0].Enclosure.Type, "audio/mpeg", a.Channel.Items[0].Enclosure)
	assert.Equal(t, a.Channel.Items[1].Source.Url, "http://scriptingnews.userland.com/xml/scriptingNews2.xml", a.Channel.Items[1].Source)
	assert.Equal(t, a.Channel.Items[1].Source.Text, "Scripting News", a.Channel.Items[1].Source)
	assert.Equal(t, a.Channel.Items[21].Enclosure.Length, "5272510", a.Channel.Items[21].Enclosure)

}

func Test_RssFeed_007(t *testing.T) {
	// RSS 2.0
	// http://static.userland.com/gems/backend/rssTwoExample2.xml
	s := `

<?xml version="1.0"?>
<!-- RSS generated by Radio UserLand v8.0.5 on 9/30/2002; 4:00:00 AM Pacific -->
<rss version="2.0" xmlns:blogChannel="http://backend.userland.com/blogChannelModule">
	<channel>
		<title>Scripting News</title>
		<link>http://www.scripting.com/</link>
		<description>A weblog about scripting and stuff like that.</description>
		<language>en-us</language>
		<blogChannel:blogRoll>http://radio.weblogs.com/0001015/userland/scriptingNewsLeftLinks.opml</blogChannel:blogRoll>
		<blogChannel:mySubscriptions>http://radio.weblogs.com/0001015/gems/mySubscriptions.opml</blogChannel:mySubscriptions>
		<blogChannel:blink>http://diveintomark.org/</blogChannel:blink>
		<copyright>Copyright 1997-2002 Dave Winer</copyright>
		<lastBuildDate>Mon, 30 Sep 2002 11:00:00 GMT</lastBuildDate>
		<docs>http://backend.userland.com/rss</docs>
		<generator>Radio UserLand v8.0.5</generator>
		<category domain="Syndic8">1765</category>
		<managingEditor>dave@userland.com</managingEditor>
		<webMaster>dave@userland.com</webMaster>
		<ttl>40</ttl>
		<item>
			<description>&quot;rssflowersalignright&quot;With any luck we should have one or two more days of namespaces stuff here on Scripting News. It feels like it's winding down. Later in the week I'm going to a &lt;a href=&quot;http://harvardbusinessonline.hbsp.harvard.edu/b02/en/conferences/conf_detail.jhtml?id=s775stg&amp;pid=144XCF&quot;&gt;conference&lt;/a&gt; put on by the Harvard Business School. So that should change the topic a bit. The following week I'm off to Colorado for the &lt;a href=&quot;http://www.digitalidworld.com/conference/2002/index.php&quot;&gt;Digital ID World&lt;/a&gt; conference. We had to go through namespaces, and it turns out that weblogs are a great way to work around mail lists that are clogged with &lt;a href=&quot;http://www.userland.com/whatIsStopEnergy&quot;&gt;stop energy&lt;/a&gt;. I think we solved the problem, have reached a consensus, and will be ready to move forward shortly.</description>
			<pubDate>Mon, 30 Sep 2002 01:56:02 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#When:6:56:02PM</guid>
			</item>
		<item>
			<description>Joshua Allen: &lt;a href=&quot;http://www.netcrucible.com/blog/2002/09/29.html#a243&quot;&gt;Who loves namespaces?&lt;/a&gt;</description>
			<pubDate>Sun, 29 Sep 2002 19:59:01 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#When:12:59:01PM</guid>
			</item>
		<item>
			<description>&lt;a href=&quot;http://www.docuverse.com/blog/donpark/2002/09/29.html#a68&quot;&gt;Don Park&lt;/a&gt;: &quot;It is too easy for engineer to anticipate too much and XML Namespace is a frequent host of over-anticipation.&quot;</description>
			<pubDate>Mon, 30 Sep 2002 01:52:02 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#When:6:52:02PM</guid>
			</item>
		<item>
			<description>&lt;a href=&quot;http://scriptingnews.userland.com/stories/storyReader$1768&quot;&gt;Three Sunday Morning Options&lt;/a&gt;. &quot;I just got off the phone with Tim Bray, who graciously returned my call on a Sunday morning while he was making breakfast for his kids.&quot; We talked about three options for namespaces in RSS 2.0, and I think I now have the tradeoffs well outlined, and ready for other developers to review. If there is now a consensus, I think we can easily move forward. </description>
			<pubDate>Sun, 29 Sep 2002 17:05:20 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#When:10:05:20AM</guid>
			</item>
		<item>
			<description>&lt;a href=&quot;http://blog.mediacooperative.com/mt-comments.cgi?entry_id=1435&quot;&gt;Mark Pilgrim&lt;/a&gt; weighs in behind option 1 on a Ben Hammersley thread. On the RSS2-Support list, Phil Ringnalda lists a set of &lt;a href=&quot;http://groups.yahoo.com/group/RSS2-Support/message/54&quot;&gt;proposals&lt;/a&gt;, the first is equivalent to option 1. </description>
			<pubDate>Sun, 29 Sep 2002 19:09:28 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#When:12:09:28PM</guid>
			</item>
		<item>
			<description>&lt;a href=&quot;http://effbot.org/zone/effnews-4.htm&quot;&gt;Fredrik Lundh breaks&lt;/a&gt; through, following Simon Fell's lead, now his Python aggregator works with Scripting News &lt;a href=&quot;http://www.scripting.com/rss.xml&quot;&gt;in&lt;/a&gt; RSS 2.0. BTW, the spec is imperfect in regards to namespaces. We anticipated a 2.0.1 and 2.0.2 in the Roadmap for exactly this purpose. Thanks for your help, as usual, Fredrik. </description>
			<pubDate>Sun, 29 Sep 2002 15:01:02 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#When:8:01:02AM</guid>
			</item>
		<item>
			<title>Law and Order</title>
			<link>http://scriptingnews.userland.com/backissues/2002/09/29#lawAndOrder</link>
			<description>
				&lt;p&gt;&lt;a href=&quot;http://www.nbc.com/Law_&amp;_Order/index.html&quot;&gt;&lt;img src=&quot;http://radio.weblogs.com/0001015/images/2002/09/29/lenny.gif&quot; width=&quot;45&quot; height=&quot;53&quot; border=&quot;0&quot; align=&quot;right&quot; hspace=&quot;15&quot; vspace=&quot;5&quot; alt=&quot;A picture named lenny.gif&quot;&gt;&lt;/a&gt;A great line in a recent Law and Order. Lenny Briscoe, played by Jerry Orbach, is interrogating a suspect. The suspect tells a story and reaches a point where no one believes him, not even the suspect himself. Lenny says: &quot;Now there's five minutes of my life that's lost forever.&quot; &lt;/p&gt;
				</description>
			<pubDate>Sun, 29 Sep 2002 23:48:33 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#lawAndOrder</guid>
			</item>
		<item>
			<title>Rule 1</title>
			<link>http://scriptingnews.userland.com/backissues/2002/09/29#rule1</link>
			<description>
				&lt;p&gt;In the discussions over namespaces in RSS 2.0, one thing I hear a lot of, that is just plain wrong, is that when you move up by a major version number, breakage is expected and is okay. In the world I come from it is, emphatically, &lt;i&gt;not okay.&lt;/i&gt; We spend huge resources to make sure that files, scripts and apps built in version N work in version N+1 without modification. Even the smallest change in the core engine can break apps. It's just not acceptable. When we make changes we have to be sure there's no breakage. I don't know where these other people come from, or if they make software that anyone uses, but the users I know don't stand for that. As we expose the tradeoffs it becomes clear that &lt;i&gt;that's the issue here.&lt;/i&gt; We are not in Year Zero. There are users. Breaking them is not an option. A conclusion to lift the confusion: Version 0.91 and 0.92 files are valid 2.0 files. This is where we started, what seems like years ago.&lt;/p&gt;
				&lt;p&gt;BTW, you can ask anyone who's worked for me in a technical job to explain rules 1 and 1b. (I'll clue you in. Rule 1 is &quot;No Breakage&quot; and Rule 1b is &quot;Don't Break Dave.&quot;)&lt;/p&gt;
				</description>
			<pubDate>Sun, 29 Sep 2002 17:24:20 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#rule1</guid>
			</item>
		<item>
			<title>Really early morning no-coffee notes</title>
			<link>http://scriptingnews.userland.com/backissues/2002/09/29#reallyEarlyMorningNocoffeeNotes</link>
			<description>
				&lt;p&gt;One of the lessons I've learned in 47.4 years: When someone accuses you of a &lt;a href=&quot;http://www.dictionary.com/search?q=deceit&quot;&gt;deceit&lt;/a&gt;, there's a very good chance the accuser practices that form of deceit, and a reasonable chance that he or she is doing it as they point the finger. &lt;/p&gt;
				&lt;p&gt;&lt;a href=&quot;http://www.docuverse.com/blog/donpark/2002/09/28.html#a66&quot;&gt;Don Park&lt;/a&gt;: &quot;He poured a barrel full of pig urine all over the Korean Congress because he was pissed off about all the dirty politics going on.&quot;&lt;/p&gt;
				&lt;p&gt;&lt;a href=&quot;http://davenet.userland.com/1995/01/04/demoingsoftwareforfunprofi&quot;&gt;1/4/95&lt;/a&gt;: &quot;By the way, the person with the big problem is probably a competitor.&quot;&lt;/p&gt;
				&lt;p&gt;I've had a fair amount of experience in the last few years with what you might call standards work. XML-RPC, SOAP, RSS, OPML. Each has been different from the others. In all this work, the most positive experience was XML-RPC, and not just because of the technical excellence of the people involved. In the end, what matters more to me is &lt;a href=&quot;http://www.dictionary.com/search?q=collegiality&quot;&gt;collegiality&lt;/a&gt;. Working together, person to person, for the sheer pleasure of it, is even more satisfying than a good technical result. Now, getting both is the best, and while XML-RPC is not perfect, it's pretty good. I also believe that if you have collegiality, technical excellence follows as a natural outcome.&lt;/p&gt;
				&lt;p&gt;One more bit of philosophy. At my checkup earlier this week, one of the things my cardiologist asked was if I was experiencing any kind of intellectual dysfunction. In other words, did I lose any of my sharpness as a result of the surgery in June. I told him yes I had and thanked him for asking. In an amazing bit of synchronicity, the next day John Robb &lt;a href=&quot;http://jrobb.userland.com/2002/09/25.html#a2598&quot;&gt;located&lt;/a&gt; an article in New Scientist that said that scientists had found a way to prevent this from happening. I hadn't talked with John about my experience or the question the doctor asked. Yesterday I was telling the story to my friend Dave Jacobs. He said it's not a problem because I always had excess capacity in that area. Exactly right Big Dave and thanks for the vote of confidence.&lt;/p&gt;
				</description>
			<pubDate>Sun, 29 Sep 2002 11:13:10 GMT</pubDate>
			<guid>http://scriptingnews.userland.com/backissues/2002/09/29#reallyEarlyMorningNocoffeeNotes</guid>
			</item>
		</channel>
	</rss>

`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "2.0", a)
	assert.Equal(t, a.Attributes[0].Name.Local, "xmlns:blogChannel", a)
	assert.Equal(t, a.Channel.ExtensionElement[2].XMLName.Local, "blogChannel:blink", a.Channel.ExtensionElement)
	assert.Equal(t, a.Channel.Ttl, "40", a.Channel)
	assert.Equal(t, len(a.Channel.Items), 9, a.Channel)
	assert.Equal(t, a.Channel.Items[8].Guid.Guid, "http://scriptingnews.userland.com/backissues/2002/09/29#reallyEarlyMorningNocoffeeNotes", a.Channel)
	assert.Equal(t, a.Channel.Items[8].PubDate, "Sun, 29 Sep 2002 11:13:10 GMT", a.Channel)

}

func Test_RssFeed_008(t *testing.T) {
	// content:encoded
	s := `
<?xml version="1.0" encoding="utf-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
<channel xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <title>喔啊 X Media</title>
  <link>https://example.com/</link>
  <description>X media article news feed.</description>
  <atom:link href="https://example.com/newsfeed/" rel="self"></atom:link>
  <language>zh-hant</language>
  <lastBuildDate>Thu, 10 Nov 2022 10:01:00 +0800</lastBuildDate>
  <link href="https://example.com/newsfeed/?page=1" rel="next"></link>

  <item>
    <title>喔喔喔</title>
    <link>https://example.com/article/example-morning-brief/</link>
    <description>喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔在12月6喔喔喔喔喔</description>
    <pubDate>Thu, 10 Nov 2022 10:01:00 +0800</pubDate>
    <guid>https://example.com/article/example-morning-brief/</guid>
    <category>喔</category>
    <content:encoded>&lt;section&gt;  &lt;article&gt;&lt;h3&gt;喔喔喔喔喔喔喔喔喔喔喔喔期&lt;/h3&gt;&lt;p&gt;喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔12月6喔喔喔喔喔喔喔喔。&lt;/p&gt;&lt;p&gt;&lt;a href="https://www.example.com/politics/2022/11/09/example-news-live-updates/#link-example"&gt;喔喔喔喔喔&lt;/a&gt;喔喔喔喔喔喔喔喔喔喔去40喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔1986喔喔喔喔喔喔&lt;/p&gt;&lt;p&gt;喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔。&lt;/p&gt;&lt;p&gt;喔&lt;a href="https://www.example.com/example/2022/results/example?example-data-id=2022-SG&gt;</content:encoded>
  </item>
</channel>
</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "2.0", a)
	assert.Equal(t, a.Channel.Items[0].Title, "喔喔喔", a)
	assert.Equal(t, a.Channel.Items[0].Content.XMLName.Local, "content:encoded", a.Channel.Items[0].Content)
	assert.Equal(t, string(a.Channel.Items[0].Content.Content), "&lt;section&gt;  &lt;article&gt;&lt;h3&gt;喔喔喔喔喔喔喔喔喔喔喔喔期&lt;/h3&gt;&lt;p&gt;喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔12月6喔喔喔喔喔喔喔喔。&lt;/p&gt;&lt;p&gt;&lt;a href=\"https://www.example.com/politics/2022/11/09/example-news-live-updates/#link-example\"&gt;喔喔喔喔喔&lt;/a&gt;喔喔喔喔喔喔喔喔喔喔去40喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔1986喔喔喔喔喔喔&lt;/p&gt;&lt;p&gt;喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔喔。&lt;/p&gt;&lt;p&gt;喔&lt;a href=\"https://www.example.com/example/2022/results/example?example-data-id=2022-SG&gt;", a.Channel.Items[0].Content)

}

func Test_RssFeed_009(t *testing.T) {
	// ISO-8859-1
	// https://github.com/SlyMarbo/rss/blob/9ae0f45449d6f6424a61969aeedeb14cd5d40094/testdata/rss_0.91
	s := `
<?xml version="1.0" encoding="ISO-8859-1" ?>
<rss version="0.91">
	<channel>
		<title>WriteTheWeb</title> 
		<link>http://writetheweb.com</link> 
		<description>News for web users that write back</description> 
		<language>en-us</language> 
		<copyright>Copyright 2000, WriteTheWeb team.</copyright> 
		<managingEditor>editor@writetheweb.com</managingEditor> 
		<webMaster>webmaster@writetheweb.com</webMaster> 
		<image>
			<title>WriteTheWeb</title> 
			<url>http://writetheweb.com/images/mynetscape88.gif</url> 
			<link>http://writetheweb.com</link> 
			<width>88</width> 
			<height>31</height> 
			<description>News for web users that write back</description> 
			</image>
		<item>
			<title>Giving the world a pluggable Gnutella</title> 
			<link>http://writetheweb.com/read.php?item=24</link> 
			<description>WorldOS is a framework on which to build programs that work like Freenet or Gnutella -allowing distributed applications using peer-to-peer routing.</description> 
			</item>
		<item>
			<title>Syndication discussions hot up</title> 
			<link>http://writetheweb.com/read.php?item=23</link> 
			<description>After a period of dormancy, the Syndication mailing list has become active again, with contributions from leaders in traditional media and Web syndication.</description> 
			</item>
		<item>
			<title>Personal web server integrates file sharing and messaging</title> 
			<link>http://writetheweb.com/read.php?item=22</link> 
			<description>The Magi Project is an innovative project to create a combined personal web server and messaging system that enables the sharing and synchronization of information across desktop, laptop and palmtop devices.</description> 
			</item>
		<item>
			<title>Syndication and Metadata</title> 
			<link>http://writetheweb.com/read.php?item=21</link> 
			<description>RSS is probably the best known metadata format around. RDF is probably one of the least understood. In this essay, published on my O'Reilly Network weblog, I argue that the next generation of RSS should be based on RDF.</description> 
			</item>
		<item>
			<title>UK bloggers get organised</title> 
			<link>http://writetheweb.com/read.php?item=20</link> 
			<description>Looks like the weblogs scene is gathering pace beyond the shores of the US. There's now a UK-specific page on weblogs.com, and a mailing list at egroups.</description> 
			</item>
		<item>
			<title>Yournamehere.com more important than anything</title> 
			<link>http://writetheweb.com/read.php?item=19</link> 
			<description>Whatever you're publishing on the web, your site name is the most valuable asset you have, according to Carl Steadman.</description> 
			</item>
		</channel>
	</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "0.91", a)
	assert.Equal(t, a.Charset, "ISO-8859-1", a)

}

func Test_RssFeed_010(t *testing.T) {
	// itunes
	// https://help.apple.com/itc/podcasts_connect/#/itcbaf351599
	s := `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>Hiking Treks</title>
    <link>https://www.apple.com/itunes/podcasts/</link>
    <language>en-us</language>
    <copyright>&#169; 2020 John Appleseed</copyright>
    <itunes:author>The Sunset Explorers</itunes:author>
    <description>
      Love to get outdoors and discover nature&apos;s treasures? Hiking Treks is the
      show for you. We review hikes and excursions, review outdoor gear and interview
      a variety of naturalists and adventurers. Look for new episodes each week.
    </description>
    <itunes:type>serial</itunes:type>
    <itunes:owner>
      <itunes:name>Sunset Explorers</itunes:name>
      <itunes:email>mountainscape@icloud.com</itunes:email>
    </itunes:owner>
    <itunes:image
      href="https://applehosted.podcasts.apple.com/hiking_treks/artwork.png"
    />
	<image>
		<title>title</title>
	</image>
    <itunes:category text="Sports">
      <itunes:category text="Wilderness"/>
    </itunes:category>
    <itunes:explicit>false</itunes:explicit>
    <item>
      <itunes:episodeType>trailer</itunes:episodeType>
      <itunes:title>Hiking Treks Trailer</itunes:title>
      <description>
          <![CDATA[The Sunset Explorers share tips, techniques and recommendations for
          great hikes and adventures around the United States. Listen on 
          <a href="https://www.apple.com/itunes/podcasts/">Apple Podcasts</a>.]]>
      </description>
      <enclosure 
        length="498537" 
        type="audio/mpeg" 
        url="http://example.com/podcasts/everything/AllAboutEverythingEpisode4.mp3"
      />
      <guid>aae20190418</guid>
      <pubDate>Tue, 8 Jan 2019 01:15:00 GMT</pubDate>
      <itunes:duration>1079</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>4</itunes:episode>
      <itunes:season>2</itunes:season>
      <title>S02 EP04 Mt. Hood, Oregon</title>
      <description>
        Tips for trekking around the tallest mountain in Oregon
      </description>
      <enclosure
        length="8727310" 
        type="audio/x-m4a" 
        url="http://example.com/podcasts/everything/mthood.m4a"
      />
      <guid>aae20190606</guid>
      <pubDate>Tue, 07 May 2019 12:00:00 GMT</pubDate>
      <itunes:duration>1024</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>3</itunes:episode>
      <itunes:season>2</itunes:season>
      <title>S02 EP03 Bouldering Around Boulder</title>
      <description>
        We explore fun walks to climbing areas about the beautiful Colorado city of Boulder.
      </description>
      <itunes:image
        href="http://example.com/podcasts/everything/AllAboutEverything/Episode2.jpg"
      />
      <link>href="http://example.com/podcasts/everything/</link>
      <enclosure 
        length="5650889" 
        type="video/mp4" 
        url="http://example.com/podcasts/boulder.mp4"
      />
      <guid>aae20190530</guid>
      <pubDate>Tue, 30 Apr 2019 13:00:00 EST</pubDate>
      <itunes:duration>3627</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>2</itunes:episode>
      <itunes:season>2</itunes:season>
      <title>S02 EP02 Caribou Mountain, Maine</title>
      <description>
        Put your fitness to the test with this invigorating hill climb.
      </description>
      <itunes:image
        href="http://example.com/podcasts/everything/AllAboutEverything/Episode3.jpg"
      />
      <enclosure 
        length="5650889"
        type="audio/x-m4v" 
        url="http://example.com/podcasts/everything/caribou.m4v"
      />
      <guid>aae20190523</guid>
      <pubDate>Tue, 23 May 2019 02:00:00 -0700</pubDate>
      <itunes:duration>2434</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>1</itunes:episode>
      <itunes:season>2</itunes:season>
      <title>S02 EP01 Stawamus Chief</title>
      <description>
        We tackle Stawamus Chief outside of Vancouver, BC and you should too!
      </description>
      <enclosure
        length="498537" 
        type="audio/mpeg" 
        url="http://example.com/podcasts/everything/AllAboutEverythingEpisode4.mp3"
      />
      <guid>aae20190516</guid>
      <pubDate>2019-02-16T07:00:00.000Z</pubDate>
      <itunes:duration>13:24</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>4</itunes:episode>
      <itunes:season>1</itunes:season>
      <title>S01 EP04 Kuliouou Ridge Trail</title>
      <description>
        Oahu, Hawaii, has some picturesque hikes and this is one of the best!
      </description>
      <enclosure
        length="498537" 
        type="audio/mpeg" 
        url="http://example.com/podcasts/everything/AllAboutEverythingEpisode4.mp3"
      />
      <guid>aae20190509</guid>
      <pubDate>Tue, 27 Nov 2018 01:15:00 +0000</pubDate>
      <itunes:duration>929</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>3</itunes:episode>
      <itunes:season>1</itunes:season>
      <title>S01 EP03 Blood Mountain Loop</title>
      <description>
        Hiking the Appalachian Trail and Freeman Trail in Georgia
      </description>
      <enclosure 
        length="498537" 
        type="audio/mpeg" 
        url="http://example.com/podcasts/everything/AllAboutEverythingEpisode4.mp3"
      />
      <guid>aae20190502</guid>
      <pubDate>Tue, 23 Oct 2018 01:15:00 +0000</pubDate>
      <itunes:duration>1440</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>2</itunes:episode>
      <itunes:season>1</itunes:season>
      <title>S01 EP02 Garden of the Gods Wilderness</title>
      <description>
        Wilderness Area Garden of the Gods in Illinois is a delightful spot for 
        an extended hike.
      </description>
      <enclosure 
        length="498537" 
        type="audio/mpeg" 
        url="http://example.com/podcasts/everything/AllAboutEverythingEpisode4.mp3"
      />
      <guid>aae20190425</guid>
      <pubDate>Tue, 18 Sep 2018 01:15:00 +0000</pubDate>
      <itunes:duration>839</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
    <item>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>1</itunes:episode>
      <itunes:season>1</itunes:season>
      <title>S01 EP01 Upper Priest Lake Trail to Continental Creek Trail</title>
      <description>
        We check out this powerfully scenic hike following the river in the Idaho
        Panhandle National Forests.
      </description>
      <enclosure 
        length="498537" 
        type="audio/mpeg" 
        url="http://example.com/podcasts/everything/AllAboutEverythingEpisode4.mp3"
      />
      <guid>aae20190418a</guid>
      <pubDate>Tue, 14 Aug 2018 01:15:00 +0000</pubDate>
      <itunes:duration>1399</itunes:duration>
      <itunes:explicit>false</itunes:explicit>
    </item>
  </channel>
</rss>
`

	a, err := RssParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.XMLName.Local, "rss", a)
	assert.Equal(t, a.Version, "2.0", a)

}
